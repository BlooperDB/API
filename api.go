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
	"github.com/BlooperDB/API/storage"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/minio/minio-go"
	"github.com/wuman/firebase-server-sdk-go"
)

func Initialize() {
	var listenPort int
	var postgresHost string
	var minioHost string

	flag.IntVar(&listenPort, "postgres-port", 8080, "sets the port to run on")
	flag.StringVar(&postgresHost, "postgres-host", "postgres", "sets the postgres host to connect to")
	flag.StringVar(&minioHost, "minio-host", "minio", "sets the minio host to connect to")
	flag.Parse()

	firebase.InitializeApp(&firebase.Options{
		ServiceAccountPath: "src/github.com/BlooperDB/API/blooper-firebase-adminsdk.json",
	})

	var (
		db_user = os.Getenv("POSTGRES_USER")
		db_name = os.Getenv("POSTGRES_DB")
		db_pass = os.Getenv("POSTGRES_PASSWORD")
	)

	orm_cmd := "host=" + postgresHost + " user=" + db_user + " dbname=" + db_name + " sslmode=disable password=" + db_pass + ""
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

	var (
		minio_access_key = os.Getenv("MINIO_ACCESS_KEY")
		minio_secret_key = os.Getenv("MINIO_SECRET_KEY")
	)

	minioClient, err := minio.New(minioHost+":9000", minio_access_key, minio_secret_key, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	storage.Initialize(minioClient, os.Getenv("STORAGE_BASE_URL"))

	router := mux.NewRouter()

	router.NotFoundHandler = api.LoggerHandler(api.NotFoundHandler())

	v1 := api.RouteHandler(router, "/v1")
	nodes.RegisterUserRoutes(v1)
	nodes.RegisterBlueprintRoutes(v1)
	nodes.RegisterCommentRoutes(v1)
	nodes.RegisterRevisionRoutes(v1)
	nodes.RegisterTagRoutes(v1)

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

	fmt.Printf("Listening on port %d\n", listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), finalRouter))
}
