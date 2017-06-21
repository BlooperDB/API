package nodes

import (
	"net/http"

	"encoding/json"

	"strconv"

	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
	"github.com/wuman/firebase-server-sdk-go"
	"gopkg.in/validator.v2"
)

type PrivateUserResponse struct {
	Id         uint             `json:"id"`
	Email      string           `json:"email"`
	Username   string           `json:"username"`
	Avatar     string           `json:"avatar"`
	CreatedAt  time.Time        `json:"register-date"`
	UpdatedAt  time.Time        `json:"register-date"`
	Blueprints []SmallBlueprint `json:"blueprints"`
}

type PublicUserResponse struct {
	Id         uint             `json:"id"`
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
	router("GET", "/user/self", api.AuthHandler(getUserSelf))
	router("GET", "/user/self/blueprints", api.AuthHandler(getUserSelfBlueprints))

	router("GET", "/user/{user}", getUser)
	router("PUT", "/user/{user}", api.AuthHandler(putUser))
	router("GET", "/user/{user}/blueprints", getUserBlueprints)
}

type UserSignInResponse struct {
	BlooperToken string `json:"blooper-token"`
	FirstLogin   bool   `json:"first-login"`
}

type UserSignInRequest struct {
	FirebaseToken string `json:"firebase-token";validate:"nonzero"`
}

func signIn(r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request UserSignInRequest
	err := decoder.Decode(&request)

	if err != nil || request.FirebaseToken == "" {
		return nil, &utils.Error_invalid_request_data
	}

	if err = validator.Validate(request); err != nil {
		return nil, &utils.ErrorResponse{
			Code:    utils.Error_invalid_request_data.Code,
			Message: utils.Error_invalid_request_data.Message + ": " + err.Error(),
			Status:  utils.Error_invalid_request_data.Status,
		}
	}

	auth, _ := firebase.GetAuth()
	decodedToken, err := auth.VerifyIDToken(request.FirebaseToken)

	if err != nil {
		return nil, &utils.ErrorResponse{
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

func getUser(r *http.Request) (interface{}, *utils.ErrorResponse) {
	userId, _ := strconv.ParseUint(mux.Vars(r)["user"], 10, 32)

	user := db.GetUserById(uint(userId))

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

	authUser := db.GetAuthUser(r)

	if authUser != nil {
		if authUser.ID == uint(userId) {
			return PrivateUserResponse{
				Id:         uint(userId),
				Email:      user.Email,
				Username:   user.Username,
				Avatar:     user.Avatar,
				CreatedAt:  user.CreatedAt,
				UpdatedAt:  user.UpdatedAt,
				Blueprints: reBlueprint,
			}, nil
		}
	}

	return PublicUserResponse{
		Id:         uint(userId),
		Username:   user.Username,
		Avatar:     user.Avatar,
		Blueprints: reBlueprint,
	}, nil
}

func getUserSelf(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprints := u.GetUserBlueprints()
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

	return PrivateUserResponse{
		Id:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		Avatar:     u.Avatar,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		Blueprints: reBlueprint,
	}, nil
}

type PutUserRequest struct {
	Username string `json:"username";validate:"nonzero"`
}

func putUser(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	decoder := json.NewDecoder(r.Body)
	var request PutUserRequest
	err := decoder.Decode(&request)

	if err != nil {
		return nil, &utils.Error_invalid_request_data
	}

	if err = validator.Validate(request); err != nil {
		return nil, &utils.ErrorResponse{
			Code:    utils.Error_invalid_request_data.Code,
			Message: utils.Error_invalid_request_data.Message + ": " + err.Error(),
			Status:  utils.Error_invalid_request_data.Status,
		}
	}

	userId, _ := strconv.ParseUint(mux.Vars(r)["user"], 10, 32)

	if u.ID != uint(userId) {
		return nil, &utils.Error_no_access
	}

	u.Username = request.Username

	u.Save()

	return nil, nil
}

type UserBlueprintResponse struct {
	Blueprints []UserBlueprintResponseBlueprint `json:"blueprints"`
}

type UserBlueprintResponseBlueprint struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created-at"`
	UpdatedAt   time.Time `json:"updated-at"`
}

func getUserBlueprints(r *http.Request) (interface{}, *utils.ErrorResponse) {
	userId, _ := strconv.ParseUint(mux.Vars(r)["user"], 10, 32)

	user := db.GetUserById(uint(userId))

	blueprints := user.GetUserBlueprints()
	reBlueprint := make([]UserBlueprintResponseBlueprint, len(blueprints))

	for i := 0; i < len(blueprints); i++ {
		blueprint := blueprints[i]

		tags := blueprint.GetTags()
		reTags := make([]string, len(tags))

		for i := 0; i < len(tags); i++ {
			tag := tags[i]
			reTags[i] = tag.Name
		}

		reBlueprint[i] = UserBlueprintResponseBlueprint{
			Id:          blueprint.ID,
			Name:        blueprint.Name,
			Description: blueprint.Description,
			Tags:        reTags,
			CreatedAt:   blueprint.CreatedAt,
			UpdatedAt:   blueprint.UpdatedAt,
		}
	}

	return UserBlueprintResponse{
		Blueprints: reBlueprint,
	}, nil
}

func getUserSelfBlueprints(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprints := u.GetUserBlueprints()
	reBlueprint := make([]UserBlueprintResponseBlueprint, len(blueprints))

	for i := 0; i < len(blueprints); i++ {
		blueprint := blueprints[i]

		tags := blueprint.GetTags()
		reTags := make([]string, len(tags))

		for i := 0; i < len(tags); i++ {
			tag := tags[i]
			reTags[i] = tag.Name
		}

		reBlueprint[i] = UserBlueprintResponseBlueprint{
			Id:          blueprint.ID,
			Name:        blueprint.Name,
			Description: blueprint.Description,
			Tags:        reTags,
			CreatedAt:   blueprint.CreatedAt,
			UpdatedAt:   blueprint.UpdatedAt,
		}
	}

	return UserBlueprintResponse{
		Blueprints: reBlueprint,
	}, nil
}
