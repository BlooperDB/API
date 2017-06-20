package db

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model

	Version   Version `gorm:"ForeignKey:VersionID;AssociationForeignKey:ID"`
	VersionID uint    `gorm:"index; not null"`

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID uint `gorm:"index; not null"`

	Message string `gorm:"not null"`
}

func (m Comment) Save() {
	db.Save(&m)
}

func (m Comment) Delete() {
	db.Delete(&m)
}

func GetCommentById(id uint) *Comment {
	var comment Comment
	db.First(&comment, id)

	if comment.ID == 0 {
		return nil
	}

	return &comment
}

func (m Comment) GetUser() User {
	var user User
	db.Model(m).Related(user)
	return user
}

func (m Comment) GetVersion() Version {
	var version Version
	db.Model(m).Related(version)
	return version
}
