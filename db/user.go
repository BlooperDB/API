package db

import (
	"time"

	"fmt"

	"github.com/BlooperDB/API/utils"
	"github.com/gocql/gocql"
	"github.com/wuman/firebase-server-sdk-go"
)

var UserTable = "user"

type User struct {
	Id           string
	Email        string
	Username     *string
	Avatar       string
	BlooperToken string
	RegisterDate int64
}

func (m User) Save() {
	i := GetSession().Query("UPDATE "+UserTable+" SET "+
		" avatar=?,"+
		" register_date=?"+
		" WHERE id=? AND email=? AND username=? AND blooper_token=?;", m.Avatar, m.RegisterDate, m.Id, m.Email, m.Username, m.BlooperToken).Exec()
	fmt.Println(i)
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
			Username:     &name,
			Avatar:       avatar,
			BlooperToken: GenerateBlooperToken(),
			RegisterDate: time.Now().Unix(),
		}

		user.Save()

		return &user, true
	}

	username := data["username"].(string)

	return &User{
		Id:           data["id"].(string),
		Email:        data["email"].(string),
		Username:     &username,
		Avatar:       data["avatar"].(string),
		BlooperToken: data["blooper_token"].(string),
		RegisterDate: data["register_date"].(int64),
	}, false
}

func GenerateBlooperToken() string {
	return utils.GenerateRandomString(32)
}
