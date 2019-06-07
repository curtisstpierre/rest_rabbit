package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zond/gotomic"
)

var tokens = gotomic.NewList()
var r *rand.Rand

// IssueToken issue a temporary token for the AMQP producer
func IssueToken(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["secret"] == "trytoguessme" {
		token := RandomString(64)
		tokens.Push(token)
		log.Printf(" [x] Authorized %s", token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

// fun main()
func main() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	router := mux.NewRouter()
	router.HandleFunc("/messages/{id}", CreateMessage).Methods("POST")
	router.HandleFunc("/auth/{secret}", IssueToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
