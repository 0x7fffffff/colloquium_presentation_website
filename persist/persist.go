package persist

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	// Linter
	_ "github.com/mattn/go-sqlite3"
)

type Question struct {
	Id 				int
	Body 			string
	Number			int
	Info			string
	CorrectIndex	int
}

type Answer struct {
	Id 			int
	QuestionId 	int
	AnswerIndex int
}

// DB lint
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	_, err = configureDatabase(DB)
	if err != nil {
		panic(err)
	}
}

func configureDatabase(database *sql.DB) (sql.Result, error) {
	bytes, err := ioutil.ReadFile("./persist/ddl.sql")
	if err != nil {
		return nil, err
	}

	return database.Exec(string(bytes))
}

		// create table if not exists admin (
		// 	id integer primary key,
		// 	handle text unique not null,
		// 	hashed_password text not null
		// );

		// create table if not exists network (
		// 	id integer primary key,
		// 	listening_port integer unique not null,
		// 	name text not null,
		// 	timeout integer not null
		// );

		// create table if not exists encoder (
		// 	id integer primary key,
		// 	ip_address text not null,
		// 	port integer not null default(23),
		// 	name text null default ('New Encoder'),
		// 	handle text not null,
		// 	password text not null,
		// 	network_id integer not null,
		// 	foreign key(network_id) references network(id)
		// );

		// create table if not exists backup (
		// 	id integer primary key,
		// 	payload text not null,
		// 	network_id integer not null,
		// 	foreign key(network_id) references network(id)
		// );


func getQuestion(number int) (*Question, error) {
	query := `
		SELECT id, body, number, info, correct_index
		FROM question
		WHERE number = ?
	`

	row := DB.QueryRow(query, number)
	if row == nil {
		return nil, errors.New("Question not found")
	}

	var question Question
	if err := row.Scan(&question.Id, &question.Body, &question.Number, &question.Info, &question.CorrectIndex); err != nil {
		return nil, errors.New("Failed to find specified Question")
	}

	return &question, nil
}

func getAllQuestions() ([]Question, error) {
	query := `
		SELECT id, body, number
		FROM question
		ORDER BY number ASC;
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	var questions = make([]Question, 0)

	for rows.Next() {
		var question Question

		if err = rows.Scan(&question.Id, &question.Body, &question.Number, &question.Info, &question.CorrectIndex); err != nil {
			log.Fatal(err)
			continue
		}

		questions = append(questions, question)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return questions, nil
}

func getAnswersForQuestionId(id int) ([]Answer, error) {
	// query := `
	// 	SELECT id, question_id, answer_index
	// 	FROM answer
	// 	ORDER BY answer_index ASC;
	// `

	return nil, nil
}



// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// )

// type quiz struct {
// 	Data []Question
// }

// type Question struct {
// 	Question		string
// 	Answers			[]string
// 	CorrectIndex	int
// }

// func getQuiz() []Question {
// 	file, err := ioutil.ReadFile("./quiz.json")
//     if err != nil {
//         fmt.Printf("JSON loading error: %v\n", err)
//         return nil
//     }

//     var quizJSON quiz
//     json.Unmarshal(file, &quizJSON)

//     return quizJSON.Data
// }