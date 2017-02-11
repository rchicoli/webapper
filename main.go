package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	usage = `
Usage: curl http://127.0.0.1:8080/[endpoint]

    Perform HTTP request against the API for a specific action. 

Endpoints:

    /            show this message
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

func main() {

	http.HandleFunc("/", endpointHandler)
	http.HandleFunc("/cow", cowsayHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/hostname", hostnameHandler)

	http.ListenAndServe(":8080", nil)

}
