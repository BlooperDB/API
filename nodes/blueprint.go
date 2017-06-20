package nodes

import (
	"net/http"

	"encoding/json"

	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
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
	router("POST", "/blueprint/{blueprint}", api.AuthHandler(postBlueprint))
	router("PUT", "/blueprint/{blueprint}", api.AuthHandler(updateBlueprint))
	router("DELETE", "/blueprint/{blueprint}", api.AuthHandler(deleteBlueprint))

	router("GET", "/blueprint/{blueprint}/comments", getComments)
	router("GET", "/blueprint/{blueprint}/comment/{comment}", getComment)
	router("POST", "/blueprint/{blueprint}/comment/{comment}", api.AuthHandler(postComment))
	router("PUT", "/blueprint/{blueprint}/comment/{comment}", api.AuthHandler(updateComment))
	router("DELETE", "/blueprint/{blueprint}/comment/{comment}", api.AuthHandler(deleteComment))

	router("GET", "/blueprint/{blueprint}/versions", getVersions)
	router("GET", "/blueprint/{blueprint}/version/{version}", getVersion)
	router("POST", "/blueprint/{blueprint}/version/{version}", api.AuthHandler(postVersion))
	router("PUT", "/blueprint/{blueprint}/version/{version}", api.AuthHandler(updateVersion))
	router("DELETE", "/blueprint/{blueprint}/version/{version}", api.AuthHandler(deleteVersion))
}

/*
Search for blueprints
*/
func searchBlueprints(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Get all blueprints (paged)
*/
func getBlueprints(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Get a specific blueprint
*/
func getBlueprint(r *http.Request) (interface{}, *utils.ErrorResponse) {
	id := mux.Vars(r)["blueprint"]
	blueprint := db.GetBlueprintById(id)

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	versions := blueprint.GetVersions()
	reVersion := make([]Version, len(versions))

	authUser := db.GetAuthUser(r)

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

			if authUser != nil {
				if rating.UserId == authUser.Id {
					if rating.ThumbsUp {
						userVote = 1
					} else {
						userVote = 2
					}
				}
			}
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

type PostBlueprintRequest struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Version     PostBlueprintRequestVersion `json:"version"`
}

type PostBlueprintRequestVersion struct {
	Version   string `json:"version"`
	Blueprint string `json:"blueprint"`
}

type PostBlueprintResponse struct {
	BlueprintId string `json:"blueprint-id"`
	VersionId   string `json:"version-id"`
}

/*
Post a new blueprint
*/
func postBlueprint(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PostBlueprintRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	blueprint := db.Blueprint{
		Id:          utils.GenerateRandomId(),
		UserId:      u.Id,
		Name:        request.Name,
		Description: request.Description,
	}

	blueprint.Save()

	version := db.Version{
		Id:          utils.GenerateRandomId(),
		BlueprintId: blueprint.Id,
		Version:     request.Version.Version,
		Changes:     "",
		Date:        time.Now().Unix(),
		Blueprint:   request.Version.Blueprint,
	}

	version.Save()

	return PostBlueprintResponse{
		BlueprintId: blueprint.Id,
		VersionId:   version.Id,
	}, nil
}

type PutBlueprintRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
Update a blueprint
*/
func updateBlueprint(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PutBlueprintRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint.UserId != u.Id {
		return nil, &utils.Error_no_access
	}

	blueprint.Name = request.Name
	blueprint.Description = request.Description

	blueprint.Save()

	return nil, nil
}

/*
Delete a blueprint
*/
func deleteBlueprint(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PutBlueprintRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint.UserId != u.Id {
		return nil, &utils.Error_no_access
	}

	blueprint.Delete()

	return nil, nil
}

/*
Get all comments
*/
func getComments(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Get specific comment
*/
func getComment(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Post a comment
*/
func postComment(u *db.User, _ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Update a comment
*/
func updateComment(u *db.User, _ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Delete a comment
*/
func deleteComment(u *db.User, _ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Get all versions
*/
func getVersions(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Get specific version
*/
func getVersion(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Post a version
*/
func postVersion(u *db.User, _ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Update a version
*/
func updateVersion(u *db.User, _ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

/*
Delete a version
*/
func deleteVersion(u *db.User, _ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}
