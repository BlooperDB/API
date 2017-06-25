package blooper

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/nodes"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/wuman/firebase-server-sdk-go"
)

func Initialize() {
	var port int
	var dbhost string
	flag.IntVar(&port, "port", 8080, "sets the port to run on")
	flag.StringVar(&dbhost, "dbhost", "postgres", "sets the db host to query")
	flag.Parse()

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
	nodes.RegisterTagRoutes(v1)

	var (
		db_user = os.Getenv("POSTGRES_USER")
		db_name = os.Getenv("POSTGRES_DB")
		db_pass = os.Getenv("POSTGRES_PASSWORD")
	)

	orm_cmd := "host=" + dbhost + " user=" + db_user + " dbname=" + db_name + " sslmode=disable password=" + db_pass + ""
	connection, err := gorm.Open("postgres", orm_cmd)

	if err != nil {
		time.Sleep(5 * time.Second)
		connection, err = gorm.Open("postgres", orm_cmd)
		if err != nil {
			panic("failed to connect database")
		}
	}

	defer connection.Close()

	utils.Initialize()
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

	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), finalRouter))
}
