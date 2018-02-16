package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

type Payload struct {
	Message string `json:"message"`
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
	log.Printf("%s", r.URL.Path)
}

func logMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var msg Payload
	err := decoder.Decode(&msg)
	if err != nil {
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	fmt.Fprint(w, msg.Message)
	fmt.Println(msg.Message)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", endpointHandler).Methods("GET")
	router.HandleFunc("/cowsay", cowsayHandler).Methods("GET")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/hostname", hostnameHandler).Methods("GET")
	router.HandleFunc("/log", logMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}
