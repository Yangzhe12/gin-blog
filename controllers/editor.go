package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditorGet(c *gin.Context) {

	c.HTML(http.StatusOK, "editor/editor.html", gin.H{
		"page": "写文章",
	})
}
