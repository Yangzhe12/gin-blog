package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	var content [5]map[string]string
	for i := 0; i < 5; i++ {
		content[i] = map[string]string{
			"title": "article_" + strconv.Itoa(i),
			"user":  "user_" + strconv.Itoa(i),
		}
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"content": content,
		"page":    "首页",
	})
}
