package nodes

import (
	"net/http"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
)

func RegisterTagRoutes(router api.RegisterRoute) {
	router("GET", "/tags/popular", popularTags)
	router("GET", "/tags/autocomplete/{tag}", autocompleteTag)
}

func popularTags(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	tags := db.PopularTags()
	reTags := make([]string, len(tags))

	for i, tag := range tags {
		reTags[i] = tag.Name
	}

	return AutocompleteTagResponse{
		Tags: reTags,
	}, nil
}

type AutocompleteTagResponse struct {
	Tags []string `json:"tags"`
}

func autocompleteTag(r *http.Request) (interface{}, *utils.ErrorResponse) {
	tag := mux.Vars(r)["tag"]

	tags := db.AutocompleteTag(tag)
	reTags := make([]string, len(tags))

	for i, tag := range tags {
		reTags[i] = tag.Name
	}

	return AutocompleteTagResponse{
		Tags: reTags,
	}, nil
}
