package model

type Comment struct {
	BaseModel
	PostID string `json:"post_id"`
	Text   string `json:"text"`
}
