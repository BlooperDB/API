package db

import "github.com/gocql/gocql"

var VersionTable = "version"

type Version struct {
	Id          string
	BlueprintId string
	Version     string
	Changes     string
	Date        int64
	Blueprint   string
}

func (m Version) Save() {
	GetSession().Query("UPDATE "+VersionTable+" SET "+
		" blueprint_id=?,"+
		" version=?,"+
		" changes=?,"+
		" date=?,"+
		" blueprint=?"+
		" WHERE id=?;",
		m.BlueprintId, m.Version, m.Changes, m.Date, m.Blueprint, m.Id).Exec()
}

func FindVersionsByBlueprint(b Blueprint) []*Version {
	r := GetSession().Query("SELECT * FROM "+VersionTable+" WHERE blueprint_id = ?;", b.Id).Consistency(gocql.All).Iter()

	result := make([]*Version, r.NumRows())

	for i := 0; i < r.NumRows(); i++ {
		data := make(map[string]interface{})
		r.MapScan(data)
		result[i] = &Version{
			Id:          data["id"].(string),
			BlueprintId: data["blueprint_id"].(string),
			Version:     data["version"].(string),
			Changes:     data["changes"].(string),
			Date:        data["date"].(int64),
			Blueprint:   data["blueprint"].(string),
		}
	}

	return result
}

func (m Version) GetComments() []*Comment {
	return FindCommentsByVersion(m)
}

func (m Version) GetRatings() []*Rating {
	return FindRatingsByVersion(m)
}
