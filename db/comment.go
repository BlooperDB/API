package db

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model

	RevisionID uint     `gorm:"index; not null"`
	UserID     uint     `gorm:"index; not null"`
	Message    string   `gorm:"not null"`
}

func (m *Comment) Save() {
	db.Save(m)
}

func (m *Comment) Delete() {
	db.Delete(m)
}

func GetCommentById(id uint) *Comment {
	var comment Comment
	db.Where("id = ?", id).Find(&comment)
	if comment.ID != 0 {
		return &comment
	}
	return nil
}

func (m Comment) GetUser() User {
	var user User
	db.Where("id = ?", m.UserID).Find(&user)
	return user
}

func (m Comment) GetRevision() Revision {
	var revision Revision
	db.Where("id = ?", m.RevisionID).Find(&revision)
	return revision
}
