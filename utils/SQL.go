package utils

const (
	// IndexSQL 首页文章列表查询SQL
	IndexSQL = "select id,title,content,pageview,pub_datetime,author_name, like_num from article order by upd_datetime desc limit %d,%d;"

	// UserArtListSQL 用户文章列表页查询SQL
	UserArtListSQL = "select id,title,content,pageview,pub_datetime,author_name, like_num from article where author_name=? order by upd_datetime desc limit %d,%d;"

	// QueryArtByIDSQL 通过文章ID来查询文章的详细信息
	QueryArtByIDSQL = "select title,content,pageview,pub_datetime,author_name from article where id=?;"

	// CountArtNumberSQL 统计数据库中文章的总数量
	CountArtNumberSQL = "select count(title) from article"

	// AddArticleSQL 添加文章
	AddArticleSQL = "insert into article (title, content, author_name) values (?,?,?)"

	// AddArtPageviewSQL 增加文章访问量
	AddArtPageviewSQL = "update article set pageview=pageview+1 where id=?;"
)
