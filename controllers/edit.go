package controllers

import (
	config "gin-blog/conf"
	"gin-blog/models"
	"gin-blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditGet(c *gin.Context) {
	c.HTML(http.StatusOK, "editor/editor.html", gin.H{
		"page": "写文章",
	})
}

func EditPost(c *gin.Context) {
	// 获取请求中的JSON数据
	var requestData models.Article
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// 获取请求Json失败
		c.JSON(http.StatusBadRequest, gin.H{
			"resno": 1,
			"msg":   "数据错误，请稍后再试！",
		})
	} else {
		// 校验数据
		if requestData.Title == "" {
			c.JSON(http.StatusNoContent, gin.H{
				"resno": 2,
				"msg":   "文章标题不能为空！",
			})
			return
		} else if requestData.Content == "" {
			c.JSON(http.StatusNoContent, gin.H{
				"resno": 3,
				"msg":   "文章内容不能为空！",
			})
			return
		}
		// 获取当前用户ID
		utils.SessionKey = config.GetConfiguration().UserInfoSessionKey
		session := utils.Default(c)
		value := session.Get("userID")
		if value == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"resno": 4,
				"msg":   "登陆失效，请重新登陆！",
			})
			return
		}
		userID := value.(int)
		// 将文章保存到数据库
		AddArtcleSQL := "insert into article (title, content, author_id) values (?,?,?)"
		_, err := utils.Db.Exec(AddArtcleSQL, requestData.Title, requestData.Content, userID)
		if err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{
				"resno": 5,
				"msg":   "服务器故障，文章保存失败，请稍后再试！",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"resno": 0,
			"msg":   "保存成功",
		})
	}

}
