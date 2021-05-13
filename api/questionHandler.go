package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// QuestionsHandler encodes questions into JSON format and writes to the response.
func QuestionsHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("Got new request for QuestionsHandler", request.URL.Path)
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		ReturnStatusBadRequest(response, "Invalid ID.")
		return
	}

	Questions, questionsExist := QuestionsDataSet[id]
	if !questionsExist {
		ReturnStatusNotFound(response)
		return
	}

	// cache it on the client side
	response.Header().Set("Cache-Control", "max-age=300") //5 minutes
	json.NewEncoder(response).Encode(Questions)

	log.Println("Responded with HTTP status 200 OK")
}

// AnswersHandler saves client's answers.
func AnswersHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("Got new request for AnswersHandler", request.URL.Path)

	// ID is not an integer
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		ReturnStatusBadRequest(response, "ID must be an integer")
		return
	}

	// invalid ID was supplied
	Questions, questionsExist := QuestionsDataSet[id]
	if !questionsExist {
		ReturnStatusNotFound(response)
		return
	}

	// parsing JSON
	var answers []struct {
		AnswerID int `json:"answerId"`
	}
	err = json.NewDecoder(request.Body).Decode(&answers)
	if err != nil {
		ReturnStatusBadRequest(response, "")
		return
	}

	// checking that all questions are answered and there are no "extra" answers
	totalQuestions := len(Questions)

	if totalQuestions != len(answers) {
		ReturnStatusBadRequest(response, "Invalid amount of answers.")
		return
	}

	// calculate total correct answers
	correctAnswers := 0
	for questionIx, answer := range answers {
		if answer.AnswerID == Questions[questionIx].IndexAnswer {
			correctAnswers++
		}
	}

	// respond with user's result
	userIsBetterThanOthers := getUserIsBetterThanOthersPerc(correctAnswers, totalQuestions)

	// updating results
	UsersAnsweredCorrectly[correctAnswers]++
	TotalAnsweredUsers++

	res := struct {
		CorrectAnswers    int `json:"correctAnswers"`
		AnsweredUsers     int `json:"answeredUsers"`
		ComparingToOthers int `json:"comparingToOthers"`
	}{
		CorrectAnswers:    correctAnswers,
		AnsweredUsers:     TotalAnsweredUsers,
		ComparingToOthers: userIsBetterThanOthers,
	}
	json.NewEncoder(response).Encode(res)

	log.Println("Responded with HTTP status 200 OK")
}

// Returns by how many percentage current user has answered better than the rest
func getUserIsBetterThanOthersPerc(correctAnswers int, totalQuestions int) int {
	if correctAnswers == 0 {
		return 0
	}
	if correctAnswers == totalQuestions || TotalAnsweredUsers == 0 {
		return 100
	}

	// calculating percentage
	usersAnsweredWorse := 0
	for i := 0; i <= correctAnswers; i++ {
		usersAnsweredWorse += UsersAnsweredCorrectly[i]
	}

	if usersAnsweredWorse == 0 {
		return 0
	}

	return int(math.Ceil(float64(usersAnsweredWorse) / float64(TotalAnsweredUsers) * 100.0))
}
