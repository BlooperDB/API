package blooper

import (
	"net/http"

	"fmt"

	"log"

	"os"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/nodes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
	nodes.RegisterCommentRoutes(v1)
	nodes.RegisterRevisionRoutes(v1)

	var (
		db_user = os.Getenv("POSTGRES_USER")
		db_name = os.Getenv("POSTGRES_DB")
		db_pass = os.Getenv("POSTGRES_PASSWORD")
	)

	connection, err := gorm.Open("postgres", "host=postgres user="+db_user+" dbname="+db_name+" sslmode=disable password="+db_pass+"")

	if err != nil {
		panic("failed to connect database")
	}

	defer connection.Close()

	db.Initialize(connection)

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
