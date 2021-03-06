package controllers

import (
	"database/sql"
	"fmt"
	config "gin-blog/conf"
	"gin-blog/models"
	"gin-blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetArticleList 根据调用的页面不同，获取不同文章列表，并处理请求
//
func GetArticleList(c *gin.Context, username string) {
	var (
		dbLikeNumberStr string              // 数据库中文章点赞数对应字符串
		currentUser     string              // 当前登陆的用户
		likedStatus     string              // 点赞状态
		sqlString       string              // sql语句
		currentPage     int                 // 当前页码
		totalPage       int                 // 总页码
		articlesNumber  int                 // 总文章数量
		respData        []map[string]string // 最终返回的数据
		dbData          models.Article      // 文章数据模型
	)

	// 获取url中指定的page页数
	currentPage, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		fmt.Println("类型转换错误")
		currentPage = 1
	}

	// 获取当前登陆用户
	currentUser = utils.GetUserInfo(c)

	// 获取数据库文章数量
	row := utils.Db.QueryRow(utils.CountArtNumberSQL)
	err = row.Scan(&articlesNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 计算文章总页数
	totalPage = countArticleTotalPage(articlesNumber)

	// 从数据库取出当前页文章数据
	if username == "" {
		sqlString = fmt.Sprintf(utils.IndexSQL, (currentPage-1)*config.ArticlesPerPage, config.ArticlesPerPage)
	} else {
		sqlString = fmt.Sprintf(utils.UserArtListSQL, (currentPage-1)*config.ArticlesPerPage, config.ArticlesPerPage)
	}

	rows, err := getCurrentPageArticles(currentPage, sqlString, username)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		// 获取文章相关数据
		rows.Scan(&dbData.ID, &dbData.Title, &dbData.Content, &dbData.Pageview, &dbData.PubDatetime, &dbData.AuthorName, &dbData.LikeNum)

		// 获取点赞状态
		if currentUser != "" {
			likedStatus = getLikedStatus(utils.JointKey(currentUser, dbData.ID))
		} else {
			likedStatus = "unliked"
		}
		// 获取点赞数
		// likedNumberHashKey :=
		dbLikeNumberStr = getLikedNumber(utils.JointKey("article", dbData.ID), strconv.Itoa(dbData.LikeNum))
		respData = append(respData, map[string]string{
			"articleID":    strconv.Itoa(dbData.ID),
			"title":        dbData.Title,
			"artContent":   dbData.Content,
			"pageView":     strconv.Itoa(dbData.Pageview),
			"pubDateTime":  utils.B2S(dbData.PubDatetime),
			"author":       dbData.AuthorName,
			"dbLikeNumber": dbLikeNumberStr,
			"isLiked":      likedStatus,
		})
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"content":     respData,
		"page":        "首页",
		"currentUser": currentUser,
		"totalPage":   totalPage,
		"currentPage": currentPage,
		"username":    username,
	})
}

// countArticleTotalPage 计算所有文章总共可以分为多少页
func countArticleTotalPage(articlesNumber int) int {
	var totalPage int
	if (articlesNumber % config.ArticlesPerPage) == 0 {
		totalPage = (articlesNumber / config.ArticlesPerPage)
		if totalPage == 0 {
			totalPage = 1
		}
	} else {
		totalPage = (articlesNumber / config.ArticlesPerPage) + 1
	}
	return totalPage
}

// getCurrentPageArticles 获取当前页文章数据
func getCurrentPageArticles(currentPage int, sqlString string, username string) (*sql.Rows, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if username == "" {
		// 首页文章列表
		rows, err = utils.Db.Query(sqlString)
	} else {
		// 用户文章列表
		rows, err = utils.Db.Query(sqlString, username)
	}
	// rows, err := utils.Db.Query(sqlString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return rows, nil
}

// getLikedStatus 获取文章点赞状态
func getLikedStatus(redisKey string) string {
	redisConn := utils.RedisPool.Get()
	defer redisConn.Close()
	isLiked, err := redisConn.Do("hget", "likedData", redisKey)
	if err != nil {
		fmt.Println(err)
		return "unliked"
	}
	if (isLiked == nil) || (utils.RedisGetStringResult(isLiked) == "0") {
		return "unliked"
	}
	return "liked"
}

// getLikedNumber 获取redis中文章点赞数，
// 如果没有，将mysql数据库中的文章点赞数设置到redis中，并返回
func getLikedNumber(hashKey, dbLikeNum string) string {
	redisConn := utils.RedisPool.Get()
	defer redisConn.Close()
	likeNumber, err := redisConn.Do("hget", "articleLikeNumber", hashKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if likeNumber == nil {
		// redis中没有文章点赞数，从数据库中取，并存入redis
		_, err = redisConn.Do("hmset", "articleLikeNumber", hashKey, dbLikeNum)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return dbLikeNum
	}
	return utils.RedisGetStringResult(likeNumber)
}
