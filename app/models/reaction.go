package models

type Reaction struct {
	PostId int64
	UserId int64
	Status int // 1 for like 0 for dislike
}
