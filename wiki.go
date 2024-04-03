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

func setupServer() *http.ServeMux {
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

	return mux
}

func main() {
	fmt.Println("starting server")

	server := http.Server{
		Addr:    ":8080",
		Handler: setupServer(),
	}

	panic(server.ListenAndServe())
}
