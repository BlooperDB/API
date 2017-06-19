package nodes

import (
	"net/http"

	"encoding/json"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
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

	router("Get", "/user/{user}", getUser)
	router("PUT", "/user/{user}", api.AuthHandler(putUser))
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
		return nil, &utils.Error_invalid_request_data
	}

	auth, _ := firebase.GetAuth()
	decodedToken, err := auth.VerifyIDToken(request.FirebaseToken)

	if err != nil {
		return nil, &api.ErrorResponse{
			Code:    utils.Error_user_token_invalid.Code,
			Message: utils.Error_user_token_invalid.Message + ": " + err.Error(),
			Status:  utils.Error_user_token_invalid.Status,
		}
	}

	_, found := decodedToken.UID()

	if !found {
		return nil, &utils.Error_user_token_invalid
	}

	user, firstLogin := db.SignIn(decodedToken)

	return UserSignInResponse{
		BlooperToken: user.BlooperToken,
		FirstLogin:   firstLogin,
	}, nil
}

func getUser(r *http.Request) (interface{}, *api.ErrorResponse) {
	userId := mux.Vars(r)["user"]

	user := db.GetUserById(userId)

	blueprints := user.GetUserBlueprints()
	reBlueprint := make([]SmallBlueprint, len(blueprints))

	for i := 0; i < len(blueprints); i++ {
		tags := blueprints[i].GetTags()
		reTags := make([]string, len(tags))

		for i := 0; i < len(tags); i++ {
			tag := tags[i]
			reTags[i] = tag.Name
		}

		reBlueprint[i] = SmallBlueprint{
			Name:        blueprints[i].Name,
			Description: blueprints[i].Description,
			Tags:        reTags,
		}
	}

	return PublicUserResponse{
		Username:   user.Username,
		Avatar:     user.Avatar,
		Blueprints: reBlueprint,
	}, nil
}

type PostUserRequest struct {
	Username string `json:"username"`
}

func putUser(u *db.User, r *http.Request) (interface{}, *api.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PostUserRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	userId := mux.Vars(r)["user"]

	if u.Id != userId {
		return nil, &utils.Error_no_access
	}

	u.Username = request.Username

	u.Save()

	return nil, nil
}
