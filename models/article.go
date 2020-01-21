package models

// Article结构体
type Article struct {
	ID          int
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ArticleTag  []string `json:"tags"`
	PubDatetime string
	UpdDatetime string
	AuthorID    int
}

// Tag 结构体
type Tag struct {
	ID      int
	TagName string
}
