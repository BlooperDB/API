package blueprint

import (
	"net/http"

	"github.com/FactorioDB/API/api"
)

//noinspection GoNameStartsWithPackageName
type BlueprintResponse struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
	Tags        []string  `json:"tags"`
}

type Version struct {
	Id         string    `json:"id"`
	Version    string    `json:"version"`
	Changes    string    `json:"changes"`
	Date       int       `json:"date"`
	Blueprint  string    `json:"blueprint"`
	ThumbsUp   int       `json:"thumbs_up"`
	ThumbsDown int       `json:"thumbs_down"`
	Comments   []Comment `json:"comments"`
}

type Comment struct {
	Id      string `json:"id"`
	UserId  string `json:"user"`
	Date    int    `json:"date"`
	Message string `json:"message"`
}

func RegisterBlueprintRoutes(router api.RegisterRoute) {
	router("GET", "/{id}", get)
	router("GET", "/search/{query}", search)
}

func get(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

func search(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}
