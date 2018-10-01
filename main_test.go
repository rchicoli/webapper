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
