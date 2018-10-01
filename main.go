package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/gorilla/mux"
)

const (
	usage = `
Usage: curl http://127.0.0.1:8080/[endpoint]

    Perform HTTP request against the API for a specific action.

Endpoints:

    /           show this message
    /echo       return the payload sent
    /headers    display all headers
    /hostname   display hostname
    /log        log the message
    /trace      dump the request
    /health     check health status
`
)

type Payload struct {
	Message string `json:"message"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", usage)
	log.Printf("%s", r.URL.Path)
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func hostnameHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hostname: %s\n", hostname)
	log.Printf("%s", r.URL.Path)
}

func logHandler(w http.ResponseWriter, r *http.Request) {

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

func traceHandler(w http.ResponseWriter, req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprint(w, string(requestDump))
	}
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", endpointHandler).Methods("GET")
	router.HandleFunc("/echo", echoHandler).Methods("GET", "POST")
	router.HandleFunc("/headers", headersHandler).Methods("GET")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/hostname", hostnameHandler).Methods("GET")
	router.HandleFunc("/log", logHandler).Methods("POST")
	router.HandleFunc("/trace", traceHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}
