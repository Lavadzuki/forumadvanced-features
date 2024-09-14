package user

import (
	"fmt"
	"forum/app/models"
	"log"
	"time"
)

func (u *userService) GetUserByToken(token string) (models.User, error) {
	userId, err := u.repository.GetUserIdByToken(token)
	if err != nil {
		return models.User{}, err
	}
	user, err := u.repository.GetUserByUserId(userId)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u userService) GetUserByEmail(email string) (models.User, error) {
	return u.repository.GetUserByEmail(email)
}

func (u *userService) GetUserByPostId(postId int) (int64, error) {
	return u.repository.GetUserIdByPostId(postId)
}

func (u *userService) SendNotification(userTo, userFrom int64, userFromUsername string, sourceId int, action string) error {
	// First, we need to identify the opposite action. There can becas where User Liked -> Disliked -> Liked. If we do not delete old Like, the last action will be Disliked since we do not create new insctance of like if it is already exists(Explained later.)

	var oppositeAction string
	if action == "liked your post" {
		oppositeAction = "disliked your post"
	} else if action == "disliked your post" {
		oppositeAction = "liked your post"
	} else if action == "disliked your comment" {
		oppositeAction = "liked your comment"
	} else if action == "liked your comment" {
		oppositeAction = "disliked your comment"
	} else {
		log.Println("Unknown action for the notification:", action)
	}

	// If user of the same post making Dislike instead of Like we should delete this dislike)
	err := u.repository.DeleteNotification(userTo, userFrom, sourceId, oppositeAction)
	if err != nil {
		log.Println("Error deleting opposite action notification:", err)
		return err
	}
	// if User two or more times in a row click on the Like, we check is it exists.
	exists, err := u.repository.NotificationExists(userTo, userFrom, sourceId, action)
	if err != nil {
		log.Println("Error checking if notification exists:", err)
		return err
	}

	//  If a notification for the same action exists, don't send it again
	if exists {
		log.Println("Notification for this action already exists, not sending again.")
		return nil
	}

	notification := models.Notification{
		Action:    action,
		Content:   fmt.Sprintf("%s %s ", userFromUsername, action),
		UserFrom:  userFrom,
		UserTo:    userTo,
		Username:  userFromUsername,
		SourceID:  sourceId,
		CreatedAt: time.Now(),
	}

	err = u.repository.Notification(&notification)
	if err != nil {
		log.Println("Error creating new notification:", err)
		return err
	}

	return nil
}
