package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutGet(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
}
