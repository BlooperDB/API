package nodes

import (
	"net/http"

	"encoding/json"

	"time"

	"strconv"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

type BlueprintResponse struct {
	Id          uint       `json:"id"`
	UserId      uint       `json:"user"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Revisions   []Revision `json:"revisions"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created-at"`
	UpdatedAt   time.Time  `json:"updated-at"`
}

func RegisterBlueprintRoutes(router api.RegisterRoute) {
	router("GET", "/blueprints", getBlueprints)
	router("GET", "/blueprints/search/{query}", searchBlueprints)

	router("POST", "/blueprint", api.AuthHandler(postBlueprint))
	router("GET", "/blueprint/{blueprint}", getBlueprint)
	router("PUT", "/blueprint/{blueprint}", api.AuthHandler(updateBlueprint))
	router("DELETE", "/blueprint/{blueprint}", api.AuthHandler(deleteBlueprint))

	router("GET", "/blueprint/{blueprint}/revisions", getRevisions)
	router("GET", "/blueprint/{blueprint}/revision/{revision}", getRevisionIncremental)
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
	blueprintId, _ := strconv.ParseUint(mux.Vars(r)["blueprint"], 10, 32)

	blueprint := db.GetBlueprintById(uint(blueprintId))

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	revisions := blueprint.GetRevisions()
	reRevision := make([]Revision, len(revisions))

	authUser := db.GetAuthUser(r)

	for i := 0; i < len(revisions); i++ {
		revision := revisions[i]

		ratings := revision.GetRatings()
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

		comments := revision.GetComments()
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

		reRevision[i] = Revision{
			Id:         revision.ID,
			Revision:   revision.Revision,
			Changes:    revision.Changes,
			CreatedAt:  revision.CreatedAt,
			UpdatedAt:  revision.UpdatedAt,
			Blueprint:  revision.BlueprintString,
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
		Revisions:   reRevision,
		Tags:        reTags,
	}, nil
}

type PostBlueprintRequest struct {
	Name            string `json:"name";validate:"nonzero"`
	Description     string `json:"description";validate:"nonzero"`
	BlueprintString string `json:"blueprint-string";validate:"nonzero"`
}

type PostBlueprintResponse struct {
	BlueprintId uint `json:"blueprint-id"`

	// Global unique revision identifier
	RevisionId uint `json:"revision-id"`

	// Blueprint incremental version
	Revision uint `json:"revision"`
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

	blueprint := &db.Blueprint{
		UserID:       u.ID,
		Name:         request.Name,
		Description:  request.Description,
		LastRevision: 1,
	}

	blueprint.Save()

	revision := &db.Revision{
		BlueprintID:     blueprint.ID,
		Revision:        blueprint.LastRevision,
		Changes:         "",
		BlueprintString: request.BlueprintString,
	}

	revision.Save()

	return PostBlueprintResponse{
		BlueprintId: blueprint.ID,
		RevisionId:  revision.ID,
		Revision:    revision.Revision,
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

	blueprintId, _ := strconv.ParseUint(mux.Vars(r)["blueprint"], 10, 32)

	blueprint := db.GetBlueprintById(uint(blueprintId))

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
	blueprintId, _ := strconv.ParseUint(mux.Vars(r)["blueprint"], 10, 32)

	blueprint := db.GetBlueprintById(uint(blueprintId))

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	blueprint.Delete()

	return nil, nil
}

type GetRevisionsResponse struct {
	Revisions []Revision `json:"revisions"`
}

/*
Get all revisions
*/
func getRevisions(r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprintId, _ := strconv.ParseUint(mux.Vars(r)["blueprint"], 10, 32)

	blueprint := db.GetBlueprintById(uint(blueprintId))

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	revisions := blueprint.GetRevisions()
	reRevision := make([]Revision, len(revisions))

	authUser := db.GetAuthUser(r)

	for i := 0; i < len(revisions); i++ {
		revision := revisions[i]

		ratings := revision.GetRatings()
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

		comments := revision.GetComments()
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

		reRevision[i] = Revision{
			Id:         revision.ID,
			Revision:   revision.Revision,
			Changes:    revision.Changes,
			CreatedAt:  revision.CreatedAt,
			UpdatedAt:  revision.UpdatedAt,
			Blueprint:  revision.BlueprintString,
			ThumbsUp:   thumbsUp,
			ThumbsDown: thumbsDown,
			UserVote:   userVote,
			Comments:   reComment,
		}
	}

	return GetRevisionsResponse{
		Revisions: reRevision,
	}, nil
}

/*
Get specific revision from blueprint
*/
func getRevisionIncremental(r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprintId, _ := strconv.ParseUint(mux.Vars(r)["blueprint"], 10, 32)

	blueprint := db.GetBlueprintById(uint(blueprintId))

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	revisionI, _ := strconv.ParseUint(mux.Vars(r)["revision"], 10, 32)

	revision := blueprint.GetRevision(uint(revisionI))

	if revision == nil {
		return nil, &utils.Error_revision_not_found
	}

	authUser := db.GetAuthUser(r)

	ratings := revision.GetRatings()
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

	comments := revision.GetComments()
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

	return Revision{
		Id:         revision.ID,
		Revision:   revision.Revision,
		Changes:    revision.Changes,
		CreatedAt:  revision.CreatedAt,
		UpdatedAt:  revision.UpdatedAt,
		Blueprint:  revision.BlueprintString,
		ThumbsUp:   thumbsUp,
		ThumbsDown: thumbsDown,
		UserVote:   userVote,
		Comments:   reComment,
	}, nil
}
