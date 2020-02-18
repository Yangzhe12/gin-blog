package controllers

import (
	"database/sql"
	"fmt"
	config "gin-blog/conf"
	"gin-blog/models"
	"gin-blog/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	utrackcsrf "github.com/utrack/gin-csrf"
)

func LoginGet(c *gin.Context) {
	utils.SessionKey = config.GetConfiguration().UserInfoSessionKey
	session := utils.Default(c)
	v := session.Get("username")
	if v != nil {
		fmt.Println("------------------------1")
		c.Redirect(http.StatusMovedPermanently, "/v1")
	} else {
		fmt.Println("------------------------2")
		csrfToken := utrackcsrf.GetToken(c)
		fmt.Println("------------------------3")
		c.HTML(http.StatusOK, "account/login.html", gin.H{
			"page":      "登陆",
			"csrfToken": csrfToken,
		})
	}
}

func LoginPost(c *gin.Context) {
	// 获取请求中的JSON数据
	var requestData models.Account
	if err := c.ShouldBindJSON(&requestData); err == nil {

		// 查询用户数据库，校验用户是否存在
		loginSQL := "select username,password from account where username = ?"
		row := utils.Db.QueryRow(loginSQL, requestData.Username)
		var dbUsername, dbPassword string

		err = row.Scan(&dbUsername, &dbPassword)
		if err == nil {
			// 验证用户名密码
			if (requestData.Username == dbUsername) && (requestData.Password == dbPassword) {
				utils.SessionKey = config.GetConfiguration().UserInfoSessionKey
				session := utils.Default(c)
				session.Options(sessions.Options{
					MaxAge: config.GetConfiguration().UserInfoSessionValidTime,
				})
				value := session.Get("username")
				if value == nil {
					session.Set("username", requestData.Username)
					session.Save()
				}

				c.JSON(http.StatusOK, gin.H{
					"resno":    0,
					"username": requestData.Username,
					"msg":      "登陆成功！",
				})

			} else {
				c.JSON(http.StatusOK, gin.H{
					"resno": 2,
					"msg":   "用户名或密码错误！",
				})
			}
		} else {
			// 未查询到用户信息，用户名不存在
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, gin.H{
					"resno": 3,
					"msg":   "用户名不存在！",
				})
			} else {
				// 数据库连接出错
				c.JSON(http.StatusOK, gin.H{
					"resno": 5,
					"msg":   "系统繁忙，请稍候再试！",
				})
			}

		}
	} else {
		// 获取请求Json失败
		c.JSON(http.StatusOK, gin.H{
			"resno": 1,
			"msg":   "数据错误，请稍后再试！",
		})
	}
}
