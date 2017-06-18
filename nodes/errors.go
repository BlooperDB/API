package nodes

import "github.com/BlooperDB/API/api"

var (
	error_invalid_request_data = api.ErrorResponse{1, "Invalid request data", 400}
)

var (
	error_user_token_invalid = api.ErrorResponse{100, "User token invalid", 400}
)

var (
	error_blueprint_not_found = api.ErrorResponse{200, "Blueprint not found", 404}
)
