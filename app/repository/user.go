package repository

import (
	"database/sql"
	"forum/app/models"
	"strings"
)

type UserQuery interface {
	CreateUser(user *models.User) error
	GetUserIdByToken(token string) (int, error)
	GetUserByUserId(userId int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	UpdateUser(user *models.User) error
	GetUserIdByPostId(postId int) (int64, error)
	Notification(notification *models.Notification) error
}

type userQuery struct {
	db *sql.DB
}

func (u *userQuery) Notification(notification *models.Notification) error {
	_, err := u.db.Exec(`insert into notifications (action,content,UserFrom,UserTo,Username,SourceId,CreatedAt) values (?,?,?,?,?,?,?)`, notification.Action, notification.Content, notification.UserFrom, notification.UserTo, notification.Username, notification.SourceID, notification.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) CreateUser(user *models.User) error {
	_, err := u.db.Exec(`insert into users (username, email, password) VALUES (?,?,?)`, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) GetUserIdByPostId(postId int) (int64, error) {
	var id int64
	err := u.db.QueryRow("SELECT user_id FROM posts WHERE post_id = ?", postId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *userQuery) UpdateUser(user *models.User) error {
	query := "Update users set username=?,email=? where user_id=?"
	_, err := u.db.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}
func (u *userQuery) GetUserIdByToken(token string) (int, error) {
	row := u.db.QueryRow(`select user_id from sessions where token=?`, token)
	var userId int
	err := row.Scan(&userId)
	if err != nil {
		return -1, err
	}
	return userId, nil
}

func (u *userQuery) GetUserByUserId(userId int) (models.User, error) {
	row := u.db.QueryRow(`select user_id, email, password,username from users where user_id=?`, userId)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userQuery) GetUserByEmail(email string) (models.User, error) {
	mail := strings.ToLower(email)
	row := u.db.QueryRow(`select user_id,email,password,username from users where email=?`, mail)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userQuery) GetUserByUsername(username string) (models.User, error) {
	userName := strings.ToLower(username)
	row := u.db.QueryRow(`select user_id,email,password,username from users where username=?`, userName)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
