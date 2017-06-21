package utils

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

var (
	Error_invalid_request_data = ErrorResponse{1, "Invalid request data", 400}
)

var (
	Error_user_token_invalid    = ErrorResponse{100, "User token invalid", 400}
	Error_blooper_token_invalid = ErrorResponse{101, "Blooper token invalid", 400}
	Error_no_access             = ErrorResponse{101, "No access", 403}
)

var (
	Error_blueprint_not_found = ErrorResponse{200, "Blueprint not found", 404}
)

var (
	Error_revision_not_found = ErrorResponse{300, "Blueprint not found", 404}
)

var (
	Error_comment_not_found = ErrorResponse{400, "Blueprint not found", 404}
)
