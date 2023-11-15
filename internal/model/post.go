package model

type Post struct {
	BaseModel
	Title string `json:"title"`
	Body  string `json:"body"`
}
