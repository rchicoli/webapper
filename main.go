package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func hostname(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hostname: %s!", hostname)
	log.Printf("/: %s", hostname)

}

func cow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Println("/health: OK")
	fmt.Fprintf(w, "%v", `
 _________________________________________
/ You are only young once, but you can    \
\ stay immature indefinitely.             /
 -----------------------------------------
   \
    \
        .--.
       |o_o |
       |:_/ |
      //   \ \
     (|     | )
    /\_   _/\
    \___)=(___/
	`)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Println("/health: OK")
}

func main() {
	http.HandleFunc("/", hostname)

	http.HandleFunc("/cow", cow)

	http.HandleFunc("/health", healthcheck)

	http.ListenAndServe(":8080", nil)
}
