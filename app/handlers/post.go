package handlers

import (
	"forum/app/models"
	"forum/pkg"
	"log"
	"net/http"
	"time"
)

func (app *App) PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pkg.RenderTemplate(w, "createpost.html", models.Data{})
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			return
		}
		title := r.FormValue("title")
		message := r.FormValue("message")
		genre := r.Form["category"]

		if len(genre) == 0 {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		user, ok := r.Context().Value(KeyUserType(keyUser)).(models.User)
		if !ok {
			pkg.ErrorHandler(w, http.StatusUnauthorized)
			return
		}
		post := models.Post{
			Title:       title,
			Content:     message,
			Category:    models.Stringslice(genre),
			Author:      user,
			CreatedTime: time.Now().Format(time.RFC822),
		}
		status, err := app.postService.CreatePost(&post)
		if err != nil {
			log.Println(err)
			switch status {
			case http.StatusInternalServerError:
				pkg.ErrorHandler(w, http.StatusInternalServerError)
				return
			case http.StatusBadRequest:
				pkg.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)

	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *App) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		// Extract post ID from URL or query parameters
		vars := mux.Vars(r) // Assuming you're using Gorilla Mux for routing
		postID, ok := vars["id"]
		if !ok {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}

		// Call the DeletePost service with the post ID
		err := app.postService.DeletePost(postID)
		if err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}

		// If successful, respond with 204 No Content
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
