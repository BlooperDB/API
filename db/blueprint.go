package db

import (
	"strconv"

	"strings"

	"fmt"

	"github.com/jinzhu/gorm"
)

type Blueprint struct {
	gorm.Model

	UserID uint `gorm:"index;not null"`

	Name         string `gorm:"not null"`
	Description  string `gorm:"not null"`
	LastRevision uint   `gorm:"not null"`
}

func SearchBlueprints(query string, offset int, limit int) []*Blueprint {
	query = strings.ToLower(query)
	split := strings.Split(query, " ")
	joined := "%(" + strings.Join(split, "|") + ")%"

	rows, err := db.Raw("SELECT *"+
		"FROM blueprints "+
		"WHERE id IN "+
		"    ( SELECT blueprint_id "+
		"     FROM blueprint_tags "+
		"     WHERE tag_id IN "+
		"         ( SELECT id "+
		"          FROM tags "+
		"          WHERE LOWER(\"name\") IN( ? ) ) ) "+
		"  OR id IN "+
		"    ( SELECT blueprint_id "+
		"     FROM revisions "+
		"     WHERE LOWER(\"changes\") SIMILAR TO ? ) "+
		"  OR LOWER(\"name\") SIMILAR TO ? "+
		"  OR LOWER(\"description\") SIMILAR TO ? "+
		"OFFSET ? "+
		"LIMIT ? ",
		split,
		joined,
		joined,
		joined,
		offset,
		limit,
	).Rows()

	if err != nil {
		fmt.Println(err)
	}

	var blueprints []*Blueprint
	i := 0
	for rows.Next() {
		var blueprint Blueprint
		db.ScanRows(rows, &blueprint)
		blueprints = append(blueprints, &blueprint)
		i++
	}

	rows.Close()

	return blueprints
}

func GetAllBlueprints(offset int, limit int) []*Blueprint {
	var blueprints []*Blueprint
	db.Offset(offset).Limit(limit).Find(&blueprints)
	db.Offset(-1).Limit(-1)
	return blueprints
}

func GetBlueprintById(id uint) *Blueprint {
	var blueprint Blueprint
	db.First(&blueprint, id)
	if blueprint.ID == 0 {
		return nil
	}
	return &blueprint
}

type blueprintLatestRevision struct {
	BlueprintId uint
	Revision    uint
}

func GetLatestBlueprintRevisions(ids ...uint) map[uint]uint {
	var revs []blueprintLatestRevision
	q := db.Table("revisions").
		Select("blueprint_id, revision")
	if len(ids) > 0 {
		idString := ""
		for i, id := range ids {
			if i != 0 {
				idString += ", "
			}
			idString += strconv.FormatUint(uint64(id), 10)
		}
		q = q.Where("blueprint_id IN (" + idString + ")")
	}
	q.Where("revision = (SELECT revision FROM revisions WHERE deleted_at IS NULL ORDER BY revision desc LIMIT 1)").
		Scan(&revs)

	ret := make(map[uint]uint)
	for _, r := range revs {
		ret[r.BlueprintId] = r.Revision
	}
	return ret
}

func (m *Blueprint) Save() {
	db.Save(m)
}

func (m *Blueprint) Delete() {
	db.Delete(m)
}

func (m Blueprint) GetAuthor() User {
	var user User
	db.Model(m).Related(&user)
	return user
}

func (m Blueprint) GetTags() []Tag {
	var tags []Tag
	db.Raw(`
		SELECT t.*
		FROM tags t
		JOIN blueprint_tags bt ON (t.id = bt.tag_id)
		WHERE bt.blueprint_id = ?`, m.ID).Scan(&tags)
	return tags
}

func (m Blueprint) GetTag(tag uint) *BlueprintTag {
	var blueprintTag BlueprintTag
	println("a")
	db.Where("blueprint_id = ? AND tag_id = ?", m.ID, tag).Find(&blueprintTag)
	println("a")
	if blueprintTag.ID != 0 {
		return &blueprintTag
	}
	return nil
}

func (m Blueprint) AddTag(tag uint) bool {
	if m.GetTag(tag) != nil {
		return false
	}
	bt := BlueprintTag{
		BlueprintId: m.ID,
		TagId:       tag,
	}
	bt.Save()
	return true
}

func (m Blueprint) RemoveTag(tag uint) bool {
	if t := m.GetTag(tag); t != nil {
		t.Delete()
		return true
	}
	return false
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
