// @APIVersion 1.0.0
// @APITitle Blooper
// @APIDescription Factorio blueprint database
// @License MIT
// @LicenseUrl https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"encoding/xml"

	"github.com/julienschmidt/httprouter"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    *interface{}   `json:"data,omitempty"`
}

type GenericHandle func(*http.Request, httprouter.Params) GenericResponse

type DataHandle func(*http.Request, httprouter.Params) (interface{}, *ErrorResponse)

func ProcessResponse(handle GenericHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		response := handle(r, p)

		if response.Error != nil {
			w.WriteHeader(response.Error.Status)
		}

		format := p.ByName("format")

		switch format {
		case "json":
		default:
			w.Header().Set("Content-Type", "application/json")

			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "    ")
			encoder.Encode(response)
			break
		case "xml":
			w.Header().Set("Content-Type", "application/xml")

			encoder := xml.NewEncoder(w)
			encoder.Indent("", "    ")
			encoder.Encode(response)
			break
		}
	}
}

func DataHandler(handle DataHandle) GenericHandle {
	return func(r *http.Request, p httprouter.Params) GenericResponse {
		data, err := handle(r, p)

		return GenericResponse{
			Success: err == nil,
			Error:   err,
			Data:    &data,
		}
	}
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello")
	})

	router.GET("/blueprint/search", ProcessResponse(DataHandler(Search)))

	log.Fatal(http.ListenAndServe(":8080", router))
}
