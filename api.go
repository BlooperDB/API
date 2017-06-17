package blooper

import (
	"net/http"

	"fmt"

	"log"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/blueprint"
	"github.com/BlooperDB/API/db"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func Initialize() {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(api.LoggerHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 page not found")
	}))

	blueprint.RegisterBlueprintRoutes(api.RouteHandler(router, "/v1"))

	cluster := gocql.NewCluster("scylladb")
	session, _ := cluster.CreateSession()
	_, err := session.KeyspaceMetadata("blooper")
	if err != nil {
		session.Query("CREATE KEYSPACE blooper WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};").Exec()
	}
	session.Close()

	cluster = gocql.NewCluster("scylladb")
	cluster.Keyspace = "blooper"
	cluster.Consistency = gocql.Quorum
	session, _ = cluster.CreateSession()
	defer session.Close()

	db.Initialize(session)

	log.Fatal(http.ListenAndServe(":8080", router))
}
