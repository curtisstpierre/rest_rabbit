package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/zond/gotomic"
)

var tokens = gotomic.NewList()
var r *rand.Rand

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Alert struct
type Alert struct {
	ID      string    `json:"id,omitempty"`
	Silence string    `json:"silence,omitempty"`
	Time    time.Time `json:"time,omitempty"`
}

// CreateAlert a new item
func CreateAlert(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auth := r.Header.Get("X-Vault-Auth")
	log.Printf(" [x] Tokens %s", tokens)
	if tokens.Search(String(auth)) == nil {
		log.Printf(" [x] Token not found")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if len(params["id"]) > 0 {
		var alert Alert
		_ = json.NewDecoder(r.Body).Decode(&alert)
		alert.ID = params["id"]
		amqpCaller(alert)
	}
}

func amqpCaller(alert Alert) {
	user := "guest"
	password := "guest"
	server := "localhost"
	port := 5672
	connectionStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, server, port)

	conn, err := amqp.Dial(connectionStr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		alert.ID, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := alert.Silence
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func RandomString(len int) string {
    bytes := make([]byte, len)
    for i := 0; i < len; i++ {
        bytes[i] = byte(65 + r.Intn(25))
    }
    return string(bytes)
}

func compStrings(i, j string) int {
        l := len(i)
        if len(j) < l {
                l = len(j)
        }
        for ind := 0; ind < l; ind++ {
                if i[ind] < j[ind] {
                        return -1
                } else if i[ind] > j[ind] {
                        return 1
                }
        }
        if len(i) < len(j) {
                return -1
        } else if len(i) > len(j) {
                return 1
        }
        if i == j {
                return 0
        }
        panic(fmt.Errorf("wtf, %v and %v are not the same!", i, j))
}


type String string

func (s String) Compare(t gotomic.Thing) int {
	return compStrings(string(s), t.(string))
}

// Issue a temporary token for the AMQP producer
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
	router.HandleFunc("/alerts/{id}", CreateAlert).Methods("POST")
	router.HandleFunc("/auth/{secret}", IssueToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
