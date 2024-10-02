package models

type ActivityData struct {
	User              User
	CreatedPosts      []Post
	LikedPosts        []Post
	DislikedPosts     []Post
	LikedComments     []Comment
	DislikedComments  []Comment
	CommentsWithPosts []CommentWithPost
}
