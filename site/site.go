package site

import (
	"net/http"
)

func Start() {
		// Add all web page routes to the router.
	router := addRoutes()

	// Start the HTTP server
	if err := http.ListenAndServe(":4000", router); err != nil {
		panic(err)
	}
}