package blooper

import (
	"log"
	"net/http"

	"github.com/FactorioDB/API/api"
	"github.com/FactorioDB/API/blueprint"
	"github.com/gorilla/mux"
)

func Initialize() {
	router := mux.NewRouter()

	blueprint.RegisterBlueprintRoutes(api.RouteHandler(router, "/v1"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
