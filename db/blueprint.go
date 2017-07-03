package db

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Blueprint struct {
	gorm.Model

	UserID       uint   `gorm:"index;not null"`
	Name         string `gorm:"not null"`
	Description  string `gorm:"not null"`
	LastRevision uint   `gorm:"not null"`
}

func SearchBlueprints(query string, offset int, limit int) []*Blueprint {
	query = strings.ToLower(query)
	split := strings.Split(query, " ")
	joined := "(" + strings.Join(split, "|") + ")%"
	fullJoined := "%" + joined

	var blueprints []*Blueprint
	db.Raw(`
		SELECT *
		FROM blueprints b
		WHERE id IN (
			SELECT blueprint_id
			FROM blueprint_tags
			WHERE tag_id IN (
				SELECT id
				FROM tags
				WHERE LOWER("name") SIMILAR TO ?
			)
		)
		OR id IN (
			SELECT blueprint_id
			FROM revisions
			WHERE LOWER("changes") SIMILAR TO ?
		)
		OR LOWER("name") SIMILAR TO ?
		OR LOWER("description") SIMILAR TO ?
		ORDER BY (
			select (
				SUM(case when thumbs_up = true then 1 else 0 end) 
				- 
				SUM(case when thumbs_up = false then 1 else 0 end)
			)
			from ratings
			where revision_id = (
				select id
				from revisions
				where blueprint_id = b.id
				limit 1
			)
		), id DESC
		OFFSET ?
		LIMIT ?
	`, joined, fullJoined, fullJoined, fullJoined, offset, limit).Scan(&blueprints)
	return blueprints
}

func GetAllBlueprints(offset int, limit int) []*Blueprint {
	var blueprints []*Blueprint
	db.Offset(offset).Limit(limit).Find(&blueprints)
	return blueprints
}

func GetBlueprintById(id uint) *Blueprint {
	var blueprint Blueprint
	db.Where("id = ?", id).Find(&blueprint)
	if blueprint.ID != 0 {
		return &blueprint
	}
	return nil
}

type blueprintLatestRevision struct {
	BlueprintId uint
	Revision    uint
}

func GetLatestBlueprintRevisions(ids ...uint) map[uint]uint {
	var revs []blueprintLatestRevision
	q := db.Table("revisions").Select("blueprint_id, revision")
	if len(ids) > 0 {
		q = q.Where("blueprint_id IN (?)", ids)
	}
	q.Where(`
		revision = (
			SELECT revision
			FROM revisions
			WHERE deleted_at IS NULL
			ORDER BY revision DESC
			LIMIT 1
		)
	`).Scan(&revs)

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
	db.Where("id = ?", m.UserID).Find(&user)
	return user
}

func (m Blueprint) GetTags() []*Tag {
	var tags []*Tag
	db.Raw(`
		SELECT t.*
		FROM tags t
		JOIN blueprint_tags bt ON (t.id = bt.tag_id)
		WHERE bt.blueprint_id = ?`, m.ID).Scan(&tags)
	return tags
}

func (m Blueprint) GetTag(tag uint) *BlueprintTag {
	var blueprintTag BlueprintTag
	db.Where("blueprint_id = ? AND tag_id = ?", m.ID, tag).Find(&blueprintTag)
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

func (m Blueprint) GetRevisions() []*Revision {
	var revisions []*Revision
	db.Where("blueprint_id = ?", m.ID).Find(&revisions)
	return revisions
}

func (m *Blueprint) CountRevisions() uint {
	var count uint
	db.Table("revisions").Where("blueprint_id = ?", m.ID).Count(&count)
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

func PopularBlueprints(offset int, limit int) []*Blueprint {
	var blueprints []*Blueprint
	db.Raw(`
		SELECT *
		FROM blueprints b
		ORDER BY
			(
				SELECT inside.hotness
				FROM (
					SELECT
					round(
						CAST(
							(
								(CASE WHEN scoring.score > 0 THEN 1 WHEN 100 = 0 THEN 0 ELSE -1 END)
								*
								log(10, greatest(abs(scoring.score), 1))
								+
								((EXTRACT(epoch FROM b.created_at) - 1400000000) / 45000)
							) AS numeric
						),
						7
					) AS hotness
					FROM (
						SELECT
						(
							SELECT (
								SUM(CASE WHEN thumbs_up = true THEN 1 ELSE 0 END)
								-
								SUM(CASE WHEN thumbs_up = false THEN 1 ELSE 0 END)
							)
							FROM ratings
							WHERE revision_id = (
								SELECT id
								FROM revisions
								WHERE blueprint_id = b.id
								LIMIT 1
							)
						) AS score
					) AS scoring
				) AS inside
			) DESC,
			id DESC
		OFFSET ?
		LIMIT ?
	`, offset, limit).Scan(&blueprints)
	return blueprints
}

func TopBlueprints(offset int, limit int) []*Blueprint {
	var blueprints []*Blueprint
	db.Raw(`
		SELECT *
		FROM blueprints b
		ORDER BY (
			select (
				SUM(case when thumbs_up = true then 1 else 0 end)
				-
				SUM(case when thumbs_up = false then 1 else 0 end)
			)
			from ratings
			where revision_id = (
				select id
				from revisions
				where blueprint_id = b.id
				limit 1
			)
		), id DESC
		OFFSET ?
		LIMIT ?
	`, offset, limit).Scan(&blueprints)
	return blueprints
}

func NewBlueprints(offset int, limit int) []*Blueprint {
	var blueprints []*Blueprint
	db.Raw(`
		SELECT *
		FROM blueprints b
		ORDER BY created_at DESC, id DESC
		OFFSET ?
		LIMIT ?
	`, offset, limit).Scan(&blueprints)
	return blueprints
}
