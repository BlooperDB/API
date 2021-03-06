package nodes

import (
	"net/http"

	"strconv"

	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
	"github.com/wuman/firebase-server-sdk-go"
)

type PrivateUserResponse struct {
	Id         uint                 `json:"id"`
	Email      string               `json:"email"`
	Username   string               `json:"username"`
	Avatar     string               `json:"avatar"`
	CreatedAt  time.Time            `json:"register-date"`
	UpdatedAt  time.Time            `json:"register-date"`
	Blueprints []*BlueprintResponse `json:"blueprints,omitempty"`
}

type PublicUserResponse struct {
	Id         uint                 `json:"id"`
	Username   string               `json:"username"`
	Avatar     string               `json:"avatar"`
	Blueprints []*BlueprintResponse `json:"blueprints,omitempty"`
}

func RegisterUserRoutes(router api.RegisterRoute) {
	router("POST", "/user/signin", signIn)

	router("GET", "/user/self", api.AuthHandler(getUserSelf, false))
	router("PUT", "/user/self", api.AuthHandler(putUserSelf, false))
	router("GET", "/user/self/blueprints", api.AuthHandler(getUserSelfBlueprints, false))

	router("GET", "/user/{user}", getUser)
	router("GET", "/user/{user}/blueprints", getUserBlueprints)
}

type UserSignInResponse struct {
	BlooperToken string `json:"blooper-token"`
	FirstLogin   bool   `json:"first-login"`
}

type UserSignInRequest struct {
	FirebaseToken string `json:"firebase-token" validate:"nonzero"`
}

func signIn(r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request UserSignInRequest
	e := utils.ValidateRequestBody(r, &request)

	if e != nil {
		return nil, e
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
	userId, err := strconv.ParseUint(mux.Vars(r)["user"], 10, 32)
	if err != nil {
		return nil, &utils.Error_user_not_found
	}

	user := db.GetUserById(uint(userId))

	if user == nil {
		return nil, &utils.Error_user_not_found
	}

	getBlueprints := len(r.URL.Query()["blueprints"]) > 0

	var reBlueprint []*BlueprintResponse

	if getBlueprints {
		blueprints := user.GetUserBlueprints()
		reBlueprint = reBlueprintData(blueprints)
	}

	authUser := db.GetAuthUser(r)

	if authUser != nil && authUser.ID == uint(userId) {
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

	return PublicUserResponse{
		Id:         uint(userId),
		Username:   user.Username,
		Avatar:     user.Avatar,
		Blueprints: reBlueprint,
	}, nil
}

func getUserSelf(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	blueprints := u.GetUserBlueprints()

	getBlueprints := len(r.URL.Query()["blueprints"]) > 0

	var reBlueprint []*BlueprintResponse

	if getBlueprints {
		reBlueprint = reBlueprintData(blueprints)
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
	Username string `json:"username"`
}

func putUserSelf(u *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	var request PutUserRequest

	if e := utils.ValidateRequestBody(r, &request); e != nil {
		return nil, e
	}
	uname := request.Username
	if !utils.UsernameRegex.MatchString(uname) {
		return nil, &utils.Error_invalid_username
	}

	if existingUser := db.GetUserByUsername(request.Username); existingUser != nil {
		if u.ID != existingUser.ID {
			return nil, &utils.Error_username_taken
		} else {
			return nil, &utils.Error_nothing_changed
		}
	}

	u.Username = uname
	u.Save()

	return nil, nil
}

type UserBlueprintResponse struct {
	Blueprints []*BlueprintResponse `json:"blueprints"`
}

func getUserBlueprints(r *http.Request) (interface{}, *utils.ErrorResponse) {
	userId, err := strconv.ParseUint(mux.Vars(r)["user"], 10, 32)
	if err != nil {
		return nil, &utils.Error_user_not_found
	}

	user := db.GetUserById(uint(userId))
	if user == nil {
		return nil, &utils.Error_user_not_found
	}
	return getBlueprintsUser(user)
}

func getUserSelfBlueprints(user *db.User, r *http.Request) (interface{}, *utils.ErrorResponse) {
	return getBlueprintsUser(user)
}

func getBlueprintsUser(user *db.User) (interface{}, *utils.ErrorResponse) {
	blueprints := user.GetUserBlueprints()
	reBlueprint := reBlueprintData(blueprints)

	return UserBlueprintResponse{
		Blueprints: reBlueprint,
	}, nil
}
