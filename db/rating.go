package db

import "github.com/gocql/gocql"

var RatingTable = "rating"

type Rating struct {
	Id        string
	UserId    string
	VersionId string
	ThumbsUp  bool
}

func (m Rating) Save() {
	GetSession().Query("UPDATE "+RatingTable+" SET "+
		" version_id=?,"+
		" user_id=?,"+
		" thumbs_up=?"+
		" WHERE id=?;",
		m.VersionId, m.UserId, m.ThumbsUp, m.Id).Exec()
}

func FindRatingsByVersion(m Version) []*Rating {
	r := GetSession().Query("SELECT * FROM "+RatingTable+" WHERE version_id = ?;", m.Id).Consistency(gocql.All).Iter()

	result := make([]*Rating, r.NumRows())

	for i := 0; i < r.NumRows(); i++ {
		data := make(map[string]interface{})
		r.MapScan(data)
		result[i] = &Rating{
			Id:        data["id"].(string),
			UserId:    data["user_id"].(string),
			VersionId: data["version_id"].(string),
			ThumbsUp:  data["thumbs_up"].(bool),
		}
	}

	return result
}
