package controllers

import (
	config "gin-blog/conf"
	"gin-blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutGet(c *gin.Context) {
	utils.SessionKey = config.GetConfiguration().UserInfoSessionKey
	session := utils.Default(c)
	session.Delete("username")
	session.Delete("userID")
	session.Save()
	c.Header("Cache-Control", "no-cache")
	c.Redirect(http.StatusFound, "/")
}
