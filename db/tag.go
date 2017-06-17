package db

import (
	"github.com/gocql/gocql"
)

var TagTable = [2]string{
	"tag",
	"CREATE TABLE IF NOT EXISTS tag (" +
		"name varchar PRIMARY KEY" +
		");",
}

type Tag struct {
	Name string
}

func (m Tag) Save() {
	GetSession().Query("UPDATE "+TagTable[0]+" SET "+
		" name=?"+
		" WHERE name=?;",
		m.Name, m.Name).Exec()
}

func GetTagById(id string) *Tag {
	var data map[string]interface{} = make(map[string]interface{})

	GetSession().Query("SELECT * FROM "+TagTable[0]+" WHERE name = ?;", id).Consistency(gocql.One).MapScan(data)

	if len(data) == 0 {
		return nil
	}

	return &Tag{
		Name: data["name"].(string),
	}
}
