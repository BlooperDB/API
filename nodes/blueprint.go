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

type BlueprintResponse struct {
	Id          uint        `json:"id"`
	UserId      uint        `json:"user"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Revisions   []*Revision `json:"revisions"`
	Tags        []string    `json:"tags"`
	CreatedAt   time.Time   `json:"created-at"`
	UpdatedAt   time.Time   `json:"updated-at"`
}

func RegisterBlueprintRoutes(router api.RegisterRoute) {
	router("GET", "/blueprints", getBlueprints)
	router("GET", "/blueprints/search/{query}", searchBlueprints)

	router("POST", "/blueprint", api.UsernameRequiredHandler(postBlueprint))
	router("GET", "/blueprint/{blueprint}", getBlueprint)
	router("PUT", "/blueprint/{blueprint}", api.UsernameRequiredHandler(updateBlueprint))
	router("DELETE", "/blueprint/{blueprint}", api.UsernameRequiredHandler(deleteBlueprint))

	router("GET", "/blueprint/{blueprint}/revisions", getRevisions)
	router("GET", "/blueprint/{blueprint}/revision/latest", getRevisionLatest)
	router("GET", "/blueprint/{blueprint}/revision/{revision}", getRevisionIncremental)
}

type SearchBlueprintsResponse struct {
	Blueprints []*SmallBlueprintResponse `json:"blueprints"`
}

/*
Search for blueprints
*/
func searchBlueprints(r *http.Request) (interface{}, *utils.ErrorResponse) {
	var (
		offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
		count, _  = strconv.Atoi(r.URL.Query().Get("count"))
	)

	if count == 0 {
		count = 20
	}
	if count > 100 {
		count = 100
	}

	query := mux.Vars(r)["query"]

	blueprints := db.SearchBlueprints(query, offset, count)
	reBlueprint := make([]*SmallBlueprintResponse, len(blueprints))

	for i, blueprint := range blueprints {
		var revId uint = 0
		if rev := blueprint.GetLatestRevision(); rev != nil {
			revId = rev.Revision
		}

		tags := blueprint.GetTags()
		reTags := make([]string, len(tags))

		for i, tag := range tags {
			reTags[i] = tag.Name
		}

		reBlueprint[i] = &SmallBlueprintResponse{
			Id:          blueprint.ID,
			Latest:      revId,
			UserId:      blueprint.UserID,
			Name:        blueprint.Name,
			Description: blueprint.Description,
			Tags:        reTags,
		}
	}

	return SearchBlueprintsResponse{
		Blueprints: reBlueprint,
	}, nil
}

type GetBlueprintsResponse struct {
	Blueprints []*SmallBlueprintResponse `json:"blueprints"`
}

type SmallBlueprintResponse struct {
	Id          uint     `json:"id"`
	Latest      uint     `json:"latest-revision"`
	UserId      uint     `json:"user-id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

/*
Get all blueprints (paged)
*/
func getBlueprints(r *http.Request) (interface{}, *utils.ErrorResponse) {
	var (
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		count, err0  = strconv.Atoi(r.URL.Query().Get("count"))
	)
	if err != nil || err0 != nil {
		return nil, &utils.Error_invalid_request_data
	}

	if count == 0 {
		count = 20
	}
	if count > 100 {
		count = 100
	}

	blueprints := db.GetAllBlueprints(offset, count)
	reBlueprint := make([]*SmallBlueprintResponse, len(blueprints))

	for i, blueprint := range blueprints {
		var revId uint = 0
		if rev := blueprint.GetLatestRevision(); rev != nil {
			revId = rev.Revision
		}

		tags := blueprint.GetTags()
		reTags := make([]string, len(tags))

		for i, tag := range tags {
			reTags[i] = tag.Name
		}

		reBlueprint[i] = &SmallBlueprintResponse{
			Id:          blueprint.ID,
			Latest:      revId,
			UserId:      blueprint.UserID,
			Name:        blueprint.Name,
			Description: blueprint.Description,
			Tags:        reTags,
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
	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	revisions := blueprint.GetRevisions()
	reRevision := make([]*Revision, len(revisions))

	authUser := db.GetAuthUser(r)

	for i, revision := range revisions {
		if revision.DeletedAt != nil {
			continue
		}
		rev, err := revisionToJSON(authUser, &revision)
		if err != nil {
			return nil, err
		}
		reRevision[i] = rev
	}

	tags := blueprint.GetTags()
	reTags := make([]string, len(tags))

	for i, tag := range tags {
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
	var request PostBlueprintRequest
	e := utils.ValidateRequestBody(r, &request)

	if e != nil {
		return nil, e
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
	var request PutBlueprintRequest
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

	blueprint.Name = request.Name
	blueprint.Description = request.Description

	blueprint.Save()

	return nil, nil
}

/*
Delete a blueprint
*/
func deleteBlueprint(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	if blueprint.UserID != u.ID {
		return nil, &utils.Error_no_access
	}

	blueprint.Delete()

	return nil, nil
}

type GetRevisionsResponse struct {
	Revisions []*Revision `json:"revisions"`
}

/*
Get all revisions
*/
func getRevisions(r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	revisions := blueprint.GetRevisions()
	reRevision := make([]*Revision, len(revisions))

	authUser := db.GetAuthUser(r)

	for i, revision := range revisions {
		rev, err := revisionToJSON(authUser, &revision)
		if err != nil {
			return nil, err
		}
		reRevision[i] = rev
	}

	return GetRevisionsResponse{
		Revisions: reRevision,
	}, nil
}

/*
Get latest revision from blueprint
*/
func getRevisionLatest(r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	authUser := db.GetAuthUser(r)
	revision := blueprint.GetLatestRevision()
	if revision == nil || revision.DeletedAt != nil {
		return nil, &utils.Error_revision_not_found
	}

	return revisionToJSON(authUser, revision)
}

/*
Get specific revision from blueprint
*/
func getRevisionIncremental(r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprint, e := parseBlueprint(r)

	if e != nil {
		return nil, e
	}

	revisionI, err := strconv.ParseUint(mux.Vars(r)["revision"], 10, 32)
	if err != nil {
		return nil, &utils.Error_revision_not_found
	}

	authUser := db.GetAuthUser(r)
	revision := blueprint.GetRevision(uint(revisionI))
	return revisionToJSON(authUser, revision)
}

func revisionToJSON(authUser *db.User, revision *db.Revision) (*Revision, *utils.ErrorResponse) {
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

	return &Revision{
		Id:          revision.ID,
		Revision:    revision.Revision,
		Changes:     revision.Changes,
		CreatedAt:   revision.CreatedAt,
		UpdatedAt:   revision.UpdatedAt,
		BlueprintID: revision.BlueprintID,
		Blueprint:   revision.BlueprintString,
		ThumbsUp:    thumbsUp,
		ThumbsDown:  thumbsDown,
		UserVote:    userVote,
		Comments:    reComment,
	}, nil
}

func parseBlueprint(r *http.Request) (*db.Blueprint, *utils.ErrorResponse) {
	blueprintId, err := strconv.ParseUint(mux.Vars(r)["blueprint"], 10, 32)

	if err != nil {
		return nil, &utils.Error_blueprint_not_found
	}

	blueprint := db.GetBlueprintById(uint(blueprintId))

	if blueprint == nil {
		return nil, &utils.Error_blueprint_not_found
	}

	return blueprint, nil
}
