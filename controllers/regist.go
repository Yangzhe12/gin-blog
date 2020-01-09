package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegistGet(c *gin.Context) {
	c.HTML(http.StatusOK, "account/regist.html", gin.H{
		"page": "注册",
	})
}

func RegistPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	repwd := c.PostForm("repwd")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	fmt.Println("--------------", username, password, repwd, email, phone)
	c.Redirect(http.StatusMovedPermanently, "/login")
}
