package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

const (
	usage = `
Usage: curl http://127.0.0.1:8080/[endpoint]

    Perform HTTP request against the API for a specific action.

Endpoints:

    /            show this message
    /cowsay      talking cow
    /health      check health status
    /hostname    display hostname
`

	cowsay = `
 _______________________________________
|  You are only young once, but you can |
|  stay immature indefinitely.          |
 ---------------------------------------
    \   ^__^
     \  (oo)\_______
        (__)\       )\/\
            ||----w |
            ||     ||
`
)

type RabbitMQ struct {
	ch *amqp.Channel
}

type Log struct {
	Payload Payload
	RabbitMQ
}

type Payload struct {
	Message string `json:"message"`
	Queue   string `json:"queue"`
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", usage)
	log.Printf("%s", r.URL.Path)
}

func hostnameHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hostname: %s\n", hostname)
	log.Printf("%s", r.URL.Path)
}

func cowsayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", cowsay)
	log.Printf("%s", r.URL.Path)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func (msg *Log) logMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// b, err := ioutil.ReadAll(r.Body)
	// err = json.Unmarshal(b, &msg.Payload)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg.Payload)
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	defer r.Body.Close()

	fmt.Fprintf(w, "%q", msg.Payload.Message)
	fmt.Printf("%q", msg.Payload.Message)

}

func (msg *Log) rabbitPublish(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// b, err := ioutil.ReadAll(r.Body)
	// err = json.Unmarshal(b, &msg.Payload)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg.Payload)
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	defer r.Body.Close()

	q, err := msg.ch.QueueDeclare(
		msg.Payload.Queue, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = msg.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(msg.Payload.Message),
		})
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, "error: publishing message to rabbitmq", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%q", msg.Payload.Message)
	fmt.Printf("%q", msg.Payload.Message)
}

func (msg *Log) rabbitConsume(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	msg.Payload.Queue = params["queue"]

	q, err := msg.ch.QueueDeclare(
		msg.Payload.Queue, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := msg.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	loop := true
	ticker := time.NewTicker(5 * time.Second)

	var bodies []byte
	bodies = make([]byte, 0)

	for loop {
		select {
		case d := <-msgs:
			bodies = append(bodies, d.Body...)
			fmt.Printf("%q", d.Body)
		case <-ticker.C:
			fmt.Printf("%q", "timeout")
			loop = false
		}
	}
	fmt.Fprintf(w, "Received a message: %s", string(bodies))

}

var re = regexp.MustCompile(`(\r\n\s+$)|(\r\n$)|(\n\s+$)|(\n$)`)

func main() {

	username := flag.String("username", "guest", "username to authenticate to rabbitmq")
	password := flag.String("password", "guest", "password to authenticate to rabbitmq")
	host := flag.String("host", "localhost", "rabbitmq host")
	port := flag.Int("port", 5672, "port which rabbitmq is listening to")

	flag.Parse()

	m := &Log{
		Payload: Payload{
			Queue: "test",
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/", endpointHandler).Methods("GET")
	router.HandleFunc("/cowsay", cowsayHandler).Methods("GET")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/hostname", hostnameHandler).Methods("GET")
	router.HandleFunc("/log", m.logMessage).Methods("POST")
	router.HandleFunc("/rabbit/publish", m.rabbitPublish).Methods("POST")
	router.HandleFunc("/rabbit/consume/{queue}", m.rabbitConsume).Methods("GET")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", *username, *password, *host, *port))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	m.ch, err = conn.Channel()
	defer m.ch.Close()
	failOnError(err, "Failed to open a channel")

	log.Fatal(http.ListenAndServe(":8080", router))

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
