package models

type TmplData struct {
	Post       Post
	Categories []string
	User       User // если требуется информация о пользователе
}
