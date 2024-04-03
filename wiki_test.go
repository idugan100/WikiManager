package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupServer(t *testing.T) {
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()
	defer deletePage("testWiki")

	requestList := []struct {
		Method     string
		Path       string
		ResultCode int
	}{
		{http.MethodGet, "/", http.StatusOK},
		{http.MethodGet, "/view/testWiki", http.StatusOK},
		{http.MethodGet, "/loginpage", http.StatusOK},
		{http.MethodGet, "/login", http.StatusUnauthorized},
	}

	for _, request := range requestList {
		req := httptest.NewRequest(request.Method, request.Path, nil)
		w := httptest.NewRecorder()
		setupServer().ServeHTTP(w, req)

		if w.Code != request.ResultCode {
			t.Errorf("expected code %d, got %d", request.ResultCode, w.Code)
		}
	}
}
