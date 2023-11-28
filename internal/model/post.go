package model

import "github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"

type Post struct {
	BaseModel
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (p Post) ToDTO() *dto.PostResponse {
	return &dto.PostResponse{
		Body:      p.Body,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		DeletedAt: p.DeletedAt,
		ID:        p.ID,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy,
		Title:     p.Title,
	}
}
