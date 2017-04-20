package site

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "html/template"
	// "log"
	// "net"
	// "net/http"
	// "path"
	// "strconv"

	"github.com/gorilla/mux"
)

// adds all the routes to the router.
func addRoutes() *mux.Router {
	router := mux.NewRouter()

	// handleLogin(router)
	// handleDashboardPage(router)
	// handleNetworksPage(router)
	// handleAccountsPage(router)

	// serveStaticFolder("/css/", router)
	// serveStaticFolder("/js/", router)
	// serveStaticFolder("/fonts/", router)

	// websocket.Start(router)

	// router.NotFoundHandler = http.HandlerFunc(generalNotFound)
	// http.Handle("/", router)

	return router
}
