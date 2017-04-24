package site

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/0x7fffffff/colloquium_presentation_website/persist"
	"github.com/0x7fffffff/colloquium_presentation_website/websocket"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

var currentQuestion = 0
var totalQuestions int
var currentSessionHeaderName string
var over = false

func init() {
	questions, err := persist.GetAllQuestions()
	if err != nil {
		totalQuestions = 0
	} else {
		totalQuestions = len(questions)
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		currentSessionHeaderName = "session"
	} else {
		currentSessionHeaderName = "session-" + uuid.String()
	}
}

func templateOnBase(path string) *template.Template {
	funcMap := template.FuncMap{
		"percentage": func(x, y int) float64 {
			return float64(x) / float64(y) * 100.0
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

func win() {
	fmt.Println("over. about to sleep")
	time.Sleep(5 * time.Second)
	fmt.Println("done sleeping")

	sessions, err := persist.FindWinners(1)
	if err != nil {
		return
	}
	fmt.Println(sessions)
	websocket.SocketMessage{
		Payload: map[string]interface{}{
			"winners": map[string]interface{}{
				"sessions": sessions,
			},
		},
	}.Send()	
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
		if over {
			writer.WriteHeader(http.StatusOK)
			return
		}

		currentQuestion++

		websocket.SocketMessage{
			Payload: map[string]interface{}{
				"next": map[string]interface{}{
					"question_number": currentQuestion,
				},
			},
		}.Send()

		if currentQuestion >= totalQuestions {
			over = true
			go win()
		}

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
		session, err := store.Get(request, currentSessionHeaderName)
		if err != nil {
			clientError(writer, errors.New("couldn't get session"))
			return
		}

		if session.IsNew {
			uuid, err := uuid.NewV4()
			if err != nil {
				serverError(writer, errors.New("Couldn't generate user key"))
				return
			}

			session.Values["token"] = uuid.String()

			if err = session.Save(request, writer); err != nil {
				clientError(writer, errors.New("couldn't update session"))
				return
			}		 	
		} 

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

		if *id > totalQuestions - 1 {
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

		question, err := persist.GetQuestion(*id)
		if err != nil {
			clientError(writer, errors.New("Invalid question id"))
			return			
		}

		answers, err := persist.GetAnswersForQuestion(*question)
		if err != nil {
			clientError(writer, errors.New("Invalid question id"))
			return						
		}

		data := map[string]interface{}{
			"Question": *question,
			"Answers": answers,
		}

		if err := quizTemplate.Execute(writer, templateParamsOnBase(data, request)); err != nil {
			serverError(writer, err)
		}
	}).Methods("GET")

	router.HandleFunc("/score", func(writer http.ResponseWriter, request *http.Request) {
		scoreTemplate := templateOnBase("templates/_score.html")

		session, err := store.Get(request, currentSessionHeaderName)
		if err != nil {
			clientError(writer, errors.New("couldn't get session"))
			return
		}

		count, err := persist.CorrectCountForSessionId(session.Values["token"].(string))
		if err != nil {
			serverError(writer, errors.New("I don't even know"))
			return
		}

		data := map[string]interface{}{
			"Percentage": float64(count) / float64(totalQuestions) * 100.0,
			"Session": session.Values["token"].(string),
		}

		if err := scoreTemplate.Execute(writer, templateParamsOnBase(data, request)); err != nil {
			serverError(writer, err)
		}
	}).Methods("GET")

	router.HandleFunc("/question/{question_id:[0-9]+}/answer/{answer_index:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
		questionId := identifierFromRequest("question_id", request)
		if questionId == nil {
			clientError(writer, errors.New("Missing question identifier"))
			return
		}

		answerIndex := identifierFromRequest("answer_index", request)
		if answerIndex == nil {
			clientError(writer, errors.New("Missing answer index"))
			return
		}

		session, err := store.Get(request, currentSessionHeaderName)
		if err != nil {
			clientError(writer, errors.New("couldn't get session"))
			return
		}

		fmt.Println(session.Values["token"])
		persist.AnswerQuestion(*questionId + 1, *answerIndex, session.Values["token"].(string))

		writer.WriteHeader(http.StatusOK)
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
