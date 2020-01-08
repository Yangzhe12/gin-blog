package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginGet(c *gin.Context) {

	c.HTML(http.StatusOK, "account/login.html", gin.H{
		"page": "登陆",
	})
}

func LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if (username == "user1") && (password == "2012") {
		cookie, err := c.Cookie(username)
		if err != nil {
			c.SetCookie(username, "test_cookie", 3600, "/", "localhost", false, true)
		} else {
			fmt.Println("===================", cookie)
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}

}
