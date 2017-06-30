package utils

import "regexp"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

var UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9]{3,30}$")

var (
	Error_invalid_request_data = ErrorResponse{1, "Invalid request data", 400}
	Error_nothing_changed      = ErrorResponse{2, "Nothing changed", 400}
	Error_internal_error       = ErrorResponse{3, "Internal error", 500}
)

var (
	Error_user_token_invalid    = ErrorResponse{100, "User token invalid", 400}
	Error_blooper_token_invalid = ErrorResponse{101, "Blooper token invalid", 400}
	Error_no_access             = ErrorResponse{102, "No access", 403}
	Error_user_not_found        = ErrorResponse{103, "User not found", 404}
	Error_invalid_username      = ErrorResponse{104, "Invalid username", 400}
	Error_username_required     = ErrorResponse{105, "A username is required to do that", 400}
	Error_username_taken        = ErrorResponse{106, "Username taken", 400}
)

var (
	Error_blueprint_not_found = ErrorResponse{200, "Blueprint not found", 404}
)

var (
	Error_revision_not_found = ErrorResponse{300, "Blueprint revision not found", 404}
)

var (
	Error_comment_not_found = ErrorResponse{400, "Blueprint comment not found", 404}
)

var (
	Error_no_search_terms = ErrorResponse{500, "No search terms given", 400}
)

var (
	Error_rating_not_found = ErrorResponse{600, "Rating not found", 404}
)
