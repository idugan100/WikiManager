package main

import (
	"net/http"
	"os"
)

func viewWikiPage(w http.ResponseWriter, r *http.Request, filename string, isAdmin bool) {
	p, err := loadPage(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Page    Page
		IsAdmin bool
	}{
		Page:    *p,
		IsAdmin: isAdmin,
	}
	err = templates.ExecuteTemplate(w, "view.html", data)
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

func allWikiPages(w http.ResponseWriter, r *http.Request, isAdmin bool) {

	search := r.URL.Query().Get("search")
	var pageList []Page
	pageList, err := loadAllPages(search)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data := struct {
		List    []Page
		IsAdmin bool
	}{
		List:    pageList,
		IsAdmin: isAdmin,
	}
	err = templates.ExecuteTemplate(w, "all.html", data)
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

func login(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if password == "password" {
		session, _ := session_store.Get(r, "admin")
		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func loginScreen(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := session_store.Get(r, "admin")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
