package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"text/template"

	"github.com/satish218/sinppetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		// http.NotFound(w, r)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/pages/home.tmpl",
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		// app.errorLog.Print(err.Error())
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Create an instance of a templateData struct holding the slice of
	// snippets.
	data := &templateData{Snippets: snippets}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		// app.errorLog.Print(err.Error())
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("hello from snippet box"))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.

	files := []string{
		"./ui/html/pages/view.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/base.tmpl",
	}

	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
	}
	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{
		Snippet: snippet,
	}
	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
	// Write the snippet data as a plain-text HTTP response body.
	//fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))
		// Instead of using w.WriteHeader() and w.Write() functions we can use
		// http.Error() shortcut which takes a given message and status code and
		// then calls the w.WriteHeader() and w.Write() methods behind the scenes
		app.clientError(w, http.StatusMethodNotAllowed)
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n -Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	w.Write([]byte("Creating a snippet"))
}
