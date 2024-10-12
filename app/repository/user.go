package repository

import (
	"database/sql"
	"fmt"
	"forum/app/models"
	"log"
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
	NotificationExists(userTo, userFrom int64, sourceId int, action string) (bool, error)
	DeleteNotification(userTo, userFrom int64, sourceId int, action string) error
	GetLikedCommentsByUserId(userId int64) ([]models.Comment, error)
	GetDislikedCommentsByUserId(userId int64) ([]models.Comment, error)
	GetCommentByCommentId(commentId int) (models.Comment, error)
	GetAllNotifications(userId int64) ([]models.Notification, error)
	DeleteNotificationById(id int) error
}

type userQuery struct {
	db *sql.DB
}

func (u *userQuery) GetAllNotifications(userId int64) ([]models.Notification, error) {

	var notifications []models.Notification
	rows, err := u.db.Query("SELECT * FROM notifications WHERE UserTo=? order by CreatedAt desc ", userId)
	if err != nil {
		return notifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.Action,
			&notification.Content,
			&notification.UserFrom,
			&notification.UserTo,
			&notification.Username,
			&notification.SourceID,
			&notification.CreatedAt,
		); err != nil {
			fmt.Println(err, "err")
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	fmt.Println(notifications)
	return notifications, nil
}

func (u *userQuery) GetCommentByCommentId(commentId int) (models.Comment, error) {
	var comment models.Comment
	err := u.db.QueryRow(`select * from comments where comment_id=?`, commentId).Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Username, &comment.Message, &comment.Like, &comment.Dislike, &comment.Born)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}
func (u *userQuery) GetDislikedCommentsByUserId(userId int64) ([]models.Comment, error) {
	query := `select SourceId from notifications where UserFrom = ? and action='disliked your comment'`
	rows, err := u.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []models.Comment
	var commentsNumbers []int
	for rows.Next() {
		var commentNumber int
		if err := rows.Scan(&commentNumber); err != nil {
			fmt.Println("72")
			return nil, err
		}
		commentsNumbers = append(commentsNumbers, commentNumber)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, commentNumber := range commentsNumbers {
		comment, err := u.GetCommentByCommentId(commentNumber)
		if err != nil {
			fmt.Println("85")
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (u *userQuery) GetLikedCommentsByUserId(userId int64) ([]models.Comment, error) {
	query := `select SourceId from notifications where UserFrom = ? and action='liked your comment'`
	rows, err := u.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []models.Comment
	var commentsNumbers []int
	for rows.Next() {
		var commentNumber int
		if err := rows.Scan(&commentNumber); err != nil {
			fmt.Println("72")
			return nil, err
		}
		commentsNumbers = append(commentsNumbers, commentNumber)
	}
	fmt.Println(commentsNumbers, "comments Number")
	if err := rows.Err(); err != nil {
		fmt.Println("79")
		log.Fatal(err)
	}
	for _, commentNumber := range commentsNumbers {
		comment, err := u.GetCommentByCommentId(commentNumber)
		if err != nil {
			fmt.Println("85")
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (u *userQuery) Notification(notification *models.Notification) error {
	_, err := u.db.Exec(`insert into notifications (action,content,UserFrom,UserTo,Username,SourceId,CreatedAt) values (?,?,?,?,?,?,?)`, notification.Action, notification.Content, notification.UserFrom, notification.UserTo, notification.Username, notification.SourceID, notification.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) NotificationExists(userTo, userFrom int64, sourceId int, action string) (bool, error) {
	var exists bool
	err := u.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM notifications WHERE UserTo = ? AND UserFrom = ? AND SourceID = ? AND Action = ?)`, userTo, userFrom, sourceId, action).Scan(&exists)
	if err != nil {
		log.Println("Error checking if notification exists:", err)
		return false, err
	}

	log.Printf("Notification exists: %v for action: %s\n", exists, action)
	return exists, nil
}

func (u *userQuery) DeleteNotification(userTo, userFrom int64, sourceId int, action string) error {
	log.Printf("Attempting to delete notification with action: %s, UserTo: %d, UserFrom: %d, SourceID: %d\n", action, userTo, userFrom, sourceId)

	res, err := u.db.Exec(`DELETE FROM notifications WHERE UserTo = ? AND UserFrom = ? AND SourceID = ? AND Action = ?`, userTo, userFrom, sourceId, action)
	if err != nil {
		log.Println("Error deleting notification:", err)
		return err
	}
	// to check is it deleted or not
	rowsAffected, _ := res.RowsAffected()
	log.Printf("Deleted %d rows for action %s\n", rowsAffected, action)
	return nil
}
func (u *userQuery) DeleteNotificationById(id int) error {
	res, err := u.db.Exec(`DELETE FROM notifications WHERE id=?`, id)
	if err != nil {
		log.Println("Error deleting notification:", err)
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	log.Printf("Deleted %d rows for action %s\n", rowsAffected)
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
