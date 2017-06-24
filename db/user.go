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

func GetAuthUser(r *http.Request) *User {
	return GetUserByBlooperToken(r.Header.Get("BLOOPER-TOKEN"))
}

func GetUserById(id uint) *User {
	var user User
	db.Where("id = ?", id).Find(&user)
	if user.ID != 0 {
		return &user
	}
	return nil
}

func GetUserByUsername(username string) *User {
	if username == "" {
		return nil
	}
	var user User
	db.Where("LOWER(username) = LOWER(?)", username).Find(&user)
	if user.ID != 0 {
		return &user
	}
	return nil
}

func GetUserByBlooperToken(token string) *User {
	var user User
	db.Where("blooper_token = ?", token).Find(&user)
	if user.ID != 0 {
		return &user
	}
	return nil
}

func (m User) GetUserBlueprints() []Blueprint {
	var blueprints []Blueprint
	db.Where("user_id = ?", m.ID).Find(&blueprints)
	return blueprints
}

func (m User) GetComments() []Comment {
	var comments []Comment
	db.Where("user_id = ?", m.ID).Find(&comments)
	return comments
}

func (m *User) Save() {
	db.Save(m)
}

func (m *User) Delete() {
	db.Delete(m)
}