package api

import (
	"encoding/json"
	"encoding/xml"

	"net/http"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

type GenericHandle func(http.ResponseWriter, *http.Request) GenericResponse

func ProcessResponse(handle GenericHandle) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := handle(w, r)

		format := r.URL.Query().Get("format")
		pretty := len(r.URL.Query()["pretty"]) > 0

		switch format {
		default:
			w.Header().Set("Content-Type", "application/json")
		case "xml":
			w.Header().Set("Content-Type", "application/xml")
		}

		if response.Error != nil {
			w.WriteHeader(response.Error.Status)
		}

		switch format {
		default:
			encoder := json.NewEncoder(w)

			if pretty {
				encoder.SetIndent("", "    ")
			}

			encoder.Encode(response)
		case "xml":
			encoder := xml.NewEncoder(w)

			if pretty {
				encoder.Indent("", "    ")
			}

			encoder.Encode(response)
		}
	}
}

type DataHandle func(*http.Request) (interface{}, *ErrorResponse)

func DataHandler(handle DataHandle) GenericHandle {
	return func(w http.ResponseWriter, r *http.Request) GenericResponse {
		data, err := handle(r)

		return GenericResponse{
			Success: err == nil,
			Error:   err,
			Data:    data,
		}
	}
}

type RegisterRoute func(method string, path string, handle DataHandle)

func RouteHandler(router *mux.Router, prefix string) RegisterRoute {
	return func(method string, path string, handle DataHandle) {
		route := router.NewRoute()
		route.PathPrefix(prefix)
		route.Path(path)
		route.Methods(method)
		route.HandlerFunc(ProcessResponse(DataHandler(handle)))
	}
}
