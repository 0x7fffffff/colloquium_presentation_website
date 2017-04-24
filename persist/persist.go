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
	Body 		string
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

func GetQuestion(number int) (*Question, error) {
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

func GetAllQuestions() ([]Question, error) {
	query := `
		SELECT id, body, number, info, correct_index
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

func GetAnswersForQuestion(question Question) ([]Answer, error) {
	query := `
		SELECT id, question_id, answer_index, body
		FROM answer
		WHERE question_id = ?
		ORDER BY answer_index ASC;
	`

	rows, err := DB.Query(query, question.Id)
	if err != nil {
		return nil, err
	}

	var answers = make([]Answer, 0)

	for rows.Next() {
		var answer Answer

		if err = rows.Scan(&answer.Id, &answer.QuestionId, &answer.AnswerIndex, &answer.Body); err != nil {
			log.Fatal(err)
			continue
		}

		answers = append(answers, answer)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return answers, nil
}

func AnswerQuestion(questionId int, answerIndex int, sessionId string) error {
	query := `
		SELECT id
		FROM question_answer
		WHERE question_id = ?
			AND session_id = ?;		
	`

	row := DB.QueryRow(query, questionId, sessionId)
	var id int
	err := row.Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		insert := `
			INSERT INTO question_answer (
				question_id, answer_index, session_id
			) VALUES (
				?, ?, ?
			);
		`

		_, err := DB.Exec(insert, questionId, answerIndex, sessionId)
		return err
	case err != nil:
		return err
	default:
		update := `
			UPDATE question_answer
				SET
					answer_index = ?
				WHERE
					id = ?;
		`

		_, err := DB.Exec(update, answerIndex, id)
		return err
	}
}

func CorrectCountForSessionId(sessionId string) (int, error) {
	query := `
		SELECT count(*) 
		FROM question_answer
		INNER JOIN question
		ON question_answer.question_id = question.id
		WHERE question.correct_index = question_answer.answer_index
		  AND question_answer.session_id = ?;
	`

	row := DB.QueryRow(query, sessionId)
	if row == nil {
		return 0, errors.New("record not found")
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, errors.New("Failed to find specified Question")
	}

	return count, nil
}

func FindWinners() []string {
	return nil
}
