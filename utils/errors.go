package utils

import "github.com/BlooperDB/API/api"

var (
	Error_invalid_request_data = api.ErrorResponse{1, "Invalid request data", 400}
)

var (
	Error_user_token_invalid    = api.ErrorResponse{100, "User token invalid", 400}
	Error_blooper_token_invalid = api.ErrorResponse{101, "Blooper token invalid", 400}
	Error_no_access             = api.ErrorResponse{101, "No access", 403}
)

var (
	Error_blueprint_not_found = api.ErrorResponse{200, "Blueprint not found", 404}
)
