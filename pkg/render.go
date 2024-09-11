package pkg

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"forum/app/models"
)

var (
	templateCache = make(map[string]*template.Template)
	once          sync.Once
)

func RenderTemplate(w http.ResponseWriter, template string, data models.Data) {
	once.Do(func() {
		err := createTemplate()
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	})

	t, ok := templateCache[template]
	if !ok {

		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

}

func createTemplate() error {
	// Load all individual page templates
	pages, err := filepath.Glob("./templates/html/*.html")
	if err != nil {
		return err
	}

	// Define the pages that do NOT need the base layout
	standalonePages := map[string]bool{
		"signin.html":        true,
		"signup.html":        true,
		"createpost.html":    true,
		"commentunauth.html": true,
		"commentview.html":   true,
		"error.html":         true,
	}

	for _, page := range pages {
		name := filepath.Base(page)

		var ts *template.Template
		var err error

		// Check if the page is a standalone page
		if standalonePages[name] {
			// Parse without the base layout
			ts, err = template.ParseFiles(page)
		} else {
			// Parse with the base layout
			ts, err = template.ParseFiles("templates/html/base.html", page)
		}

		if err != nil {
			return err
		}

		// Cache the parsed template
		templateCache[name] = ts
	}

	return nil
}
