package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	/jsonp		beautify json object
	/raw		post raw request
	/trace		trace the request
    /echo       return the payload sent
    /headers    display all headers
    /health     check health status
    /hostname   display hostname
    /log        log the message
    /trace      dump the request
`
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", usage)
	log.Printf("%s", r.URL.Path)
}
func rawHandler(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s", b)
	log.Printf("%s", b)
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

func jsonPrettyPrintHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	var jsonObject interface{}
	json.Unmarshal(b, &jsonObject)
	jsonPretty, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s", jsonPretty)
	log.Printf("%s", jsonPretty)
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
	router.HandleFunc("/jsonp", jsonPrettyPrintHandler).Methods("POST")
	router.HandleFunc("/log", rawHandler).Methods("POST")
	router.HandleFunc("/raw", rawHandler).Methods("POST")
	router.HandleFunc("/trace", traceHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}
