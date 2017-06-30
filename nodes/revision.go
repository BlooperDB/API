package nodes

import (
	"net/http"
	"time"

	"strconv"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
)

type Revision struct {
	Id          uint       `json:"id"`
	Revision    uint       `json:"revision"`
	Changes     string     `json:"changes"`
	CreatedAt   time.Time  `json:"created-at"`
	UpdatedAt   time.Time  `json:"updated-at"`
	BlueprintID uint       `json:"blueprint-id"`
	Blueprint   string     `json:"blueprint"`
	ThumbsUp    int        `json:"thumbs-up"`
	ThumbsDown  int        `json:"thumbs-down"`
	UserVote    int        `json:"user-vote"`
	Comments    []*Comment `json:"comments,omitempty"`
	Version     int        `json:"version"`
}

func RegisterRevisionRoutes(router api.RegisterRoute) {
	router("POST", "/revision", api.AuthHandler(postRevision, true))
	router("GET", "/revision/{revision}", getRevision)
	router("PUT", "/revision/{revision}", api.AuthHandler(updateRevision, true))
	router("DELETE", "/revision/{revision}", api.AuthHandler(deleteRevision, true))

	router("GET", "/revision/{revision}/comments", getRevisionComments)

	router("POST", "/revision/{revision}/rating", api.AuthHandler(postRevisionRating, true))
	router("DELETE", "/revision/{revision}/rating", api.AuthHandler(deleteRevisionRating, true))
}

/*
Get specific revision
*/
func getRevision(r *http.Request) (interface{}, *utils.ErrorResponse) {
	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	getComments := len(r.URL.Query()["comments"]) > 0

	authUser := db.GetAuthUser(r)
	return revisionToJSON(authUser, revision, getComments)
}

type PostRevisionRequest struct {
	BlueprintId uint   `json:"blueprint-id" validate:"nonzero"`
	Changes     string `json:"changes" validate:"nonzero"`
	Blueprint   string `json:"blueprint" validate:"nonzero,blueprint_string"`
}

type PostRevisionResponse struct {
	// Global unique revision identifier
	RevisionId uint `json:"revision-id"`

	// Blueprint incremental version
	Revision uint `json:"revision"`
}

/*
Post a revision
*/
func postRevision(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request PostRevisionRequest
	e := utils.ValidateRequestBody(r, &request)

	if e != nil {
		return nil, e
	}

	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	i := blueprint.IncrementAndGetRevision()

	bpVersion, _ := strconv.Atoi(request.Blueprint[0:1])

	revision := &db.Revision{
		BlueprintID:      request.BlueprintId,
		Revision:         i,
		Changes:          request.Changes,
		BlueprintString:  request.Blueprint,
		BlueprintVersion: bpVersion,
	}

	revision.Save()

	return PostRevisionResponse{
		RevisionId: revision.ID,
		Revision:   revision.Revision,
	}, nil
}

type PutRevisionRequest struct {
	Changes string `json:"changes" validate:"nonzero"`
}

/*
Update a revision
*/
func updateRevision(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request PutRevisionRequest
	e := utils.ValidateRequestBody(r, &request)

	if e != nil {
		return nil, e
	}

	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	revision.Changes = request.Changes
	revision.Save()

	return nil, nil
}

/*
Delete a revision
*/
func deleteRevision(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	revision.Delete()

	if blueprint.CountRevisions() == 0 {
		blueprint.Delete()
	}

	return nil, nil
}

type GetRevisionCommentsResponse struct {
	Comments []*Comment `json:"comments"`
}

/*
Get all comments
*/
func getRevisionComments(r *http.Request) (interface{}, *utils.ErrorResponse) {
	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	comments := revision.GetComments()
	reComment := make([]*Comment, len(comments))

	for i, comment := range comments {
		reComment[i] = &Comment{
			Id:        comment.ID,
			UserId:    comment.UserID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Message:   comment.Message,
		}
	}

	return GetRevisionCommentsResponse{
		Comments: reComment,
	}, nil
}

type PostRevisionRating struct {
	ThumbsUp bool `json:"thumbs-up"`
}

/*
Post rating revision
*/
func postRevisionRating(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request PostRevisionRating
	utils.ValidateRequestBody(r, &request)

	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	rating := db.FindRating(u.ID, revision.ID)

	rating.UserID = u.ID
	rating.RevisionID = revision.ID
	rating.ThumbsUp = request.ThumbsUp
	rating.DeletedAt = nil
	rating.Save()

	return nil, nil
}

/*
Post rating revision
*/
func deleteRevisionRating(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request PostRevisionRating
	utils.ValidateRequestBody(r, &request)

	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	rating := db.FindRating(u.ID, revision.ID)

	if rating.ID == 0 {
		return nil, &utils.Error_rating_not_found
	}

	rating.Delete()

	return nil, nil
}

func parseRevision(r *http.Request) (*db.Revision, *utils.ErrorResponse) {
	revisionId, err := strconv.ParseUint(mux.Vars(r)["revision"], 10, 32)

	if err != nil {
		return nil, &utils.Error_revision_not_found
	}

	revision := db.GetRevisionById(uint(revisionId))

	if revision == nil {
		return nil, &utils.Error_revision_not_found
	}

	return revision, nil
}
