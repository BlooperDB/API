package db

import (
	"github.com/jinzhu/gorm"
)

type Rating struct {
	gorm.Model

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID uint `gorm:"index; not null"`

	Version   Version `gorm:"ForeignKey:VersionID;AssociationForeignKey:ID"`
	VersionID uint    `gorm:"index; not null"`

	ThumbsUp bool `gorm:"not null"`
}

func (m Rating) Save() {
	db.Save(m)
}

func (m Rating) Delete() {
	db.Delete(m)
}

func (m Rating) GetUser() User {
	var user User
	db.Model(m).Related(user)
	return user
}

func (m Rating) GetVersion() Version {
	var version Version
	db.Model(m).Related(version)
	return version
}
