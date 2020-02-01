package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ArticleGet(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
