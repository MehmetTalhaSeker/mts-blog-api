package model

type Post struct {
	BaseModel
	Title string `json:"title"`
	Text  string `json:"text"`
}
