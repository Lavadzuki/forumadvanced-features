package models

import "time"

type Notification struct {
	ID        int
	Action    string
	Content   string
	UserFrom  int64
	UserTo    int64
	Username  string
	Source    string
	SourceID  int
	CreatedAt time.Time
}
