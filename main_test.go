package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestMainHealthCheckHandler(t *testing.T) {

	response := httptest.NewRecorder()
	atomic.StoreInt32(&healthy, 1)

	request, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(healthCheckHandler)
	handler.ServeHTTP(response, request)

	if result := response.Code; result != 204 {
		t.Errorf("healthcheck hander returned wrong response string: got %d expected OK", result)
	}

}
