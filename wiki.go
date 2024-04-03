package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/sessions"
)

var templates = getTemplates("./tmpl/*")
var fileValidator = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
var session_store = sessions.NewCookieStore(secret)
var secret, password = parseEnv()

func getTemplates(pathToTemplates string) *template.Template {
	return template.Must(template.ParseGlob(pathToTemplates))
}

func parseEnv() ([]byte, string) {
	return []byte(os.Getenv("GOWIKISECRET")), os.Getenv("GOWIKIPASSWORD")
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
	fmt.Println("starting server")

	server := http.Server{
		Addr:    ":8080",
		Handler: setupServer(),
	}

	panic(server.ListenAndServe())
}
