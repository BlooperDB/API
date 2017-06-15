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

	"github.com/julienschmidt/httprouter"
)

type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type GenericResponse struct {
	Success bool           `json:"success,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    *interface{}   `json:"data,omitempty"`
}

type DataHandle func(*http.Request, httprouter.Params) (interface{}, *ErrorResponse)

func DataHandler(handle DataHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, err := handle(r, p)

		response := GenericResponse{
			Success: err == nil,
			Error:   err,
			Data:    &data,
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		encoder.Encode(response)
	}
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello")
	})

	router.GET("/blueprint/search", DataHandler(Search))

	log.Fatal(http.ListenAndServe(":8080", router))
}
