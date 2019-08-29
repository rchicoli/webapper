package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	stdlog "log"

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

var (
	healthy int32
)

func rawHandler(w http.ResponseWriter, r *http.Request) {}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", usage)
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func hostnameHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf(r.Context(), "error: retrieving hostname")
		io.WriteString(w, err.Error())
		return
	}
	fmt.Fprintf(w, "hostname: %s\n", hostname)
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		log.Printf(r.Context(), "error: retrieving ip address")
		io.WriteString(w, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "ip: %s\n", string(body))
}

func jsonPrettyPrintHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var input interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
		http.Error(w, "error: reading request body", http.StatusInternalServerError)
		return
	}

	jsonP, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		http.Error(w, "error: marshaling json", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s", jsonP)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func traceHandler(w http.ResponseWriter, req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(requestDump)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", log.Decorate(endpointHandler)).Methods("GET")
	router.HandleFunc("/echo", log.Decorate(echoHandler)).Methods("GET", "POST")
	router.HandleFunc("/headers", log.Decorate(headersHandler)).Methods("GET")
	router.HandleFunc("/health", log.Decorate(healthCheckHandler)).Methods("GET")

	router.HandleFunc("/livez", log.Decorate(healthCheckHandler)).Methods("GET")
	router.HandleFunc("/healthz", log.Decorate(healthCheckHandler)).Methods("GET")

	router.HandleFunc("/ip", log.Decorate(ipHandler)).Methods("GET")
	router.HandleFunc("/hostname", log.Decorate(hostnameHandler)).Methods("GET")
	router.HandleFunc("/jsonp", log.Decorate(jsonPrettyPrintHandler)).Methods("POST")
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
		atomic.StoreInt32(&healthy, 0)

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			stdlog.Printf("HTTP server Shutdown: %v", err)
		}

		stdlog.Println("HTTP server shutdown gracefully")
		close(idleConnsClosed)
	}()

	atomic.StoreInt32(&healthy, 1)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		stdlog.Fatalf("HTTP server ListenAndServe: %v\n", err)
	}

	<-idleConnsClosed

}
