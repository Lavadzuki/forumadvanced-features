package user

import (
	"fmt"
	"forum/app/models"
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

	notification := models.Notification{
		Action:    action,
		Content:   fmt.Sprintf("%s %s ", userFromUsername, action),
		UserFrom:  userFrom,
		UserTo:    userTo,
		Username:  userFromUsername,
		SourceID:  sourceId,
		CreatedAt: time.Now(),
	}
	err := u.repository.Notification(&notification)
	if err != nil {
		return err
	}
	return nil
}
