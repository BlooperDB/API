package db

import (
	"github.com/gocql/gocql"
)

var session *gocql.Session

func Initialize(s *gocql.Session) {
	session = s

	createTables()
}

func GetSession() *gocql.Session {
	return session
}

func createTables() {
	session.Query(BlueprintTable[1]).Exec()
	session.Query(TagTable[1]).Exec()
	session.Query(VersionTable[1]).Exec()
	session.Query(RatingTable[1]).Exec()
	session.Query(CommentTable[1]).Exec()
	session.Query(BlueprintTagTable[1]).Exec()
	session.Query(UserTable[1]).Exec()
}
