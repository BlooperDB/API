package nodes

import (
	"net/http"
	"time"

	"strconv"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/storage"
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

	blueprint, e := findBlueprintById(request.BlueprintId)

	if e != nil {
		return nil, e
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	i := blueprint.IncrementAndGetRevision()

	bpVersion, _ := strconv.Atoi(request.Blueprint[0:1])

	sha265 := utils.SHA265(request.Blueprint)

	if db.FindRevisionByChecksum(sha265) != nil {
		return nil, &utils.Error_blueprint_string_already_exists
	}

	revision := &db.Revision{
		BlueprintID:       request.BlueprintId,
		Revision:          i,
		Changes:           request.Changes,
		BlueprintVersion:  bpVersion,
		BlueprintChecksum: sha265,
	}

	revision.Save()

	storage.SaveRevision(revision.ID, request.Blueprint)

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
	reComment := reCommentData(comments)

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

func revisionToJSON(authUser *db.User, revision *db.Revision, getComments bool) (*Revision, *utils.ErrorResponse) {
	if revision == nil || revision.DeletedAt != nil {
		return nil, &utils.Error_revision_not_found
	}

	ratings := revision.GetRatings()
	thumbsUp, thumbsDown, userVote := 0, 0, 0

	for _, rating := range ratings {
		if rating.ThumbsUp {
			thumbsUp++
		} else {
			thumbsDown++
		}

		if authUser != nil && authUser.ID == rating.UserID {
			if rating.ThumbsUp {
				userVote = 1
			} else {
				userVote = 2
			}
		}
	}

	var reComment []*Comment

	if getComments {
		comments := revision.GetComments()
		reComment = reCommentData(comments)
	}

	blueprintString, err := storage.LoadRevision(revision.ID)

	if err != nil {
		return nil, &utils.Error_internal_error
	}

	return &Revision{
		Id:          revision.ID,
		Revision:    revision.Revision,
		Changes:     revision.Changes,
		CreatedAt:   revision.CreatedAt,
		UpdatedAt:   revision.UpdatedAt,
		BlueprintID: revision.BlueprintID,
		Blueprint:   blueprintString,
		ThumbsUp:    thumbsUp,
		ThumbsDown:  thumbsDown,
		UserVote:    userVote,
		Comments:    reComment,
		Version:     revision.BlueprintVersion,
	}, nil
}

func reRevisionData(authUser *db.User, revisions []*db.Revision, getComments bool) ([]*Revision, *utils.ErrorResponse) {
	reRevision := make([]*Revision, len(revisions))

	for i, revision := range revisions {
		if revision.DeletedAt != nil {
			continue
		}

		rev, err := revisionToJSON(authUser, revision, getComments)

		if err != nil {
			return nil, err
		}

		reRevision[i] = rev
	}

	return reRevision, nil
}
