package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	// 检查是否登陆
	// fmt.Println("****************1")
	// username, _ := utils.IsLogin(c)
	// fmt.Println("****************2")
	username := "aaa"
	var content [5]map[string]string
	for i := 0; i < 5; i++ {
		content[i] = map[string]string{
			"title": "article_" + strconv.Itoa(i),
			"user":  "user_" + strconv.Itoa(i),
		}
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"content":  content,
		"page":     "首页",
		"username": username,
	})
}
