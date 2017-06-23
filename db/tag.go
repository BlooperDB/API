package db

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model

	Name string `gorm:"not null;unique"`
}

type BlueprintTag struct {
	gorm.Model

	BlueprintId uint `gorm:"not null;unique_index:idx_bp_tag"`
	TagId       uint `gorm:"not null;unique_index:idx_bp_tag"`
}

func GetTagById(id uint) *Tag {
	var tag Tag
	db.Where("id = ?", id).
		Find(&tag)
	if tag.ID != 0 {
		return &tag
	}
	return nil
}

func GetTagByName(name string) *Tag {
	var tag Tag
	db.Where("name = ?", name).
		Find(&tag)
	if tag.ID != 0 {
		return &tag
	}
	return nil
}

func (m *Tag) Save() {
	db.Save(m)
}

func (m *Tag) Delete() {
	db.Delete(m)
}

func (m *Tag) GetBlueprints() []Blueprint {
	var blueprints []Blueprint
	db.Raw(`
		SELECT b.*
		FROM blueprint_tags bt
		JOIN blueprints b ON (b.id = bt.blueprint_id)
		WHERE bt.tag_id = ?
	`, m.ID).Scan(&blueprints)
	return blueprints
}

func (m *BlueprintTag) Save() {
	db.Save(m)
}

func (m *BlueprintTag) Delete() {
	db.Delete(m)
}

func (m *BlueprintTag) GetBlueprint() *Blueprint {
	return GetBlueprintById(m.BlueprintId)
}

func (m *BlueprintTag) GetTag() *Tag {
	return GetTagById(m.TagId)
}
