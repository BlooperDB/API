package db

import (
	"net/http"

	"github.com/BlooperDB/API/utils"
	"github.com/jinzhu/gorm"
	"github.com/wuman/firebase-server-sdk-go"
)

type User struct {
	gorm.Model

	Email        string `gorm:"unique_index;not null"`
	Username     string `gorm:"unique_index;not null"`
	Avatar       string `gorm:"not null"`
	BlooperToken string `gorm:"unique_index;not null"`
	Blueprints   []Blueprint
	Comments     []Comment
}

func SignIn(token *firebase.Token) (User, bool) {
	email, _ := token.Email()

	var user User
	db.Where("email = ?", email).Find(&user)

	if user.ID == 0 {
		avatar, _ := token.Picture()

		user = User{
			Email:        email,
			Username:     "",
			Avatar:       avatar,
			BlooperToken: GenerateBlooperToken(),
		}

		user.Save()
		return user, true
	}

	return user, false
}

func GenerateBlooperToken() string {
	return utils.GenerateRandomString(32)
}

func GetUserById(userId uint) *User {
	var user User
	db.First(&user, userId)
	if user.ID == 0 {
		return nil
	}
	return &user
}

func GetAuthUser(r *http.Request) *User {
	return GetUserByBlooperToken(r.Header.Get("BLOOPER-TOKEN"))
}

func GetUserByBlooperToken(token string) *User {
	var user User
	db.First(&user, "blooper_token = ?", token)
	if user.ID == 0 {
		return nil
	}
	return &user
}

func (m *User) Save() {
	db.Save(m)
}

func (m *User) Delete() {
	db.Delete(m)
}

func (m User) GetUserBlueprints() []Blueprint {
	var blueprints []Blueprint
	db.Model(m).Related(&blueprints)
	return blueprints
}

func (m User) GetComments() []Comment {
	var comments []Comment
	db.Model(m).Related(&comments)
	return comments
}
