package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleRequestOK(t *testing.T) {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()

	handler := &LambdaHandler{
		Config: &Config{},
	}
	handler.Config = &Config{
		Target: ts.URL,
	}
	res, err := handler.HandleRequest(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res != "OK" {
		t.Errorf("Expected OK, got %s", res)
	}
}

func TestHandleRequestTimeout(t *testing.T) {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK?"))
	}))
	defer ts.Close()

	handler := &LambdaHandler{}
	handler.Config = &Config{
		Target:  ts.URL,
		Timeout: 500 * time.Microsecond,
	}
	_, err := handler.HandleRequest(ctx)
	if err == nil {
		t.Errorf("Expected error")
	}

	t.Logf("error: %s", err)
}

func TestHandleRequestError(t *testing.T) {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
	}))
	defer ts.Close()

	handler := &LambdaHandler{}
	handler.Config = &Config{
		Target: ts.URL,
	}
	_, err := handler.HandleRequest(ctx)
	if err == nil {
		t.Errorf("Expected error")
	}

	t.Logf("error: %s", err)
}
