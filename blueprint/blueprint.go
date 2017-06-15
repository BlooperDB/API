package blueprint

import (
	"net/http"

	"github.com/FactorioDB/API"
	"github.com/julienschmidt/httprouter"
)

type SearchResponse struct {
	Hello string `json:"hello,omitempty"`
	World string `json:"world,omitempty"`
}

func RegisterBlueprintRoutes(router blooper.RegisterRoute) {
	router("GET", "/blueprint/search", search)
}

func search(_ *http.Request, _ httprouter.Params) (interface{}, *blooper.ErrorResponse) {
	return SearchResponse{
		Hello: "herro",
		World: "wowd",
	}, nil
}
