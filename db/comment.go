package db

import (
	"github.com/gocql/gocql"
)

var CommentTable = "comment"

type Comment struct {
	Id        string
	VersionId string
	UserId    string
	Date      int64
	Message   string
	Updated   int64
}

func (m Comment) Save() {
	GetSession().Query("UPDATE "+CommentTable+" SET "+
		" version_id=?,"+
		" user_id=?,"+
		" date=?,"+
		" message=?,"+
		" updated=?"+
		" WHERE id=?;",
		m.VersionId, m.UserId, m.Date, m.Message, m.Updated, m.Id).Exec()
}

func FindCommentsByVersion(m Version) []*Comment {
	r := GetSession().Query("SELECT * FROM "+CommentTable+" WHERE version_id = ?;", m.Id).Consistency(gocql.All).Iter()

	result := make([]*Comment, r.NumRows())

	for i := 0; i < r.NumRows(); i++ {
		data := make(map[string]interface{})
		r.MapScan(data)
		result[i] = &Comment{
			Id:        data["id"].(string),
			VersionId: data["version_id"].(string),
			UserId:    data["user_id"].(string),
			Date:      data["date"].(int),
			Message:   data["message"].(string),
			Updated:   data["updated"].(int),
		}
	}

	return result
}
