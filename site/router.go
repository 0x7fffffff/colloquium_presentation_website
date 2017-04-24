package site

import (
	// "encoding/json"
	"errors"
	"fmt"
	"html/template"
	// "log"
	// "net"
	"net/http"
	// "path"
	"strconv"
	// "bytes"

	"github.com/0x7fffffff/colloquium_presentation_website/websocket"
	"github.com/gorilla/mux"
)

var currentQuestion = 0
var totalQuestions = len(getQuiz())

func templateOnBase(path string) *template.Template {
	funcMap := template.FuncMap{
		"percentage": func(x, y int) int {
			return x / y * 100.0
		},
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
		"SocketURL":   			"ws://" + request.Host + "/socket",
		"CurrentQuestionIndex": currentQuestion,
		"TotalQuestions":		totalQuestions,
	}

	for k, v := range base {
		new[k] = v
	}

	return new
}

func handleControlPage(router *mux.Router) {
	router.HandleFunc("/control", func(writer http.ResponseWriter, request *http.Request) {
		controlTemplate := templateOnBase("templates/_control.html")
		data := map[string]interface{}{}
		if err := controlTemplate.Execute(writer, templateParamsOnBase(data, request)); err != nil {
			serverError(writer, err)
		}
	}).Methods("GET")

	router.HandleFunc("/control/next", func(writer http.ResponseWriter, request *http.Request) {
		currentQuestion++

		websocket.SocketMessage{
			Payload: map[string]interface{}{
				"next": map[string]interface{}{
					"question_number": currentQuestion,
				},
			},
		}.Send()

		writer.WriteHeader(http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/control/show", func(writer http.ResponseWriter, request *http.Request) {
		websocket.SocketMessage{
			Payload: map[string]interface{}{
				"show": map[string]interface{}{},
			},
		}.Send()

		writer.WriteHeader(http.StatusOK)
	}).Methods("POST")
}

func handleQuizPage(router *mux.Router) {
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := fmt.Sprintf("/question/%v", currentQuestion)
		http.Redirect(writer, request, path, http.StatusSeeOther)
	}).Methods("GET")

	router.HandleFunc("/question/{question_id:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
		quizTemplate := templateOnBase("templates/_quiz.html")
	
		id := identifierFromRequest("question_id", request)
		if id == nil {
			clientError(writer, errors.New("Missing question identifier"))
			return
		}

		quiz := getQuiz()

		if *id > len(quiz) - 1 {
			http.Redirect(writer, request, "/score", http.StatusSeeOther)
			return
		}

		if *id < 0 {
			clientError(writer, errors.New("Invalid question id"))
			return
		}

		if *id != currentQuestion {
			path := fmt.Sprintf("/question/%v", currentQuestion)
			http.Redirect(writer, request, path, http.StatusSeeOther)
			return			
		}

		data := map[string]interface{}{
			"Question": quiz[*id],
		}

		session, err := store.Get(request, "session")
		if err != nil {
			clientError(writer, errors.New("couldn't get session"))
			return
		}

		if session.IsNew {
			if err = session.Save(request, writer); err != nil {
				clientError(writer, errors.New("couldn't update session"))
				return
			}		 	
		} 

		if err := quizTemplate.Execute(writer, templateParamsOnBase(data, request)); err != nil {
			serverError(writer, err)
		}
	}).Methods("GET")

	router.HandleFunc("/score", func(writer http.ResponseWriter, request *http.Request) {
		scoreTemplate := templateOnBase("templates/_score.html")

		data := map[string]interface{}{}
		if err := scoreTemplate.Execute(writer, templateParamsOnBase(data, request)); err != nil {
			serverError(writer, err)
		}
	}).Methods("GET")

	router.HandleFunc("/question/{question_id:[0-9]+}/answer", func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			clientError(writer, err)
			return
		}

		id := identifierFromRequest("question_id", request)
		if id == nil {
			clientError(writer, errors.New("Missing question identifier"))
			return
		}

		_, err := store.Get(request, "session")
		if err != nil {
			clientError(writer, errors.New("couldn't get session"))
			return
		}

	}).Methods("POST")
}
// adds all the routes to the router.
func addRoutes() *mux.Router {
	router := mux.NewRouter()

	handleControlPage(router)
	handleQuizPage(router)

	serveStaticFolder("/css/", router)
	serveStaticFolder("/js/", router)
	serveStaticFolder("/fonts/", router)

	websocket.Start(router)

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

// parses the given identifier out of the request path.
func identifierFromRequest(identifier string, request *http.Request) *int {
	vars := mux.Vars(request)
	idString := vars[identifier]

	if idString == "" {
		return nil
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil
	}

	return &id
}
