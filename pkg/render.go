package pkg

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

var (
	templateCache = make(map[string]*template.Template)
	once          sync.Once
)

func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	once.Do(func() {
		err := createTemplate()
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	})
	fmt.Println(7)

	t, ok := templateCache[templateName]
	if !ok {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

func createTemplate() error {
	pages, err := filepath.Glob("./templates/html/*.html")
	if err != nil {
		return err
	}

	standalonePages := map[string]bool{
		"signin.html":        true,
		"signup.html":        true,
		"createpost.html":    true,
		"commentunauth.html": true,
		"commentview.html":   true,
		"error.html":         true,
		"edit_post.html":     true,
	}

	for _, page := range pages {
		name := filepath.Base(page)

		var ts *template.Template
		var err error

		if standalonePages[name] {
			ts, err = template.ParseFiles(page)
		} else {
			ts, err = template.ParseFiles("templates/html/base.html", page)
		}

		if err != nil {
			return err
		}

		templateCache[name] = ts
	}

	return nil
}
