package controllers

import (
	"database/sql"
	"fmt"
	config "gin-blog/conf"
	"gin-blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IndexGet 处理首页Get请求
func IndexGet(c *gin.Context) {
	GetArticleList(c, "")
}

// GetArticleList 根据调用的页面不同，获取不同文章列表，并处理请求
func GetArticleList(c *gin.Context, username string) {
	var (
		dbArtIDStr      string              // 数据库中文章id对应字符串
		dbTitle         string              // 数据库中文章标题
		dbContent       string              // 数据库中文章内容
		dbAuthor        string              // 数据库中文章作者
		dbPageView      int                 // 数据库中文章访问量
		dbPageViewStr   string              // 数据库中文章访问量对应字符串
		dbLikeNumberStr string              // 数据库中文章点赞数对应字符串
		currentUser     string              // 当前登陆的用户
		hashKeyPrefix   string              // redis哈希类型键的前缀
		likedStatus     string              // 点赞状态
		sqlString       string              // sql语句
		currentPage     int                 // 当前页码
		totalPage       int                 // 总页码
		articlesNumber  int                 // 总文章数量
		dbArticleID     int                 // 数据库中文章的id，整型
		dbLikeNumber    int                 // 数据库中文章的点赞数
		dbUpdDatetime   []uint8             // 数据库中文章更新时间
		respData        []map[string]string // 最终返回的数据
	)
	currentPage, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		fmt.Println("类型转换错误")
		currentPage = 1
	}

	// 获取当前登陆用户
	currentUser = utils.GetUserInfo(c)

	// 获取数据库文章数量
	articlesNumberSQL := "select count(title) from article"
	row := utils.Db.QueryRow(articlesNumberSQL)
	err = row.Scan(&articlesNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 计算文章总页数
	totalPage = countArticleTotalPage(articlesNumber)

	// 从数据库去除当前页文章数据
	if username == "" {
		// 首页文章查询
		sqlString = fmt.Sprintf("select title,content,pageview,upd_datetime,author_name, like_num from article order by upd_datetime desc limit %d,%d;", (currentPage-1)*config.ArticlesPerPage, config.ArticlesPerPage)
	} else {
		// 当前登陆用户文章列表查询
		sqlString = fmt.Sprintf("select title,content,pageview,upd_datetime,author_name, like_num from article where author_name=? order by upd_datetime desc limit %d,%d;", (currentPage-1)*config.ArticlesPerPage, config.ArticlesPerPage)
	}
	rows, err := getCurrentPageArticles(currentPage, sqlString, username)
	if err != nil {
		fmt.Println(err)
		return
	}

	// hash键值前缀
	if currentUser != "" {
		hashKeyPrefix = fmt.Sprintf("%s::", currentUser)
	}

	for rows.Next() {
		// 获取文章相关数据
		rows.Scan(&dbTitle, &dbContent, &dbPageView, &dbUpdDatetime, &dbAuthor, &dbLikeNumber)
		dbArtIDStr = strconv.Itoa(dbArticleID)
		dbPageViewStr = strconv.Itoa(dbPageView)
		// 获取点赞状态
		if hashKeyPrefix != "" {
			likedStatus = getLikedStatus(hashKeyPrefix, dbArtIDStr)
		} else {
			likedStatus = "unliked"
		}
		// 获取点赞数
		likedNumberHashKey := "article::" + dbArtIDStr
		dbLikeNumberStr = getLikedNumber(likedNumberHashKey, strconv.Itoa(dbLikeNumber))
		respData = append(respData, map[string]string{
			"articleID":    dbArtIDStr,
			"title":        dbTitle,
			"artContent":   dbContent,
			"pageView":     dbPageViewStr,
			"updDateTime":  utils.B2S(dbUpdDatetime),
			"author":       dbAuthor,
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
		rows, err = utils.Db.Query(sqlString)
	} else {
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
func getLikedStatus(hashKeyPrefix, dbArtIDStr string) string {
	redisConn := utils.RedisPool.Get()
	defer redisConn.Close()
	likedDataHashKey := hashKeyPrefix + dbArtIDStr
	isLiked, err := redisConn.Do("hget", "likedData", likedDataHashKey)
	if err != nil {
		fmt.Println(err)
		return "unliked"
	}
	if (isLiked == nil) || (utils.RedisGetStringResult(isLiked) == "0") {
		return "unliked"
	}
	return "liked"
}

// getLikedNumber 获取redis中文章点赞数，如果没有，将mysql数据库中的文章点赞数设置到redis中，并返回
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
