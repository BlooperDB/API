package db

import (
	"github.com/gocql/gocql"
)

var BlueprintTagTable = "blueprint_tag"

type BlueprintToTag struct {
	BlueprintId string
	TagId       string
}

func (m BlueprintToTag) Save() {
	GetSession().Query("UPDATE "+BlueprintTagTable+" SET "+
		" blueprint_id=?,"+
		" tag_id=?"+
		" WHERE blueprint_id=? AND tag_id=? ;",
		m.BlueprintId, m.TagId, m.BlueprintId, m.TagId).Exec()
}

func FindTagsByBlueprint(b Blueprint) []*BlueprintToTag {
	r := GetSession().Query("SELECT * FROM "+BlueprintTagTable+" WHERE blueprint_id = ?;", b.Id).Consistency(gocql.All).Iter()

	result := make([]*BlueprintToTag, r.NumRows())

	for i := 0; i < r.NumRows(); i++ {
		data := make(map[string]interface{})
		r.MapScan(data)
		result[i] = &BlueprintToTag{
			BlueprintId: data["blueprint_id"].(string),
			TagId:       data["tag_id"].(string),
		}
	}

	return result
}

func (m BlueprintToTag) GetTag() *Tag {
	return GetTagById(m.TagId)
}
