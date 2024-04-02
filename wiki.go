package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/sessions"
)

var templates = template.Must(template.ParseFiles("./tmpl/login.html", "./tmpl/create.html", "./tmpl/edit.html", "./tmpl/view.html", "./tmpl/all.html"))
var fileValidator = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
var secret = []byte(os.Getenv("GOWIKISECRET"))
var session_store = sessions.NewCookieStore(secret)

func main() {

	fmt.Println("starting server")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /view/{path}", isAdminAndValidatePath(viewWikiPage))
	mux.HandleFunc("GET /edit/{path}", requireAdmin(validatePath(editWikiPage)))
	mux.HandleFunc("POST /save/{path}", requireAdmin(validatePath(saveWikiPage)))
	mux.HandleFunc("GET /delete/{path}", requireAdmin(validatePath(deleteWikiPage)))
	mux.HandleFunc("GET /create", requireAdmin(createWikiPage))
	mux.HandleFunc("POST /store", requireAdmin(storeWikiPage))
	mux.HandleFunc("GET /", isAdmin(allWikiPages))
	mux.HandleFunc("GET /loginpage", loginScreen)
	mux.HandleFunc("GET /login", login)
	mux.HandleFunc("GET /logout", requireAdmin(logout))

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		_ = fmt.Errorf("error with server: %w", err)
	}
}
