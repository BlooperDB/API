package db

import (
	"github.com/jinzhu/gorm"
)

type Blueprint struct {
	gorm.Model

	User   User `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	UserID uint `gorm:"index;not null"`

	Name         string `gorm:"not null"`
	Description  string `gorm:"not null"`
	LastRevision uint   `gorm:"not null"`

	Revisions []Revision
	Tags      []Tag `gorm:"many2many:blueprint_tags;"`
}

func (m *Blueprint) Save() {
	db.Save(m)
}

func (m *Blueprint) Delete() {
	db.Delete(m)
}

func GetBlueprintById(id uint) *Blueprint {
	var blueprint Blueprint
	db.First(&blueprint, id)

	if blueprint.ID == 0 {
		return nil
	}

	return &blueprint
}

func (m Blueprint) GetTags() []Tag {
	var tags []Tag
	db.Model(m).Related(&tags)
	return tags
}

func GetAllBlueprints() []Blueprint {
	var blueprint []Blueprint
	db.Find(&blueprint)
	return blueprint
}

func (m Blueprint) IncrementAndGetRevision() uint {
	m.LastRevision++
	i := m.LastRevision
	m.Save()
	return i
}

func (m Blueprint) GetRevisions() []Revision {
	var revisions []Revision
	db.Model(m).Related(&revisions)
	return revisions
}

func (m *Blueprint) CountRevisions() uint {
	var count uint
	db.Table("revisions").
		Where("blueprint_id = ?", m.ID).
		Count(&count)
	return count
}

func (m Blueprint) GetRevision(id uint) *Revision {
	var revisions []Revision
	db.Where("blueprint_id = ? AND revision = ?", m.ID, id).
		Limit(1).Find(&revisions)
	if len(revisions) > 0 {
		return &revisions[0]
	}
	return nil
}

func (m *Blueprint) GetLatestRevision() *Revision {
	rev := m.GetRevision(m.LastRevision)
	if rev != nil {
		return rev
	}
	var revisions []Revision
	db.Where("blueprint_id = ?", m.ID).
		Order("revision desc").Limit(1).
		Find(&revisions)
	if len(revisions) > 0 {
		return &revisions[0]
	}
	return nil
}