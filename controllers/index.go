package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"yzBlog/models.go"

	"yzBlog/utils"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	sqlStr := "select id,username,password,email from users where username = ?"
	row := utils.Db.QueryRow(sqlStr, "user1")
	user := &models.User{}
	row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	fmt.Println("------------------------")
	fmt.Println(user.ID, user.Username, user.Password, user.Email)
	var content [5]map[string]string
	for i := 0; i < 5; i++ {
		content[i] = map[string]string{
			"title": "article_" + strconv.Itoa(i),
			"user":  "user_" + strconv.Itoa(i),
		}
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"content": content,
	})
}
