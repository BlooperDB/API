package db

import (
	"github.com/jinzhu/gorm"
)

type Rating struct {
	gorm.Model

	UserID     uint `gorm:"index;not null;unique_index:idx_uid_rev"`
	RevisionID uint `gorm:"index;not null;unique_index:idx_uid_rev"`
	ThumbsUp   bool `gorm:"not null"`
}

func (m *Rating) Save() {
	db.Save(m)
}

func (m *Rating) Delete() {
	db.Delete(m)
}

func (m Rating) GetUser() User {
	var user User
	db.Where("id = ?", m.UserID).Find(&user)
	return user
}

func (m Rating) GetRevision() Revision {
	var revision Revision
	db.Where("id = ?", m.RevisionID).Find(&revision)
	return revision
}
