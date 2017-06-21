package db

import (
	"github.com/jinzhu/gorm"
)

type Revision struct {
	gorm.Model

	Blueprint   Blueprint `gorm:"ForeignKey:BlueprintID;AssociationForeignKey:ID"`
	BlueprintID uint      `gorm:"index not null"`

	Revision        uint `gorm:"not null"`
	Changes         string
	BlueprintString string `gorm:"not null"`
	Comments        []Comment
}

func (m Revision) Save() {
	db.Save(&m)
}

func (m Revision) Delete() {
	db.Delete(&m)
}

func (m Revision) GetComments() []Comment {
	var comments []Comment
	db.Model(m).Related(&comments)
	return comments
}

func (m Revision) GetRatings() []Rating {
	var ratings []Rating
	db.Model(m).Related(&ratings)
	return ratings
}

func (m Revision) GetBlueprint() Blueprint {
	var blueprint Blueprint
	db.Model(m).Related(&blueprint)
	return blueprint
}

func (m Revision) GetRevision() Revision {
	var revision Revision
	db.Model(m).Related(&revision)
	return revision
}

func GetRevisionById(id uint) *Revision {
	var revision Revision
	db.First(&revision, id)

	if revision.ID == 0 {
		return nil
	}

	return &revision
}
