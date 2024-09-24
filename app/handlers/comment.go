package handlers

import (
	"fmt"
	"forum/app/models"
	"forum/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *App) CommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		parts := strings.Split(r.URL.Path, "/")
		id, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Println(err)
			pkg.ErrorHandler(w, http.StatusNotFound)
			return
		}
		initialPost, status := app.postService.GetAllCommentsAndPostsByPostId(int64(id))
		switch status {
		case http.StatusMethodNotAllowed:
			pkg.ErrorHandler(w, http.StatusInternalServerError)
			return
		case http.StatusBadRequest:
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		data := models.Data{
			Comment:     initialPost.Comment,
			InitialPost: initialPost,
		}

		pkg.RenderTemplate(w, "commentview.html", data)
	case http.MethodPost:
		parts := strings.Split(r.URL.Path, "/")
		id, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Println(err)
			pkg.ErrorHandler(w, http.StatusNotFound)
			return
		}
		message := r.FormValue("comment")
		path := "/post/comment/" + parts[3]
		user, ok := r.Context().Value(KeyUserType(keyUser)).(models.User)
		if !ok {
			pkg.ErrorHandler(w, http.StatusUnauthorized) // http.StatusUnauthorized-status unauthorized
			return
		}
		comment := models.Comment{
			PostId:   int64(id),
			UserId:   user.ID,
			Username: user.Username,
			Message:  message,
			Born:     time.Now().Format(time.RFC822),
		}
		userTo, err := app.userService.GetUserByPostId(id)
		if err != nil {
			log.Println(err)
			pkg.ErrorHandler(w, http.StatusInternalServerError)
		}

		userFrom := user.ID
		userFromUsername := user.Username
		sourceId := id
		action := "commented your post"
		content := fmt.Sprintf("%s %s", user.ID, action)
		notification := models.Notification{
			Action:    action,
			Content:   content,
			UserFrom:  userFrom,
			UserTo:    userTo,
			Username:  userFromUsername,
			Source:    "Comment",
			SourceID:  sourceId,
			CreatedAt: time.Now(),
		}
		status, err := app.postService.CreateComment(&comment)
		err = app.userService.SendNotification(&notification)
		if err != nil {
			log.Println(err)
		}
		switch status {
		case http.StatusInternalServerError:
			pkg.ErrorHandler(w, http.StatusInternalServerError)
			return
		case http.StatusBadRequest:
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		case http.StatusOK:
			http.Redirect(w, r, path, http.StatusFound)
		}
	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}

func (app *App) WelcomeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	initialPost, status := app.postService.GetAllCommentsAndPostsByPostId(int64(id))
	switch status {
	case http.StatusInternalServerError:
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	case http.StatusBadRequest:
		pkg.ErrorHandler(w, http.StatusBadRequest)
		return
	}
	data := models.Data{
		Comment:     initialPost.Comment,
		InitialPost: initialPost,
	}
	pkg.RenderTemplate(w, "commentunauth.html", data)
}
