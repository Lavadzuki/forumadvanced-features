package models

type TmplData struct {
	Post       Post
	Categories []string
	User       User            // User information, if needed
	IsSelected map[string]bool // A map to track selected categories
}
