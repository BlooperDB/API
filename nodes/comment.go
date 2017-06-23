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

type Comment struct {
	Id         uint      `json:"id"`
	UserId     uint      `json:"user"`
	CreatedAt  time.Time `json:"created-at"`
	UpdatedAt  time.Time `json:"updated-at"`
	Message    string    `json:"message"`
	RevisionId uint      `json:"revision-id"`
}

func RegisterCommentRoutes(router api.RegisterRoute) {
	router("POST", "/comment", api.UsernameRequiredHandler(postComment))
	router("GET", "/comment/{comment}", getComment)
	router("PUT", "/comment/{comment}", api.UsernameRequiredHandler(updateComment))
	router("DELETE", "/comment/{comment}", api.UsernameRequiredHandler(deleteComment))
}

/*
Get specific comment
*/
func getComment(r *http.Request) (interface{}, *utils.ErrorResponse) {
	comment, err := parseComment(r)

	if err != nil {
		return nil, err
	}

	return Comment{
		Id:         comment.ID,
		UserId:     comment.UserID,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
		Message:    comment.Message,
		RevisionId: comment.RevisionID,
	}, nil
}

type PostCommentRequest struct {
	Message    string `json:"message";validate:"nonzero"`
	RevisionId uint   `json:"revision-id";validate:"nonzero"`
}

type PostCommentResponse struct {
	CommentId uint `json:"comment-id"`
}

/*
Post a comment
*/
func postComment(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request PostCommentRequest
	e := utils.ValidateRequestBody(r, &request)

	if e != nil {
		return nil, e
	}

	comment := &db.Comment{
		RevisionID: request.RevisionId,
		UserID:     u.ID,
		Message:    request.Message,
	}

	comment.Save()

	return PostCommentResponse{
		CommentId: comment.ID,
	}, nil
}

type PutCommentRequest struct {
	Message string `json:"message";validate:"nonzero"`
}

/*
Update a comment
*/
func updateComment(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	comment, err := parseComment(r)

	if err != nil {
		return nil, err
	}

	if comment.ID != u.ID {
		return nil, &utils.Error_no_access
	}

	var request PutCommentRequest
	e := utils.ValidateRequestBody(r, &request)

	if e != nil {
		return nil, e
	}

	comment.Message = request.Message
	comment.Save()

	return nil, nil
}

/*
Delete a comment
*/
func deleteComment(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	comment, err := parseComment(r)

	if err != nil {
		return nil, err
	}

	if comment.ID != u.ID {
		return nil, &utils.Error_no_access
	}

	comment.Delete()

	return nil, nil
}

func parseComment(r *http.Request) (*db.Comment, *utils.ErrorResponse) {
	commentId, err := strconv.ParseUint(mux.Vars(r)["comment"], 10, 32)

	if err != nil {
		return nil, &utils.Error_comment_not_found
	}

	comment := db.GetCommentById(uint(commentId))

	if comment == nil {
		return nil, &utils.Error_comment_not_found
	}

	return comment, nil
}
