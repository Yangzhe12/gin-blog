package controllers

import (
	"fmt"
	"net/http"

	"gin-blog/utils"

	"gin-blog/models.go"

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

	sqlStr := "select username,password from account where username = ?"
	row := utils.Db.QueryRow(sqlStr, username)
	user := &models.Account{}
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusOK, "account/login.html", gin.H{
			"msg": "用户名或密码错误！",
		})
	} else {
		if (username == user.Username) && (password == user.Password) {
			cookie, err := c.Cookie(username)
			if err != nil {
				cookie = "test_cookie"
				c.SetCookie(username, cookie, 3600, "/", "localhost", false, true)
			}
			c.Redirect(http.StatusMovedPermanently, "/")
		} else {
			c.HTML(http.StatusOK, "account/login.html", gin.H{
				"msg": "用户名或密码错误！",
			})
		}
	}

}
