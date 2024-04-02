package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestPageSaveLoadDelete(t *testing.T) {
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()

	p1, err := loadPage("testWiki")

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if string(p1.Body) != string(p.Body) {
		t.Errorf("Body is incorrect saved: %s loaded: %s", string(p.Body), string(p1.Body))
	}

	if p1.Title != p.Title {
		t.Errorf("Title is incorrect saved: %s loaded: %s", string(p.Title), string(p1.Title))
	}

	err = deletePage("testWiki")

	if err != nil {
		t.Errorf("Error when deleting Wikipage")
	}
}

func TestAllWikiRoute(t *testing.T) {
	//setup and tear down
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()
	defer deletePage("testWiki")

	//make request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	allWikiPages(w, req, true)

	//check status code
	if w.Code != http.StatusOK {

		t.Errorf("error fetching / route status code of %d", w.Result().StatusCode)
	}

	//check rendered html
	body, _ := io.ReadAll(w.Result().Body)

	if !strings.Contains(string(body), "testWiki") {
		t.Errorf("all wiki page / does not contain all wikis")
	}

}

func TestLoginScreen(t *testing.T) {
	req := httptest.NewRequest("GET", "/loginscreen", nil)
	w := httptest.NewRecorder()
	loginScreen(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error fetching login screen, status code of %d", w.Result().StatusCode)
	}

	body, _ := io.ReadAll(w.Result().Body)

	if !strings.Contains(string(body), "log in") {
		t.Errorf("incorrect html served back")
	}
}

func TestLogin(t *testing.T) {
	password = os.Getenv("GOWIKIPASSWORD")

	req := httptest.NewRequest("POST", "/login", nil)
	req.Form = url.Values{}
	req.Form.Add("password", password)

	w := httptest.NewRecorder()

	login(w, req)
	if w.Code != http.StatusFound {
		t.Errorf("unexpected status code returned from login %d", w.Code)
	}
}
