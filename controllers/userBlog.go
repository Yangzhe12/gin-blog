package controllers

import (
	"fmt"
	config "gin-blog/conf"
	"gin-blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserBlogGet(c *gin.Context) {
	var dbArtIDStr, dbTitle, dbContent, dbAuthor string
	var dbArticleID int
	var dbUpdDatetime []uint8
	var respData []map[string]string
	var currentUser string
	var articlesNumber int
	var currentPage, totalPage int

	username := c.Param("username")
	if username == "" {
		fmt.Println("用户不存在")
		c.Redirect(http.StatusMovedPermanently, "/index")
		return
	}
	currentUser = utils.GetUserInfo(c)

	articlesNumberSQL := "select count(title) from article where author_name=? "
	row := utils.Db.QueryRow(articlesNumberSQL, username)
	err := row.Scan(&articlesNumber)
	if err != nil {
		fmt.Println(err)
		return
	}
	if articlesNumber == 0 {
		totalPage = 0
		currentPage = 0
		respData = nil
	} else {
		currentPage, err = strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			fmt.Println("类型转换错误")
			currentPage = 1
		}
		totalPage = (articlesNumber / config.ArticlesPerPage) + 1
		indexSQL := fmt.Sprintf("select id,title,content,upd_datetime,author_name from article where author_name=? order by upd_datetime desc limit %d,%d;", (currentPage-1)*config.ArticlesPerPage, config.ArticlesPerPage)
		rows, err := utils.Db.Query(indexSQL, username)
		if err != nil {
			fmt.Println(err)
			respData = nil
		} else {
			for rows.Next() {
				rows.Scan(&dbArticleID, &dbTitle, &dbContent, &dbUpdDatetime, &dbAuthor)
				dbArtIDStr = strconv.Itoa(dbArticleID)
				respData = append(respData, map[string]string{
					"articleID":   dbArtIDStr,
					"title":       dbTitle,
					"artContent":  dbContent,
					"updDateTime": utils.B2S(dbUpdDatetime),
					"author":      dbAuthor,
				})
			}
		}
	}
	c.HTML(http.StatusOK, "blog/blog.html", gin.H{
		"content":     respData,
		"page":        "首页",
		"username":    username,
		"totalPage":   totalPage,
		"currentPage": currentPage,
		"currentUser": currentUser,
	})
}
