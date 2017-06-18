package nodes

import (
	"net/http"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/gorilla/mux"
	"github.com/wuman/firebase-server-sdk-go"
)

var (
	error_user_token_invalid = api.ErrorResponse{100, "User token invalid", 400}
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

func signIn(r *http.Request) (interface{}, *api.ErrorResponse) {
	firebase_token := mux.Vars(r)["firebase_token"]

	auth, _ := firebase.GetAuth()
	decodedToken, err := auth.VerifyIDToken(firebase_token)

	if err != nil {
		return nil, &error_user_token_invalid
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
