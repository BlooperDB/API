// @APIVersion 1.0.0
// @APITitle Blooper
// @APIDescription Factorio blueprint database
// @License MIT
// @LicenseUrl https://opensource.org/licenses/MIT

package blooper

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"encoding/xml"

	"github.com/FactorioDB/API/blueprint"
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
	Data    interface{}    `json:"data,omitempty"`
}

type GenericHandle func(*http.Request, httprouter.Params) GenericResponse

func ProcessResponse(handle GenericHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		response := handle(r, p)

		format := r.URL.Query().Get("format")
		pretty := r.URL.Query().Get("pretty") != ""

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

type DataHandle func(*http.Request, httprouter.Params) (interface{}, *ErrorResponse)

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

type RegisterRoute func(method string, path string, handle DataHandle)

func RouteHandler(router *httprouter.Router) RegisterRoute {
	return func(method string, path string, handle DataHandle) {
		router.Handle(method, path, ProcessResponse(DataHandler(handle)))
	}
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello")
	})

	blueprint.RegisterBlueprintRoutes(RouteHandler(router))

	log.Fatal(http.ListenAndServe(":8080", router))
}
