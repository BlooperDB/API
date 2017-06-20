package db

type Tag struct {
	Name       string      `gorm:"primary_key"`
	Blueprints []Blueprint `gorm:"many2many:blueprint_tags;"`
}

func (m Tag) Save() {
	db.Save(&m)
}

func (m Tag) Delete() {
	db.Delete(&m)
}

func GetTagById(id string) *Tag {
	var tag Tag
	db.First(&tag, id)

	if tag.Name == "" {
		return nil
	}

	return &tag
}
