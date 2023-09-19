package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/pages/home.tmpl",
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("hello from snippet box"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Displaying specific snippet"))
	fmt.Fprintf(w, "Displaying a specific snippet with id %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))
		// Instead of using w.WriteHeader() and w.Write() functions we can use
		// http.Error() shortcut which takes a given message and status code and
		// then calls the w.WriteHeader() and w.Write() methods behind the scenes

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Creating a snippet"))
}
