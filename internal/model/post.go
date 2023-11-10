package model

type Post struct {
	Model
	Title string `json:"title"`
	Text  string `json:"text"`
}
