package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rchicoli/webapper/log"

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
    /health     check health status
    /hostname   display hostname
    /jsonp      beautify json object
    /log        log the message
    /metrics    prometheus metrics
    /raw        post raw request
    /trace      dump the request
`
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", usage)
	stdlog.Printf(r.URL.Path)
}

func rawHandler(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s", b)
	stdlog.Printf("%s", b)
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
	ctx := r.Context()
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hostname: %s\n", hostname)
	log.Printf(ctx, fmt.Sprintf("[context] %p: %s", ctx, r.URL.Path))
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
	stdlog.Printf("%s", jsonPretty)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
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

	router.HandleFunc("/", log.Decorate(endpointHandler)).Methods("GET")
	router.HandleFunc("/echo", log.Decorate(echoHandler)).Methods("GET", "POST")
	router.HandleFunc("/headers", log.Decorate(headersHandler)).Methods("GET")
	router.HandleFunc("/health", log.Decorate(healthCheckHandler)).Methods("GET")
	router.HandleFunc("/hostname", log.Decorate(hostnameHandler)).Methods("GET")
	router.HandleFunc("/jsonp", log.Decorate(jsonPrettyPrintHandler)).Methods("POST")
	router.HandleFunc("/log", log.Decorate(rawHandler)).Methods("POST")
	router.HandleFunc("/metrics", log.Decorate(metricsHandler)).Methods("GET")
	router.HandleFunc("/raw", log.Decorate(rawHandler)).Methods("POST")
	router.HandleFunc("/trace", log.Decorate(traceHandler)).Methods("GET", "POST")

	srv := http.Server{Addr: ":8080", Handler: router}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			stdlog.Printf("HTTP server Shutdown: %v", err)
		}

		stdlog.Println("HTTP server shutdown gracefully")
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		stdlog.Fatalf("HTTP server ListenAndServe: %v\n", err)
	}

	<-idleConnsClosed

}
