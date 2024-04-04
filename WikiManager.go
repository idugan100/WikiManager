package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/gorilla/sessions"
)

var templates = getTemplates("./tmpl/*")
var fileValidator = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
var secret []byte
var password string
var session_store = sessions.NewCookieStore(secret)

func getTemplates(pathToTemplates string) *template.Template {
	return template.Must(template.ParseGlob(pathToTemplates))
}

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
	p := flag.String("port", "8080", "port to run server on")
	secret = []byte(*flag.String("secret", "secret", "session secret"))
	password = *flag.String("password", "password", "admin password")

	flag.Parse()

	fmt.Printf("starting server on port %s", *p)

	server := http.Server{
		Addr:    ":" + *p,
		Handler: setupServer(),
	}

	panic(server.ListenAndServe())
}
