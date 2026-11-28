package domain

import (
	"errors"
	"time"
)

// ErrCommentEmpty empty comment
var ErrCommentEmpty = errors.New("el comentario no puede estar vacío")

// Comment represent a comment
type Comment struct {
	ID        string
	TaskID    string
	AuthorID  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewComment constructor
func NewComment(taskID, authorID, content string) (*Comment, error) {
	if content == "" {
		return nil, ErrCommentEmpty
	}

	now := time.Now()

	return &Comment{
		TaskID:    taskID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
