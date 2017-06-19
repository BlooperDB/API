package db

import (
	"github.com/gocql/gocql"
)

var TagTable = "tag"

type Tag struct {
	Name string
}

func (m Tag) Save() {
	GetSession().Query("UPDATE "+TagTable+" SET "+
		" name=?"+
		" WHERE name=?;",
		m.Name, m.Name).Exec()
}

func GetTagById(id string) *Tag {
	var data map[string]interface{} = make(map[string]interface{})

	GetSession().Query("SELECT * FROM "+TagTable+" WHERE name = ?;", id).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		return nil
	}

	return &Tag{
		Name: data["name"].(string),
	}
}
