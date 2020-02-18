package controllers

import (
	"fmt"
	config "gin-blog/conf"
	"gin-blog/models"
	"gin-blog/utils"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
)

func ArticleGet(c *gin.Context) {
	var dbTitle, dbContent, dbAuthor string
	var dbArticleID int
	var dbPubDatetime, dbUpdDatetime []uint8
	articleID, err := strconv.Atoi(c.Param("articleID"))
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusNotFound, "找不到文章！")
		return
	}
	currentUser := utils.GetUserInfo(c)
	queryArtSQL := "select id,title,content,pub_datetime,upd_datetime,author_name from article where id=?;"
	row := utils.Db.QueryRow(queryArtSQL, articleID)
	err = row.Scan(&dbArticleID, &dbTitle, &dbContent, &dbPubDatetime, &dbUpdDatetime, &dbAuthor)
	// 让字符串不转义为html格式
	if err != nil {
		fmt.Println(err)
	}
	articleCotent := template.HTML(dbContent)
	c.HTML(http.StatusOK, "article/article.html", gin.H{
		"page":          "文章内容",
		"username":      currentUser,
		"articleID":     dbArticleID,
		"title":         dbTitle,
		"articleCotent": articleCotent,
		"pubDateTime":   utils.B2S(dbPubDatetime),
		"updDateTime":   utils.B2S(dbUpdDatetime),
		"author":        dbAuthor,
	})
}

// 不转义html
func unescaped(x string) interface{} {
	return template.HTML(x)
}

func LikePost(c *gin.Context) {
	// 获取请求中的JSON数据
	var likedStatus, likedNumberStr string
	var likedNumber int
	var likedData models.LikedData
	if err := c.ShouldBindJSON(&likedData); err != nil {
		// 获取请求Json失败
		c.JSON(http.StatusBadRequest, gin.H{
			"resno": 1,
			"msg":   "数据错误，请稍后再试！",
		})
	}

	// 连接redis
	redisConn, err := redis.Dial("tcp", config.GetConfiguration().RedisAddress)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"resno": 2,
			"msg":   "服务器繁忙，请稍后再试！",
		})
		return
	}
	defer redisConn.Close()
	// 拼接redis中hash键的格式
	hashKey := fmt.Sprintf("%s::%s", likedData.CurrentUsername, likedData.LikedArticleID)
	// 获取redis中相关点赞数据
	redisLikedData, err := redisConn.Do("hget", "likedData", hashKey)
	if err != nil {
		fmt.Println(err)
	}

	// 处理查询结果
	if redisLikedData != nil {
		likedStatus = utils.RedisGetStringResult(redisLikedData)
	} else {
		likedStatus = "0"
	}

	// 获取redis中的文章点赞数量
	likedNumberHashKey := "article::" + likedData.LikedArticleID
	redisLikedNumber, err := redisConn.Do("hget", "articleLikeNumber", likedNumberHashKey)
	if err != nil {
		fmt.Println(err)
	} else {
		// redis中无点赞数
		if redisLikedNumber == nil {
			getLikedNumSQL := "select like_num from article where id=?;"
			likedRow := utils.Db.QueryRow(getLikedNumSQL, likedData.LikedArticleID)
			err = likedRow.Scan(&likedNumber)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			likedNumberStr = utils.RedisGetStringResult(redisLikedNumber)
			likedNumber, _ = strconv.Atoi(likedNumberStr)
		}
	}

	if likedStatus == "0" {
		// 当前登陆用户没有点赞过该文章
		_, err = redisConn.Do("hmset", "likedData", hashKey, "1")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"resno": 3,
				"msg":   "点赞失败",
			})
		} else {
			// 更新点赞数
			likedNumberStr = strconv.Itoa(likedNumber + 1)
			_, err = redisConn.Do("hmset", "articleLikeNumber", likedNumberHashKey, likedNumberStr)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(http.StatusOK, gin.H{
				"resno":       0,
				"likedStatus": "liked",
				"likedNumber": likedNumberStr,
				"msg":         "点赞成功",
			})
		}
	} else {
		// 当前登陆用户已经点赞过该文章，再次点击取消点赞
		_, err = redisConn.Do("hmset", "likedData", hashKey, "0")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"resno": 4,
				"msg":   "取消点赞失败",
			})
		} else {
			likedNumberStr = strconv.Itoa(likedNumber - 1)
			_, err = redisConn.Do("hmset", "articleLikeNumber", likedNumberHashKey, likedNumberStr)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(http.StatusOK, gin.H{
				"resno":       0,
				"likedStatus": "unliked",
				"likedNumber": likedNumberStr,
				"msg":         "取消点赞",
			})
		}
	}
}
