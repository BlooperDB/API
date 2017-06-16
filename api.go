package blooper

import (
	"log"
	"net/http"

	"fmt"

	"github.com/BlooperDB/API/api"
	"github.com/BlooperDB/API/blueprint"
	"github.com/gorilla/mux"
)

func Initialize() {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(api.LoggerHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 page not found")
	}))

	blueprint.RegisterBlueprintRoutes(api.RouteHandler(router, "/v1"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
