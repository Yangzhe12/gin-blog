package controllers

import (
	"fmt"
	"gin-blog/utils"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ArticleGet(c *gin.Context) {
	var dbTitle, dbContent, dbAuthor string
	var dbArticleID int
	var dbPubDatetime, dbUpdDatetime []uint8
	articleID, err := strconv.Atoi(c.Param("articleID"))
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusNotFound, "找不到文章！")
		return
	}
	currentUser := utils.GetUserInfo(c)
	queryArtSQL := "select id,title,content,pub_datetime,upd_datetime,author_name from article where id=?;"
	row := utils.Db.QueryRow(queryArtSQL, articleID)
	if row == nil {
		c.String(http.StatusNotFound, "找不到文章！")
	} else {
		err = row.Scan(&dbArticleID, &dbTitle, &dbContent, &dbPubDatetime, &dbUpdDatetime, &dbAuthor)
		// 让字符串不转义为html格式
		articleCotent := template.HTML(dbContent)
		c.HTML(http.StatusOK, "article/article.html", gin.H{
			"page":          "文章内容",
			"username":      currentUser,
			"articleID":     dbArticleID,
			"title":         dbTitle,
			"articleCotent": articleCotent,
			"pubDateTime":   utils.B2S(dbPubDatetime),
			"updDateTime":   utils.B2S(dbUpdDatetime),
			"author":        dbAuthor,
		})
	}
}

// 不转义html
func unescaped(x string) interface{} {
	return template.HTML(x)
}
