package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
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

func TestViewWikiPage(t *testing.T) {
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()
	defer deletePage(p.Title)

	req := httptest.NewRequest(http.MethodGet, "/view/testWiki", nil)
	w := httptest.NewRecorder()
	viewWikiPage(w, req, p.Title, true)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
	}

	body, _ := io.ReadAll(w.Result().Body)
	if !strings.Contains(string(body), string(p.Body)) {
		t.Errorf("returned view html does not contain body")
	}
}

func TestEditWikiPage(t *testing.T) {
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()
	defer deletePage(p.Title)
	req := httptest.NewRequest(http.MethodGet, "/edit/testWiki", nil)
	w := httptest.NewRecorder()

	editWikiPage(w, req, "testWiki")

	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code on edit page %d", w.Code)
	}

}

func TestSaveWikiPage(t *testing.T) {
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()
	defer deletePage(p.Title)
	req := httptest.NewRequest(http.MethodPost, "/save/testWiki", nil)
	req.Form = url.Values{}
	req.Form.Add("body", "new test wiki body")
	w := httptest.NewRecorder()

	saveWikiPage(w, req, "testWiki")

	if w.Code != http.StatusFound {
		t.Errorf("unexpected status code on save page %d", w.Code)
	}
	newPage, err := loadPage("testWiki")

	if err != nil {
		t.Errorf("error loading saved page: %s", err.Error())
	}

	if string(newPage.Body) != "new test wiki body" {
		t.Errorf("new body does not match \"new test wiki body\": %s", newPage.Body)
	}
}

func TestCreateWikiPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/create", nil)
	w := httptest.NewRecorder()
	createWikiPage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code on create page %d", w.Code)
	}

	body, _ := io.ReadAll(w.Result().Body)

	if !strings.Contains(string(body), "Create new wiki") {
		t.Errorf("incorrect html")
	}
}

func TestStoreWikiPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/store", nil)
	req.Form = url.Values{}
	req.Form.Add("title", "testWiki")
	req.Form.Add("body", "this is the testwiki body")
	w := httptest.NewRecorder()

	storeWikiPage(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("unexpected status code %d", w.Code)
	}

	p, err := loadPage("testWiki")

	if err != nil {
		t.Errorf("error when loading created wiki: %s", err.Error())
	}
	if p.Title != "testWiki" || string(p.Body) != "this is the testwiki body" {
		t.Errorf("created wiki page does not match the submitted wikie")
	}

	deletePage("testWiki")
}

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
