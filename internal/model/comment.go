package model

import (
	"time"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

type Comment struct {
	Model
	PostID string `json:"postID"`
	Text   string `json:"text"`
}

func NewComment(PostID, Text string) *Comment {
	return &Comment{
		Model: Model{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			DeletedAt: nil,
			Status:    types.Active,
			CreatedBy: "",
			UpdatedBy: "",
			DeletedBy: "",
		},
		PostID: PostID,
		Text:   Text,
	}
}
