package db

import (
	"github.com/gocql/gocql"
)

var BlueprintTable = "blueprint"

type Blueprint struct {
	Id          string
	UserId      string
	Name        string
	Description string
}

func (m Blueprint) Save() {
	GetSession().Query("UPDATE "+BlueprintTable+" SET "+
		" user_id=?,"+
		" name=?,"+
		" description=?"+
		" WHERE id=?;",
		m.UserId, m.Name, m.Description, m.Id).Exec()
}

func GetBlueprintById(id string) *Blueprint {
	var data map[string]interface{} = make(map[string]interface{})

	GetSession().Query("SELECT * FROM "+BlueprintTable+" WHERE id = ?;", id).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		return nil
	}

	return &Blueprint{
		Id:          data["id"].(string),
		UserId:      data["user_id"].(string),
		Name:        data["name"].(string),
		Description: data["description"].(string),
	}
}

func (m Blueprint) GetVersions() []*Version {
	return FindVersionsByBlueprint(m)
}

func (m Blueprint) GetTags() []*Tag {
	tags := FindTagsByBlueprint(m)

	result := make([]*Tag, len(tags))

	for i := 0; i < len(tags); i++ {
		result[i] = tags[i].GetTag()
	}

	return result
}
