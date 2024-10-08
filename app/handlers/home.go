package handlers

import (
	"forum/app/models"
	"forum/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		pkg.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	posts, err := app.postService.GetAllPosts()
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	data := models.Data{
		Posts: posts,
		// User:  user,
		Genre: "/",
	}

	pkg.RenderTemplate(w, "index.html", data)
}

func (app *App) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	posts, err := app.postService.GetAllPosts()
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	data := models.Data{
		Posts: posts,
	}
	pkg.RenderTemplate(w, "welcome.html", data)
}

func (app *App) DeletePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[3])

	err = app.postService.DeletePost(id)
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/filter/created-post", http.StatusFound)
}

func (app *App) EditPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		// Получаем ID поста из URL
		parts := strings.Split(r.URL.Path, "/")
		postID, err := strconv.Atoi(parts[3])
		if err != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		// Получаем пост по ID
		post, err := app.postService.GetPostByPostId(postID)
		if err != nil {
			pkg.ErrorHandler(w, http.StatusNotFound)
			return
		}

		// Формируем данные для шаблона
		tmplData := models.TmplData{
			Post:       post,
			Categories: []string{"Fantasy", "Drama", "Comedy", "Adventure", "Romance"},
			User:       models.User{}, // Вы можете передать реального пользователя
		}

		// Рендерим HTML-шаблон для редактирования поста
		pkg.RenderTemplate(w, "edit_post.html", tmplData)

		return

	}

	if r.Method == http.MethodPost {
		// Обрабатываем данные формы
		err := r.ParseForm()
		if err != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
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

		// Получаем ID поста из URL
		parts := strings.Split(r.URL.Path, "/")
		postID, err := strconv.Atoi(parts[3])
		if err != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		// Обновляем пост
		changedPost := models.Post{
			Id:          int64(postID),
			Title:       title,
			Content:     message,
			Category:    models.Stringslice(genre),
			Author:      user,
			CreatedTime: time.Now().Format(time.RFC822),
		}

		status, err := app.postService.UpdatePost(changedPost)
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

	}
}
