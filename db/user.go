package db

import (
	"time"

	"github.com/BlooperDB/API/utils"
	"github.com/gocql/gocql"
	"github.com/wuman/firebase-server-sdk-go"
)

var UserTable = [2]string{
	"user",
	"CREATE TABLE IF NOT EXISTS user (" +
		"id varchar PRIMARY KEY," +
		"email varchar," +
		"username varchar," +
		"blooper_token varchar," +
		"register_date int" +
		");",
}

type User struct {
	Id           string
	Email        string
	Username     *string
	Avatar       string
	BlooperToken string
	RegisterDate int64
}

func (m User) Save() {
	GetSession().Query("UPDATE "+UserTable[0]+" SET "+
		" id=?,"+
		" email=?,"+
		" username=?,"+
		" avatar=?,"+
		" register_date=?,"+
		" WHERE id=?;", m.Id, m.Email, m.Username, m.BlooperToken, m.RegisterDate, m.Id)
}

func SignIn(token *firebase.Token) (*User, bool) {
	data := make(map[string]interface{})
	GetSession().Query("SELECT * FROM "+UserTable[0]+" WHERE email = ?;", token.Email).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		name, _ := token.Name()
		email, _ := token.Email()
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
