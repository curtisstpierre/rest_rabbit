package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

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

// fun main()
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/alerts/{id}", CreateAlert).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
