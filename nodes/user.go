package nodes

import (
	"net/http"

	"encoding/json"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/wuman/firebase-server-sdk-go"
)

type PublicUserResponse struct {
	Username   string           `json:"username"`
	Avatar     string           `json:"avatar"`
	Blueprints []SmallBlueprint `json:"blueprints"`
}

type SmallBlueprint struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func RegisterUserRoutes(router api.RegisterRoute) {
	router("POST", "/user/signin", signIn)
}

type UserSignInResponse struct {
	BlooperToken string `json:"blooper-token"`
	FirstLogin   bool   `json:"first-login"`
}

type UserSignInRequest struct {
	FirebaseToken string `json:"firebase-token"`
}

func signIn(r *http.Request) (interface{}, *api.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request UserSignInRequest
	err := decoder.Decode(&request)

	if err != nil || request.FirebaseToken == "" {
		return nil, &error_invalid_request_data
	}

	auth, _ := firebase.GetAuth()
	decodedToken, err := auth.VerifyIDToken(request.FirebaseToken)

	if err != nil {
		return nil, &api.ErrorResponse{
			Code:    error_user_token_invalid.Code,
			Message: error_user_token_invalid.Message + ": " + err.Error(),
			Status:  error_user_token_invalid.Status,
		}
	}

	_, found := decodedToken.UID()

	if !found {
		return nil, &error_user_token_invalid
	}

	user, firstLogin := db.SignIn(decodedToken)

	return UserSignInResponse{
		BlooperToken: user.BlooperToken,
		FirstLogin:   firstLogin,
	}, nil
}
