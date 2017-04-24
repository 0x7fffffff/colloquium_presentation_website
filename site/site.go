package site

import (
	"net/http"

	"github.com/0x7fffffff/colloquium_presentation_website/persist"
	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
)

var store *sqlitestore.SqliteStore

const (
	cookieSeed = "4'852b9FtL(_61R!q]La1d_BtEi8(*"
)

func init() {
	var err error
	store, err = sqlitestore.NewSqliteStoreFromConnection(
		persist.DB,
		"session",
		"/",
		86400 * 7,
		[]byte(cookieSeed))

	if err != nil {
		panic(err)
	}
}

func Start() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: false,
	}

	// Add all web page routes to the router.
	router := addRoutes()

	// Start the HTTP server
	if err := http.ListenAndServe(":4000", router); err != nil {
		panic(err)
	}
}