package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Main
func main() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/v1").Subrouter()

	//registering router
	subRouter.HandleFunc("/questions/{id:\\d+}", QuestionsHandler).Methods("GET")
	subRouter.HandleFunc("/questions/{id:\\d+}/answers", AnswersHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":123", router))
	log.Println("Started server on port 123")
}
