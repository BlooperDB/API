// @APIVersion 1.0.0
// @APITitle Blooper
// @APIDescription Factorio blueprint database
// @License MIT
// @LicenseUrl https://opensource.org/licenses/MIT

package blooper

import (
	"log"
	"net/http"

	"github.com/FactorioDB/API/api"
	"github.com/FactorioDB/API/blueprint"
	"github.com/julienschmidt/httprouter"
)

func Initialize() {
	router := httprouter.New()

	blueprint.RegisterBlueprintRoutes(api.RouteHandler(router))

	log.Fatal(http.ListenAndServe(":8080", router))
}
