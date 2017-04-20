package site

import (
	// "encoding/json"
	// "errors"
	"fmt"
	"html/template"
	// "log"
	// "net"
	"net/http"
	// "path"
	// "strconv"
	// "bytes"

	"github.com/0x7fffffff/colloquium_presentation_website/websocket"
	"github.com/gorilla/mux"
)

func templateOnBase(path string) *template.Template {
	funcMap := template.FuncMap{
		// "inc": func(i int) int {
		// 	return i + 1
		// },
	}

	return template.Must(template.New("_base.html").Funcs(funcMap).ParseFiles(
		"templates/_base.html",
		path,
	))
}

// creates the base params that will be passed to all templates when
// they are rendered.
func templateParamsOnBase(new map[string]interface{}, request *http.Request) map[string]interface{} {
	base := map[string]interface{}{
		"SocketURL":   "ws://" + request.Host + "/socket",
	}

	for k, v := range base {
		new[k] = v
	}

	return new
}

func handleQuizPage(router *mux.Router) {
	quizTemplate := templateOnBase("templates/_quiz.html")
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := map[string]interface{}{}

	    fmt.Printf("test: %v\n", getQuiz())

		if err := quizTemplate.Execute(writer, templateParamsOnBase(data, request)); err != nil {
			serverError(writer, err)
		}
	}).Methods("GET")
}
// adds all the routes to the router.
func addRoutes() *mux.Router {
	router := mux.NewRouter()

	handleQuizPage(router)

	serveStaticFolder("/css/", router)
	serveStaticFolder("/js/", router)
	serveStaticFolder("/fonts/", router)

	websocket.Start(router)

	// router.NotFoundHandler = http.HandlerFunc(generalNotFound)
	// http.Handle("/", router)

	return router
}

// used to server static files, like CSS/JavaScript/fonts/etc.
func serveStaticFolder(folder string, router *mux.Router) {
	static := "static" + folder
	fileServer := http.FileServer(http.Dir(static))
	router.PathPrefix(folder).Handler(http.StripPrefix(folder, fileServer))
}

func clientError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusBadRequest)
}

func serverError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusInternalServerError)
}
