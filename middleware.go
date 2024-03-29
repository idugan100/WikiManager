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
