package controllers

import (
	"database/sql"
	"net/http"

	"gin-blog/models"
	"gin-blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func LoginGet(c *gin.Context) {

	csrfToken := csrf.GetToken(c)
	c.HTML(http.StatusOK, "account/login.html", gin.H{
		"page":      "登陆",
		"csrfToken": csrfToken,
	})
}

func LoginPost(c *gin.Context) {
	// 获取请求中的JSON数据
	var req_date models.Account
	if err := c.ShouldBindJSON(&req_date); err == nil {

		// 查询用户数据库，校验用户是否存在
		loginSQL := "select username,password from account where username = ?"
		row := utils.Db.QueryRow(loginSQL, req_date.Username)
		var db_username, db_password string
		err = row.Scan(&db_username, &db_password)
		if err == nil {
			// session := sessions.Default(c)
			// v := session.Get("mysession")
			// fmt.Println("--------------", v)
			// session.Set("username", req_date.Username)
			// session.Save()

			// 验证用户名密码
			if (req_date.Username == db_username) && (req_date.Password == db_password) {
				c.JSON(http.StatusOK, gin.H{
					"resno":    0,
					"username": req_date.Username,
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
