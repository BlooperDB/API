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

type TagListResponse struct {
	Tags []string `json:"tags"`
}

func popularTags(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	tags := db.PopularTags()
	reTags := make([]string, len(tags))

	for i, tag := range tags {
		reTags[i] = tag.Name
	}

	return TagListResponse{
		Tags: reTags,
	}, nil
}

func autocompleteTag(r *http.Request) (interface{}, *utils.ErrorResponse) {
	tag := mux.Vars(r)["tag"]

	tags := db.AutocompleteTag(tag)
	reTags := make([]string, len(tags))

	for i, tag := range tags {
		reTags[i] = tag.Name
	}

	return TagListResponse{
		Tags: reTags,
	}, nil
}

func reTagData(tags []*db.Tag) []string {
	reTags := make([]string, len(tags))

	for i, tag := range tags {
		reTags[i] = tag.Name
	}

	return reTags
}
