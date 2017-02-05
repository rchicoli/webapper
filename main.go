package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	usage = `
    Endpoints:
        /            show this response message
        /health      check health status
        /hostname    display hostname`

	cowsay = `
      _______________________________________
    /  You are only young once, but you can   \
    \  stay immature indefinitely.            /
      ---------------------------------------
            \   ^__^
             \  (oo)\_______
                (__)\       )\/\
                     ||----w |
                     ||     ||`
)

func endpoints(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", usage)
	log.Printf("%s", r.URL.Path)
}

func hostname(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hostname: %s", hostname)
	log.Printf("%s", r.URL.Path)
}

func cow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", cowsay)
	log.Printf("%s", r.URL.Path)

}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Printf("%s", r.URL.Path)
}

func main() {

	http.HandleFunc("/", endpoints)
	http.HandleFunc("/cow", cow)
	http.HandleFunc("/health", health)
	http.HandleFunc("/hostname", hostname)

	http.ListenAndServe(":8080", nil)

}
