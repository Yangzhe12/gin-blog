package controllers

import (
	"fmt"
	config "gin-blog/conf"
	"gin-blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	var dbArtIDStr, dbTitle, dbContent, dbAuthor string
	var dbArticleID int
	var dbDatetime []uint8
	var respData []map[string]string
	var currentUser string
	var articlesNumber int
	var currentPage int
	currentPage, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		fmt.Println("类型转换错误")
		currentPage = 1
	}

	utils.SessionKey = config.GetConfiguration().UserInfoSessionKey
	session := utils.Default(c)
	username := session.Get("username")
	if username != nil {
		currentUser = username.(string)
	} else {
		currentUser = ""
	}

	articlesNumberSQL := "select count(title) from article"
	row := utils.Db.QueryRow(articlesNumberSQL)
	err = row.Scan(&articlesNumber)
	if err != nil {
		fmt.Println(err)
	}
	totalPage := (articlesNumber / config.ArticlesPerPage) + 1

	indexSQL := fmt.Sprintf("select id,title,content,upd_datetime,author_name from article order by upd_datetime desc limit %d,%d;", (currentPage-1)*config.ArticlesPerPage, config.ArticlesPerPage)
	rows, err := utils.Db.Query(indexSQL)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&dbArticleID, &dbTitle, &dbContent, &dbDatetime, &dbAuthor)
		dbArtIDStr = strconv.Itoa(dbArticleID)
		respData = append(respData, map[string]string{
			"articleID":   dbArtIDStr,
			"title":       dbTitle,
			"artContent":  dbContent,
			"pubDateTime": utils.B2S(dbDatetime),
			"author":      dbAuthor,
		})
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"content":     respData,
		"page":        "首页",
		"username":    currentUser,
		"totalPage":   totalPage,
		"currentPage": currentPage,
	})
}
