package nodes

import (
	"net/http"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/gorilla/mux"
)

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
	UserVote   int       `json:"user-vote"`
	Comments   []Comment `json:"comments"`
}

type Comment struct {
	Id      string `json:"id"`
	UserId  string `json:"user"`
	Date    int    `json:"date"`
	Message string `json:"message"`
	Updated int    `json:"updated"`
}

func RegisterBlueprintRoutes(router api.RegisterRoute) {
	router("GET", "/blueprints", getBlueprints)
	router("GET", "/blueprints/search/{query}", searchBlueprints)

	router("GET", "/blueprint/{blueprint}", getBlueprint)
	router("POST", "/blueprint/{blueprint}", postBlueprint)
	router("PUT", "/blueprint/{blueprint}", updateBlueprint)
	router("DELETE", "/blueprint/{blueprint}", deleteBlueprint)

	router("GET", "/blueprint/{blueprint}/comments", getComments)
	router("GET", "/blueprint/{blueprint}/comment/{comment}", getComment)
	router("POST", "/blueprint/{blueprint}/comment/{comment}", postComment)
	router("PUT", "/blueprint/{blueprint}/comment/{comment}", updateComment)
	router("DELETE", "/blueprint/{blueprint}/comment/{comment}", deleteComment)

	router("GET", "/blueprint/{blueprint}/versions", getVersions)
	router("GET", "/blueprint/{blueprint}/version/{version}", getVersion)
	router("POST", "/blueprint/{blueprint}/version/{version}", postVersion)
	router("PUT", "/blueprint/{blueprint}/version/{version}", updateVersion)
	router("DELETE", "/blueprint/{blueprint}/version/{version}", deleteVersion)
}

/*
Search for blueprints
*/
func searchBlueprints(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Get all blueprints (paged)
*/
func getBlueprints(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Get a specific blueprint
*/
func getBlueprint(r *http.Request) (interface{}, *api.ErrorResponse) {
	id := mux.Vars(r)["blueprint"]
	blueprint := db.GetBlueprintById(id)

	if blueprint == nil {
		return nil, &error_blueprint_not_found
	}

	versions := blueprint.GetVersions()
	reVersion := make([]Version, len(versions))

	for i := 0; i < len(versions); i++ {
		version := versions[i]

		ratings := version.GetRatings()
		thumbsUp, thumbsDown, userVote := 0, 0, 0

		for j := 0; j < len(ratings); j++ {
			rating := ratings[j]

			if rating.ThumbsUp {
				thumbsUp++
			} else {
				thumbsDown++
			}

			// TODO Check for user id
		}

		comments := version.GetComments()
		reComment := make([]Comment, len(comments))

		for j := 0; j < len(ratings); j++ {
			comment := comments[j]
			reComment[j] = Comment{
				Id:      comment.Id,
				UserId:  comment.UserId,
				Date:    comment.Date,
				Message: comment.Message,
				Updated: comment.Updated,
			}
		}

		reVersion[i] = Version{
			Id:         version.Id,
			Version:    version.Version,
			Changes:    version.Changes,
			Date:       version.Date,
			Blueprint:  version.Blueprint,
			ThumbsUp:   thumbsUp,
			ThumbsDown: thumbsDown,
			UserVote:   userVote,
			Comments:   reComment,
		}
	}

	tags := blueprint.GetTags()
	reTags := make([]string, len(tags))

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		reTags[i] = tag.Name
	}

	return BlueprintResponse{
		Id:          blueprint.Id,
		UserId:      blueprint.UserId,
		Name:        blueprint.Name,
		Description: blueprint.Description,
		Versions:    reVersion,
		Tags:        reTags,
	}, nil
}

/*
Post a new blueprint
*/
func postBlueprint(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Update a blueprint
*/
func updateBlueprint(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Delete a blueprint
*/
func deleteBlueprint(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Get all comments
*/
func getComments(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Get specific comment
*/
func getComment(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Post a comment
*/
func postComment(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Update a comment
*/
func updateComment(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Delete a comment
*/
func deleteComment(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Get all versions
*/
func getVersions(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Get specific version
*/
func getVersion(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Post a version
*/
func postVersion(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Update a version
*/
func updateVersion(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}

/*
Delete a version
*/
func deleteVersion(_ *http.Request) (interface{}, *api.ErrorResponse) {
	return nil, nil
}
