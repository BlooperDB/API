package db

import (
	"time"

	"net/http"

	"github.com/BlooperDB/API/utils"
	"github.com/gocql/gocql"
	"github.com/wuman/firebase-server-sdk-go"
)

var UserTable = "user"

type User struct {
	Id           string
	Email        string
	Username     string
	Avatar       string
	BlooperToken string
	RegisterDate int64
	LastAction   int64
}

func (m User) Save() {
	GetSession().Query("UPDATE "+UserTable+" SET "+
		" avatar=?,"+
		" register_date=?,"+
		" last_action=?"+
		" WHERE id=? AND email=? AND username=? AND blooper_token=?;", m.Avatar, m.RegisterDate, m.LastAction, m.Id, m.Email, m.Username, m.BlooperToken).Exec()
}

func SignIn(token *firebase.Token) (*User, bool) {
	email, _ := token.Email()

	data := make(map[string]interface{})
	GetSession().Query("SELECT * FROM "+UserTable+" WHERE email = ? ALLOW FILTERING;", email).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		name, _ := token.Name()
		avatar, _ := token.Picture()

		user := User{
			Id:           utils.GenerateRandomString(8),
			Email:        email,
			Username:     name,
			Avatar:       avatar,
			BlooperToken: GenerateBlooperToken(),
			RegisterDate: time.Now().Unix(),
			LastAction:   time.Now().Unix(),
		}

		user.Save()

		return &user, true
	}

	return &User{
		Id:           data["id"].(string),
		Email:        data["email"].(string),
		Username:     data["username"].(string),
		Avatar:       data["avatar"].(string),
		BlooperToken: data["blooper_token"].(string),
		RegisterDate: data["register_date"].(int64),
		LastAction:   data["last_action"].(int64),
	}, false
}

func GenerateBlooperToken() string {
	return utils.GenerateRandomString(32)
}

func GetUserByBlooperToken(token string) *User {
	data := make(map[string]interface{})
	GetSession().Query("SELECT * FROM "+UserTable+" WHERE blooper_token = ? ALLOW FILTERING;", token).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		return nil
	}

	return &User{
		Id:           data["id"].(string),
		Email:        data["email"].(string),
		Username:     data["username"].(string),
		Avatar:       data["avatar"].(string),
		BlooperToken: data["blooper_token"].(string),
		RegisterDate: data["register_date"].(int64),
		LastAction:   data["last_action"].(int64),
	}
}

func GetUserById(userId string) *User {
	data := make(map[string]interface{})
	GetSession().Query("SELECT * FROM "+UserTable+" WHERE id = ? ALLOW FILTERING;", userId).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		return nil
	}

	return &User{
		Id:           data["id"].(string),
		Email:        data["email"].(string),
		Username:     data["username"].(string),
		Avatar:       data["avatar"].(string),
		BlooperToken: data["blooper_token"].(string),
		RegisterDate: data["register_date"].(int64),
		LastAction:   data["last_action"].(int64),
	}
}

func (m User) GetUserBlueprints() []*Blueprint {
	return GetBlueprintsByUserId(m.Id)
}

func GetAuthUser(r *http.Request) *User {
	return GetUserByBlooperToken(r.Header.Get("BLOOPER-TOKEN"))
}
