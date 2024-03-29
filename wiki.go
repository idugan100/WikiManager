package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("./tmpl/create.html", "./tmpl/edit.html", "./tmpl/view.html", "./tmpl/all.html"))
var fileValidator = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

func main() {
	fmt.Println("starting server")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /view/{path}", validatePath(viewWikiPage))
	mux.HandleFunc("GET /edit/{path}", validatePath(editWikiPage))
	mux.HandleFunc("POST /save/{path}", validatePath(saveWikiPage))
	mux.HandleFunc("GET /delete/{path}", validatePath(deleteWikiPage))
	mux.HandleFunc("GET /create", createWikiPage)
	mux.HandleFunc("POST /store", storeWikiPage)
	mux.HandleFunc("GET /", allWikiPages)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Errorf("Error with server: %w", err)
	}
}
