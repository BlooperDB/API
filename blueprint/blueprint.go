package blueprint

import (
	"net/http"

	"github.com/FactorioDB/API/api"
	"github.com/julienschmidt/httprouter"
)

type SearchResponse struct {
	Hello string `json:"hello,omitempty"`
	World string `json:"world,omitempty"`
}

func RegisterBlueprintRoutes(router api.RegisterRoute) {
	router("GET", "/blueprint/search", search)
}

func search(_ *http.Request, _ httprouter.Params) (interface{}, *api.ErrorResponse) {
	return SearchResponse{
		Hello: "herro",
		World: "wowd",
	}, nil
}
