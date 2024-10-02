package handlers

import (
	"fmt"
	"forum/app/models"
	"forum/pkg"
	"log"
	"net/http"
)

func (app *App) ActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/activity" {
		pkg.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	user, err := app.userService.GetUserByToken(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}
	fmt.Println(user, "user")

	userID := user.ID

	createdPosts, err := app.postService.GetPostsByUserID(userID)
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	likedPosts, err := app.postService.GetLikedPostsByUserID(userID)
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	dislikedPosts, err := app.postService.GetDislikedPostsByUserID(userID)
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	commentsWithPosts, err := app.postService.GetCommentsByUserID(userID)
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	likedComments, err := app.userService.GetLikedCommentsByUserId(userID)
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	dislikedComments, err := app.userService.GetDislikedCommentsByUserId(userID)
	fmt.Println(dislikedComments, "disliked comments")
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	activityData := models.ActivityData{
		User:              user,
		CreatedPosts:      createdPosts,
		LikedPosts:        likedPosts,
		DislikedPosts:     dislikedPosts,
		LikedComments:     likedComments,
		DislikedComments:  dislikedComments,
		CommentsWithPosts: commentsWithPosts,
	}

	pkg.RenderTemplate(w, "activity.html", activityData)
}
