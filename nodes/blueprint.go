package nodes

import (
	"net/http"

	"encoding/json"

	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

type BlueprintResponse struct {
	Id          uint      `json:"id"`
	UserId      uint      `json:"user"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created-at"`
	UpdatedAt   time.Time `json:"updated-at"`
}

func RegisterBlueprintRoutes(router api.RegisterRoute) {
	router("GET", "/blueprints", getBlueprints)
	router("GET", "/blueprints/search/{query}", searchBlueprints)

	router("POST", "/blueprint", api.AuthHandler(postBlueprint))
	router("GET", "/blueprint/{blueprint}", getBlueprint)
	router("PUT", "/blueprint/{blueprint}", api.AuthHandler(updateBlueprint))
	router("DELETE", "/blueprint/{blueprint}", api.AuthHandler(deleteBlueprint))

	router("GET", "/blueprint/{blueprint}/versions", getVersions)
}

/*
Search for blueprints
*/
func searchBlueprints(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	return nil, nil
}

type GetBlueprintsResponse struct {
	Blueprints []SmallBlueprintResponse `json:"blueprints"`
}

type SmallBlueprintResponse struct {
	Id          uint   `json:"id"`
	UserId      uint   `json:"user-id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
Get all blueprints (paged)
*/
func getBlueprints(_ *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprints := db.GetAllBlueprints()
	reBlueprint := make([]SmallBlueprintResponse, len(blueprints))

	for i := 0; i < len(blueprints); i++ {
		reBlueprint[i] = SmallBlueprintResponse{
			Id:          blueprints[i].ID,
			UserId:      blueprints[i].UserID,
			Name:        blueprints[i].Name,
			Description: blueprints[i].Description,
		}
	}

	return GetBlueprintsResponse{
		Blueprints: reBlueprint,
	}, nil
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
				if rating.UserID == authUser.ID {
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
				Id:        comment.ID,
				UserId:    comment.UserID,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
				Message:   comment.Message,
			}
		}

		reVersion[i] = Version{
			Id:         version.ID,
			Version:    version.Version,
			Changes:    version.Changes,
			CreatedAt:  version.CreatedAt,
			UpdatedAt:  version.UpdatedAt,
			Blueprint:  version.BlueprintString,
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
		Id:          blueprint.ID,
		UserId:      blueprint.UserID,
		Name:        blueprint.Name,
		Description: blueprint.Description,
		CreatedAt:   blueprint.CreatedAt,
		UpdatedAt:   blueprint.UpdatedAt,
		Versions:    reVersion,
		Tags:        reTags,
	}, nil
}

type PostBlueprintRequest struct {
	Name        string                      `json:"name";validate:"nonzero"`
	Description string                      `json:"description";validate:"nonzero"`
	Version     PostBlueprintRequestVersion `json:"version";validate:"nonzero"`
}

type PostBlueprintRequestVersion struct {
	Version   string `json:"version";validate:"nonzero"`
	Blueprint string `json:"blueprint";validate:"nonzero"`
}

type PostBlueprintResponse struct {
	BlueprintId uint `json:"blueprint-id"`
	VersionId   uint `json:"version-id"`
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

	if err = validator.Validate(request); err != nil {
		return nil, &utils.ErrorResponse{
			Code:    utils.Error_invalid_request_data.Code,
			Message: utils.Error_invalid_request_data.Message + ": " + err.Error(),
			Status:  utils.Error_invalid_request_data.Status,
		}
	}

	blueprint := db.Blueprint{
		UserID:      u.ID,
		Name:        request.Name,
		Description: request.Description,
	}

	blueprint.Save()

	version := db.Version{
		BlueprintID:     blueprint.ID,
		Version:         request.Version.Version,
		Changes:         "",
		BlueprintString: request.Version.Blueprint,
	}

	version.Save()

	return PostBlueprintResponse{
		BlueprintId: blueprint.ID,
		VersionId:   version.ID,
	}, nil
}

type PutBlueprintRequest struct {
	Name        string `json:"name";validate:"nonzero"`
	Description string `json:"description";validate:"nonzero"`
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

	if err = validator.Validate(request); err != nil {
		return nil, &utils.ErrorResponse{
			Code:    utils.Error_invalid_request_data.Code,
			Message: utils.Error_invalid_request_data.Message + ": " + err.Error(),
			Status:  utils.Error_invalid_request_data.Status,
		}
	}

	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	if blueprint.UserID != u.ID {
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
	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	blueprint.Delete()

	return nil, nil
}

type GetVersionsResponse struct {
	Versions []Version `json:"versions"`
}

/*
Get all versions
*/
func getVersions(r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprintId := mux.Vars(r)["blueprint"]

	blueprint := db.GetBlueprintById(blueprintId)

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
				if rating.UserID == authUser.ID {
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
				Id:        comment.ID,
				UserId:    comment.UserID,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
				Message:   comment.Message,
			}
		}

		reVersion[i] = Version{
			Id:         version.ID,
			Version:    version.Version,
			Changes:    version.Changes,
			CreatedAt:  version.CreatedAt,
			UpdatedAt:  version.UpdatedAt,
			Blueprint:  version.BlueprintString,
			ThumbsUp:   thumbsUp,
			ThumbsDown: thumbsDown,
			UserVote:   userVote,
			Comments:   reComment,
		}
	}

	return GetVersionsResponse{
		Versions: reVersion,
	}, nil
}
