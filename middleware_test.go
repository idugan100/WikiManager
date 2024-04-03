package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

func TestRequireAdmin(t *testing.T) {
	//test when user is not an admin
	req := httptest.NewRequest(http.MethodGet, "/create", nil)
	w := httptest.NewRecorder()

	handler := requireAdmin(createWikiPage)
	handler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("unexpected status code %d", w.Code)
	}

	//test when user is an admin
	req = httptest.NewRequest(http.MethodGet, "/create", nil)
	w = httptest.NewRecorder()

	// Set up a test session
	store := sessions.NewCookieStore([]byte("secret"))
	session, _ := store.Get(req, "admin")
	session.Values["authenticated"] = true
	session.Save(req, w)

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
	}
}
