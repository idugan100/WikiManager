package main

import "net/http"

// middleware to ensure all paths are in current dir
func validatePath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("path")
		if fileValidator.MatchString(filename) {
			fn(w, r, filename)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func requireAdmin(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := session_store.Get(r, "admin")

		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth {
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		} else {
			fn(w, r)
		}
	}
}

func isAdmin(fn func(http.ResponseWriter, *http.Request, bool)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := session_store.Get(r, "admin")

		auth, ok := session.Values["authenticated"].(bool)

		if !auth || !ok {
			fn(w, r, false)
		} else {
			fn(w, r, true)
		}
	}
}

func isAdminAndValidatePath(fn func(http.ResponseWriter, *http.Request, string, bool)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		filename := r.PathValue("path")
		if fileValidator.MatchString(filename) {

			session, _ := session_store.Get(r, "admin")
			auth, ok := session.Values["authenticated"].(bool)

			if !auth || !ok {
				fn(w, r, filename, false)
			} else {
				fn(w, r, filename, true)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
