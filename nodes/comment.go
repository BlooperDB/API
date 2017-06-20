package nodes

import (
	"net/http"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
)

func RegisterCommentRoutes(router api.RegisterRoute) {
	router("GET", "/comment/{comment}", getComment)
	router("POST", "/comment/{comment}", api.AuthHandler(postComment))
	router("PUT", "/comment/{comment}", api.AuthHandler(updateComment))
	router("DELETE", "/comment/{comment}", api.AuthHandler(deleteComment))
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
