package model

import "github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"

type Comment struct {
	BaseModel
	Author string `json:"author"`
	PostID uint64 `json:"post_id"`
	UserID uint64 `json:"user_id"`
	Text   string `json:"text"`
}

func (p Comment) ToDTO() *dto.CommentResponse {
	return &dto.CommentResponse{
		Author:    p.Author,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		DeletedAt: p.DeletedAt,
		ID:        p.ID,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy,
		Text:      p.Text,
		PostID:    p.PostID,
		UserID:    p.UserID,
	}
}
