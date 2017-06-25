package db

import (
	"github.com/jinzhu/gorm"
)

type Revision struct {
	gorm.Model

	BlueprintID      uint   `gorm:"index;not null;unique_index:idx_bp_rev"`
	Revision         uint   `gorm:"not null;unique_index:idx_bp_rev"`
	Changes          string `gorm:"not null"`
	BlueprintString  string `gorm:"not null"`
	BlueprintVersion int    `gorm:"not null"`
}

func (m *Revision) Save() {
	db.Save(m)
}

func (m *Revision) Delete() {
	db.Delete(m)
}

func (m Revision) GetComments() []Comment {
	var comments []Comment
	db.Where("revision_id = ?", m.ID).Find(&comments)
	return comments
}

func (m Revision) GetRatings() []Rating {
	var ratings []Rating
	db.Where("revision_id = ?", m.ID).Find(&ratings)
	return ratings
}

func (m Revision) GetBlueprint() Blueprint {
	var blueprint Blueprint
	db.Where("id = ?", m.BlueprintID).Find(&blueprint)
	return blueprint
}

func GetRevisionById(id uint) *Revision {
	var revision Revision
	db.Where("id = ?", id).Find(&revision)
	if revision.ID != 0 {
		return &revision
	}
	return nil
}
