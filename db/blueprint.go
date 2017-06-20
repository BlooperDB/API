package db

import (
	"github.com/jinzhu/gorm"
)

type Blueprint struct {
	gorm.Model

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID uint `gorm:"index;not null"`

	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Versions    []Version
	Tags        []Tag `gorm:"many2many:blueprint_tags;"`
}

func (m Blueprint) Save() {
	db.Save(m)
}

func (m Blueprint) Delete() {
	db.Delete(m)
}

func GetBlueprintById(id string) *Blueprint {
	var blueprint Blueprint
	db.First(&blueprint, id)

	if blueprint.ID == 0 {
		return nil
	}

	return &blueprint
}

func (m Blueprint) GetVersions() []Version {
	var versions []Version
	db.Model(m).Related(versions)
	return versions
}

func (m Blueprint) GetTags() []Tag {
	var tags []Tag
	db.Model(m).Related(tags)
	return tags
}

func GetAllBlueprints() []Blueprint {
	var blueprint []Blueprint
	db.Find(&blueprint)
	return blueprint
}
