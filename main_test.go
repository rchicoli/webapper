package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHealthCheckHandler(t *testing.T) {

	response := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(healthCheckHandler)
	handler.ServeHTTP(response, request)

	if result := response.Body.String(); result != "OK" {
		t.Errorf("healthcheck hander returned wrong response string: got %s expected OK", result)
	}

}

func TestMainCowsayHandler(t *testing.T) {

	response := httptest.NewRecorder()
	cowsay := `
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

	request, err := http.NewRequest("GET", "/cowsay", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(cowsayHandler)
	handler.ServeHTTP(response, request)

	if result := response.Body.String(); result != cowsay {
		t.Errorf("healthcheck hander returned wrong response string: got %s expected %s", result, cowsay)
	}

}
