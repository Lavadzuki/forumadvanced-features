package post

import (
	"database/sql"
	"errors"
	"forum/app/models"
	"forum/app/repository"
	"net/http"
)

type PostService interface {
	// reaction.go
	LikePost(postId, userId int) int
	DislikePost(postId, userId int) int
	LikeComment(commentId, userId int) int
	DislikeComment(commentId, userId int) int
	// post.go
	GetAllPosts() ([]models.Post, error)
	CreatePost(post *models.Post) (int, error)
	GetAllCommentsAndPostsByPostId(id int64) (models.Post, int)
	CreateComment(comment *models.Comment) (int, error)
	// filter.go
	GetFilterPosts(genre string, user models.User) (models.Data, int)
	GetWelcomeFilterPosts(genre string) (models.Data, int)

	GetPostsByUserID(userID int64) ([]models.Post, error)
	GetLikedPostsByUserID(userID int64) ([]models.Post, error)
	GetDislikedPostsByUserID(userID int64) ([]models.Post, error)
	GetCommentsByUserID(userID int64) ([]models.CommentWithPost, error)

	DeletePost(postId int) error
	GetPostByPostId(postId int) (models.Post, error)
	GetAllCategory() ([]models.Category, error)
	UpdatePost(post models.Post) (int, error)
}

type postService struct {
	repository repository.PostQuery
}

func (p postService) UpdatePost(post models.Post) (int, error) {
	status, err := p.repository.UpdatePost(post)

	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("post not found")
		}
		return http.StatusInternalServerError, err
	}
	return status, nil
}

func (p postService) GetAllCategory() ([]models.Category, error) {
	return p.repository.GetCategory()
}

func (p postService) GetPostByPostId(postId int) (models.Post, error) {
	id := int64(postId)
	post, err := p.repository.GetPostById(id)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (p postService) DeletePost(postId int) error {
	err := p.repository.DeletePost(postId)
	return err

}

func NewPostService(repo repository.Repo) PostService {
	return &postService{repository: repo.NewPostQuery()}
}

func (s *postService) GetPostsByUserID(userID int64) ([]models.Post, error) {
	return s.repository.GetPostsByUserID(userID)
}

func (s *postService) GetLikedPostsByUserID(userID int64) ([]models.Post, error) {
	return s.repository.GetLikedPostsByUserID(userID)
}

func (s *postService) GetDislikedPostsByUserID(userID int64) ([]models.Post, error) {
	return s.repository.GetDislikedPostsByUserID(userID)
}

func (s *postService) GetCommentsByUserID(userID int64) ([]models.CommentWithPost, error) {
	return s.repository.GetCommentsByUserID(userID)
}
