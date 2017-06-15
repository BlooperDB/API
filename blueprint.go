package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type BlueprintSearchResponse struct {
	Hello string `json:"hello,omitempty"`
	World string `json:"world,omitempty"`
}

func Search(_ *http.Request, _ httprouter.Params) (interface{}, *ErrorResponse) {
	return BlueprintSearchResponse{
		Hello: "herro",
		World: "wowd",
	}, nil
}
