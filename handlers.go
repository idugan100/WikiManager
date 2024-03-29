package main

import (
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func viewWikiPage(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "view.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func editWikiPage(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "edit.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveWikiPage(w http.ResponseWriter, r *http.Request, filename string) {
	p := &Page{Title: filename, Body: []byte(r.FormValue("body"))}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+filename, http.StatusFound)
}

func createWikiPage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func storeWikiPage(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("title")
	body := []byte(r.FormValue("body"))
	p := &Page{Title: filename, Body: body}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+filename, http.StatusFound)
}

func allWikiPages(w http.ResponseWriter, r *http.Request) {
	//load all pages from content folder
	path := "./content/"
	fileSystem := os.DirFS(path)
	var pageList []Page
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			p := &Page{}
			p, err := loadPage(strings.TrimSuffix(d.Name(), ".txt"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			pageList = append(pageList, *p)

		}
		return nil
	})

	err := templates.ExecuteTemplate(w, "all.html", pageList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func deleteWikiPage(w http.ResponseWriter, r *http.Request, filename string) {
	err := os.Remove("./content/" + filename + ".txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
