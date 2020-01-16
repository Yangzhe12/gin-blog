package controllers

import (
	"database/sql"
	"gin-blog/models"
	"gin-blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func RegistGet(c *gin.Context) {
	// 生成csrfToken
	csrfToken := csrf.GetToken(c)

	c.HTML(http.StatusOK, "account/regist.html", gin.H{
		"page":      "注册",
		"csrfToken": csrfToken,
	})
}

func RegistPost(c *gin.Context) {
	// 获取请求中的JSON数据
	var req_date models.Account
	if err := c.ShouldBindJSON(&req_date); err == nil {

		// 查询用户数据库，校验用户是否存在
		registSQL := "select username, email from account where username = ? or email = ?"
		row := utils.Db.QueryRow(registSQL, req_date.Username, req_date.Email)
		var db_username, db_email string
		err = row.Scan(&db_username, &db_email)
		if err == nil {
			if db_username == req_date.Username {
				c.JSON(http.StatusOK, gin.H{
					"resno": 3,
					"msg":   "用户名已存在！",
				})
			} else if db_email == req_date.Email {
				c.JSON(http.StatusOK, gin.H{
					"resno": 4,
					"msg":   "邮箱已存在！",
				})
			}
		} else {
			// 用户不存在，可以注册，插入注册信息
			if err == sql.ErrNoRows {
				addUserSQL := "insert into account (username, password, email, phone) values (?,?,?,?)"
				_, err := utils.Db.Exec(addUserSQL, req_date.Username, req_date.Password, req_date.Email, req_date.Phone)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"resno": 2,
						"msg":   "注册失败，请重新注册！",
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"resno": 0,
						"msg":   "注册成功！",
					})
				}
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
			"msg":   "注册失败：数据错误！",
		})
	}

}
