package db

import (
	"github.com/jinzhu/gorm"
)

type Rating struct {
	gorm.Model

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID uint `gorm:"index; not null"`

	Revision   Revision `gorm:"ForeignKey:RevisionID;AssociationForeignKey:ID"`
	RevisionID uint     `gorm:"index; not null"`

	ThumbsUp bool `gorm:"not null"`
}

func (m *Rating) Save() {
	db.Save(m)
}

func (m *Rating) Delete() {
	db.Delete(m)
}

func (m Rating) GetUser() User {
	var user User
	db.Model(m).Related(&user)
	return user
}

func (m Rating) GetRevision() Revision {
	var revision Revision
	db.Model(m).Related(&revision)
	return revision
}
