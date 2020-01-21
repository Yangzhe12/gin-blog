package controllers

import (
	"fmt"
	config "gin-blog/conf"
	"gin-blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	utils.SessionKey = config.GetConfiguration().UserInfoSessionKey
	session := utils.Default(c)
	username := session.Get("username")
	var currentUser string
	if username != nil {
		currentUser = username.(string)
	} else {
		currentUser = ""
	}
	var DbTitle, DbContent string
	var DbDatetime []uint8
	var respData []map[string]string

	indexSQL := fmt.Sprintf("select title,content,pub_datetime from article order by pub_datetime desc limit %d;", config.ArticlesPerPage)
	rows, err := utils.Db.Query(indexSQL)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&DbTitle, &DbContent, &DbDatetime)
		respData = append(respData, map[string]string{
			"title":       DbTitle,
			"artContent":  DbContent,
			"pubDateTime": utils.B2S(DbDatetime),
		})
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"content":  respData,
		"page":     "首页",
		"username": currentUser,
	})
}
