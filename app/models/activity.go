package models

type ActivityData struct {
	User              User
	CreatedPosts      []Post
	LikedPosts        []Post
	DislikedPosts     []Post
	CommentsWithPosts []CommentWithPost
}
