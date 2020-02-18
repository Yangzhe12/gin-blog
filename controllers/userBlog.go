package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserBlogGet 处理用户文章列表请求
func UserBlogGet(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		fmt.Println("用户不存在")
		c.Redirect(http.StatusMovedPermanently, "/index")
		return
	}
	GetArticleList(c, username)
}
