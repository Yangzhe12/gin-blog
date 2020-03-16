package models

// Article结构体
type Article struct {
	ID          int
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ArticleTag  []string `json:"tags"`
	PubDatetime []uint8
	UpdDatetime []uint8
	AuthorName  string
	Pageview    int
	LikeNum     int
}

// Tag 结构体
type Tag struct {
	ID      int
	TagName string
}

// LikedData 点赞数据结构体
type LikedData struct {
	CurrentUsername string `json:"currentUser"`    // 当前登陆的用户名
	LikedArticleID  string `json:"likedArticleID"` // 点赞的文章id
	LikedStatus     string `json:"likedStatus"`    // 当前点赞状态
}
