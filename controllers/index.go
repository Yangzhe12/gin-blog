package controllers

import (
	"github.com/gin-gonic/gin"
)

// IndexGet 处理首页Get请求
func IndexGet(c *gin.Context) {
	GetArticleList(c, "")
}
