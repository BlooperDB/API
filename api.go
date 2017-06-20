package blooper

import (
	"net/http"

	"fmt"

	"log"

	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/nodes"
	"github.com/BlooperDB/API/utils"
	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wuman/firebase-server-sdk-go"
)

func Initialize() {
	firebase.InitializeApp(&firebase.Options{
		ServiceAccountPath: "src/github.com/BlooperDB/API/blooper-firebase-adminsdk.json",
	})

	router := mux.NewRouter()

	router.NotFoundHandler = api.LoggerHandler(api.NotFoundHandler())

	v1 := api.RouteHandler(router, "/v1")
	nodes.RegisterUserRoutes(v1)
	nodes.RegisterBlueprintRoutes(v1)

	cluster := gocql.NewCluster("cassandra")

	for i := 0; i < 60; i++ {
		stop := true

		utils.Block{
			Try: func() {
				session, _ := cluster.CreateSession()
				_, err := session.KeyspaceMetadata("blooper")
				if err != nil {
					session.Query("CREATE KEYSPACE blooper WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};").Exec()
				}
				session.Close()
			},
			Catch: func(e utils.Exception) {
				stop = false
				time.Sleep(1 * time.Second)
			},
			Finally: func() {
			},
		}.Do()

		if stop {
			break
		}
	}

	cluster = gocql.NewCluster("cassandra")
	cluster.Keyspace = "blooper"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	db.Initialize(session)

	CORSHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"content-type", "blooper-token"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}),
	)

	var finalRouter http.Handler = router
	finalRouter = CORSHandler(finalRouter)
	finalRouter = api.LoggerHandler(finalRouter)
	finalRouter = handlers.CompressHandler(finalRouter)
	finalRouter = handlers.ProxyHeaders(finalRouter)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", finalRouter))
}
