package db

import (
	"fmt"

	"os"

	_ "github.com/gemnasium/migrate/driver/cassandra"
	"github.com/gemnasium/migrate/migrate"
	"github.com/gocql/gocql"
)

var session *gocql.Session

func Initialize(s *gocql.Session) {
	session = s

	migrations()
}

func GetSession() *gocql.Session {
	return session
}

func migrations() {
	allErrors, ok := migrate.UpSync(
		"cassandra://scylladb:9042/blooper?protocol=4&consistency=all&disable_init_host_lookup",
		"./src/github.com/BlooperDB/API/migrations",
	)

	if !ok || len(allErrors) > 0 {
		fmt.Println(allErrors)
		os.Exit(1)
	}
}
