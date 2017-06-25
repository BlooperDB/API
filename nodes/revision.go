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
	Comments    []*Comment `json:"comments"`
}

func RegisterRevisionRoutes(router api.RegisterRoute) {
	router("POST", "/revision", api.UsernameRequiredHandler(postRevision))
	router("GET", "/revision/{revision}", getRevision)
	router("PUT", "/revision/{revision}", api.UsernameRequiredHandler(updateRevision))
	router("DELETE", "/revision/{revision}", api.UsernameRequiredHandler(deleteRevision))
	router("GET", "/revision/{revision}/comments", getRevisionComments)
}

/*
Get specific revision
*/
func getRevision(r *http.Request) (interface{}, *utils.ErrorResponse) {
	revision, e := parseRevision(r)

	if e != nil {
		return nil, e
	}

	authUser := db.GetAuthUser(r)
	return revisionToJSON(authUser, revision)
}

type PostRevisionRequest struct {
	BlueprintId uint   `json:"blueprint-id" validate:"nonzero"`
	Changes     string `json:"changes" validate:"nonzero"`
	Blueprint   string `json:"blueprint" validate:"nonzero"`
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

	revision := &db.Revision{
		BlueprintID:     request.BlueprintId,
		Revision:        i,
		Changes:         request.Changes,
		BlueprintString: request.Blueprint,
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
