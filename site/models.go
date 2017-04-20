package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type quiz struct {
	Data []Question
}

type Question struct {
	Question		string
	Answers			[]string
	CorrectIndex	int
}

func getQuiz() []Question {
	file, err := ioutil.ReadFile("./quiz.json")
    if err != nil {
        fmt.Printf("JSON loading error: %v\n", err)
        return nil
    }

    var quizJSON quiz
    json.Unmarshal(file, &quizJSON)

    return quizJSON.Data
}