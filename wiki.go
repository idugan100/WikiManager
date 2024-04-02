package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/sessions"
)

var templates = template.Must(template.ParseGlob("./tmpl/*"))
var fileValidator = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
var secret = []byte(os.Getenv("GOWIKISECRET"))
var session_store = sessions.NewCookieStore(secret)
var password = os.Getenv("GOWIKIPASSWORD")

func main() {
	fmt.Println("starting server")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /view/{path}", isAdminAndValidatePath(viewWikiPage))
	mux.HandleFunc("GET /edit/{path}", requireAdmin(validatePath(editWikiPage)))
	mux.HandleFunc("POST /save/{path}", requireAdmin(validatePath(saveWikiPage)))
	mux.HandleFunc("GET /delete/{path}", requireAdmin(validatePath(deleteWikiPage)))
	mux.HandleFunc("GET /create", requireAdmin(createWikiPage))
	mux.HandleFunc("POST /store", requireAdmin(storeWikiPage))
	mux.HandleFunc("GET /logout", requireAdmin(logout))
	mux.HandleFunc("GET /", isAdmin(allWikiPages))
	mux.HandleFunc("GET /loginpage", loginScreen)
	mux.HandleFunc("GET /login", login)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	panic(server.ListenAndServe())
}
