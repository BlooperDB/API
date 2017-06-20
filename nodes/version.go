package nodes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
)

func RegisterVersionRoutes(router api.RegisterRoute) {
	router("GET", "/version/{version}", getVersion)
	router("POST", "/version/{version}", api.AuthHandler(postVersion))
	router("PUT", "/version/{version}", api.AuthHandler(updateVersion))
	router("DELETE", "/version/{version}", api.AuthHandler(deleteVersion))
}

/*
Get specific version
*/
func getVersion(r *http.Request) (interface{}, *utils.ErrorResponse) {
	versionId := mux.Vars(r)["version"]

	version := db.GetVersionById(versionId)

	ratings := version.GetRatings()
	thumbsUp, thumbsDown, userVote := 0, 0, 0

	authUser := db.GetAuthUser(r)

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

	for j := 0; j < len(comments); j++ {
		comment := comments[j]
		reComment[j] = Comment{
			Id:      comment.Id,
			UserId:  comment.UserId,
			Date:    comment.Date,
			Message: comment.Message,
			Updated: comment.Updated,
		}
	}

	return Version{
		Id:         version.Id,
		Version:    version.Version,
		Changes:    version.Changes,
		Date:       version.Date,
		Blueprint:  version.Blueprint,
		ThumbsUp:   thumbsUp,
		ThumbsDown: thumbsDown,
		UserVote:   userVote,
		Comments:   reComment,
	}, nil
}

type PostVersionRequest struct {
	BlueprintId string `json:"blueprint-id"`
	Version     string `json:"version"`
	Changes     string `json:"changes"`
	Blueprint   string `json:"blueprint"`
}

type PostVersionResponse struct {
	Id string `json:"id"`
}

/*
Post a version
*/
func postVersion(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PostVersionRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint.UserId != u.Id {
		return nil, &utils.Error_no_access
	}

	version := db.Version{
		Id:          utils.GenerateRandomId(),
		BlueprintId: request.BlueprintId,
		Version:     request.Version,
		Changes:     request.Changes,
		Date:        time.Now().Unix(),
		Blueprint:   request.Blueprint,
	}

	version.Save()

	return PostVersionResponse{
		Id: version.Id,
	}, nil
}

type PutVersionRequest struct {
	Version   string `json:"version"`
	Changes   string `json:"changes"`
	Blueprint string `json:"blueprint"`
}

/*
Update a version
*/
func updateVersion(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PutVersionRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint.UserId != u.Id {
		return nil, &utils.Error_no_access
	}

	versionId := mux.Vars(r)["version"]

	version := blueprint.GetVersion(versionId)

	version.Version = request.Version
	version.Changes = request.Changes
	version.Blueprint = request.Blueprint

	return nil, nil
}

/*
Delete a version
*/
func deleteVersion(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint.UserId != u.Id {
		return nil, &utils.Error_no_access
	}

	versionId := mux.Vars(r)["version"]

	version := blueprint.GetVersion(versionId)

	version.Delete()

	return nil, nil
}
