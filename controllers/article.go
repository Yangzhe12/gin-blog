package controllers

import (
	"fmt"
	config "gin-blog/conf"
	"gin-blog/models"
	"gin-blog/utils"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ArticleGet(c *gin.Context) {
	var (
		dbTitle       string  // 从数据库中获取的文章标题
		dbContent     string  // 从数据库中获取的文章内容
		dbAuthor      string  // 从数据库中获取的文章作者
		dbPageView    int     // 从数据库中获取的文章访问量
		dbPubDatetime []uint8 // 从数据库中获取的文章发表时间
	)

	// 取出请求中文章的ID，并转为int类型
	articleID, err := strconv.Atoi(c.Param("articleID"))
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusNotFound, "找不到文章！")
		return
	}
	// 获取当前登陆的用户
	currentUser := utils.GetUserInfo(c)

	// 增加访问次数
	// 查询10分钟内是否访问过该文章
	tooOften := isPageViewInRedis(strconv.Itoa(articleID), currentUser)
	if !tooOften {
		// 没有太频繁的时候，才更新数据库中的数据
		incrPageViewSQL := "update article set pageview=pageview+1 where id=?;"
		_, err = utils.Db.Exec(incrPageViewSQL, articleID)
		if err != nil {
			fmt.Println("增加访问次数失败：", err)
		}
	}

	// 从数据库中查询指定的文章的数据
	queryArtSQL := "select title,content,pageview,pub_datetime,author_name from article where id=?;"
	row := utils.Db.QueryRow(queryArtSQL, articleID)
	err = row.Scan(&dbTitle, &dbContent, &dbPageView, &dbPubDatetime, &dbAuthor)
	if err != nil {
		fmt.Println(err)
	}

	// 让字符串不转义为html格式
	articleCotent := template.HTML(dbContent)
	c.HTML(http.StatusOK, "article/article.html", gin.H{
		"page":          dbTitle,
		"username":      currentUser,
		"articleID":     articleID,
		"title":         dbTitle,
		"articleCotent": articleCotent,
		"pageView":      strconv.Itoa(dbPageView),
		"pubDateTime":   utils.B2S(dbPubDatetime),
		"author":        dbAuthor,
	})
}

// 不转义html
func unescaped(x string) interface{} {
	return template.HTML(x)
}

// 处理点赞文章
func LikePost(c *gin.Context) {
	// 获取请求中的JSON数据，当前登陆的用户名，点赞的文章的id
	var likedData models.LikedData
	if err := c.ShouldBindJSON(&likedData); err != nil {
		// 获取请求Json失败
		c.JSON(http.StatusBadRequest, gin.H{
			"resno": 1,
			"msg":   "数据错误，请稍后再试！",
		})
		return
	}

	// 拼接redis中hash键的格式
	statusKey := fmt.Sprintf("%s::%s", likedData.CurrentUsername, likedData.LikedArticleID)
	numberKey := fmt.Sprintf("article::%s", likedData.LikedArticleID)

	// 如果当前已经点赞，则取消点赞，并设置redis中的数据
	newStatus, newNumber, err := handleRedisLikedData(likedData.LikedStatus, statusKey, numberKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resno":       0,
		"likedStatus": newStatus,
		"likedNumber": newNumber,
		"msg":         "操作成功",
	})
}

// handleRedisLikedData 处理点赞/取消点赞
func handleRedisLikedData(curStatus string, statusKey string, numberKey string) (newStatus string, newNumber string, err error) {
	redisConn := utils.RedisPool.Get()
	defer redisConn.Close()
	// 当前用户已经点赞过本文章
	if curStatus == "liked" {
		_, err = redisConn.Do("HSET", "likedData", statusKey, 0)
		if err != nil {
			fmt.Println("1----", err)
			return "", "", err
		}
		newStatus = "unliked"
		_, err = redisConn.Do("HINCRBY", "articleLikeNumber", numberKey, -1)
		if err != nil {
			fmt.Println("2----", err)
			return "", "", err
		}
	} else {
		// 当前用户未点赞过本文章
		_, err = redisConn.Do("HSET", "likedData", statusKey, 1)
		if err != nil {
			fmt.Println("3----", err)
			return "", "", err
		}
		newStatus = "liked"
		_, err = redisConn.Do("HINCRBY", "articleLikeNumber", numberKey, 1)
		if err != nil {
			fmt.Println("4----", err)
			return "", "", err
		}
	}
	newNumberObj, err := redisConn.Do("HGET", "articleLikeNumber", numberKey)
	if err != nil {
		fmt.Println("5----", err)
		return "", "", err
	}
	newNumber = utils.RedisGetStringResult(newNumberObj)
	return newStatus, newNumber, nil
}

// isPageViewInRedis 查询redis中有无当前用户访问文章的数据
// 如果规定时间重复访问，不增加访问量数据
func isPageViewInRedis(articleID string, username string) bool {
	hashKey := username + "::" + articleID
	redisConn := utils.RedisPool.Get()
	defer redisConn.Close()
	res, err := redisConn.Do("GET", hashKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// redis中不存在访问量数据，设置;
	// redis中存在访问量数据，更新过期时间
	_, _ = redisConn.Do("SETEX", hashKey, config.GetConfiguration().PageViewFlagExpiTime, "1")
	return (res != nil)
}
