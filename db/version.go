package db

import (
	"github.com/jinzhu/gorm"
)

type Version struct {
	gorm.Model

	Blueprint   Blueprint `gorm:"ForeignKey:BlueprintID;AssociationForeignKey:ID"`
	BlueprintID uint      `gorm:"index not null"`

	Version         string `gorm:"not null"`
	Changes         string
	BlueprintString string `gorm:"not null"`
	Comments        []Comment
}

func (m Version) Save() {
	db.Save(m)
}

func (m Version) Delete() {
	db.Delete(m)
}

func (m Version) GetComments() []Comment {
	var comments []Comment
	db.Model(m).Related(comments)
	return comments
}

func (m Version) GetRatings() []Rating {
	var ratings []Rating
	db.Model(m).Related(ratings)
	return ratings
}

func (m Version) GetBlueprint() Blueprint {
	var blueprint Blueprint
	db.Model(m).Related(blueprint)
	return blueprint
}

func (m Version) GetVersion() Version {
	var version Version
	db.Model(m).Related(version)
	return version
}

func GetVersionById(id string) *Version {
	var version Version
	db.First(&version, id)

	if version.ID == 0 {
		return nil
	}

	return &version
}
