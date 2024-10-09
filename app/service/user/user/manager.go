package user

import (
	"forum/app/models"
	"forum/app/repository"
)

type UserService interface {
	GetUserByToken(token string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByPostId(postId int) (int64, error)
	SendNotification(notification *models.Notification) error
	GetLikedCommentsByUserId(userId int64) ([]models.Comment, error)
	GetDislikedCommentsByUserId(userId int64) ([]models.Comment, error)
	GetAllNotificationsByUserId(userId int64) ([]models.Notification, error)
}

type userService struct {
	repository repository.UserQuery
}

func (u *userService) GetAllNotificationsByUserId(userId int64) ([]models.Notification, error) {
	return u.repository.GetAllNotifications(userId)
}

func NewUserService(repo repository.Repo) UserService {
	return &userService{repo.NewUserQuery()}
}
