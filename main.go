package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gin-blog/conf"
	"gin-blog/controllers"
	"gin-blog/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func main() {
	// 加载配置文件
	configFilePath := flag.String("C", "conf/conf.yaml", "config file path")
	flag.Parse()
	if err := config.LoadConfiguration(*configFilePath); err != nil {
		fmt.Println("err parsing config log file", err)
		return
	}

	configCon := config.GetConfiguration()

	// 数据库连接初始化
	db, err := utils.InitDB()
	fmt.Println(db)
	if err != nil {
		fmt.Println("err open databases: ", err)
		return
	}
	defer db.Close()

	router := gin.Default()

	// 生成csrf所需Cookie
	store := cookie.NewStore([]byte(configCon.CookieSecret))
	var csrfOption sessions.Options
	csrfOption.MaxAge = configCon.CsrfTokenValidTime
	store.Options(csrfOption)

	// 用户信息session
	// redisStore, err := redis.NewStore(10, "tcp", configCon.RedisAddress, "", []byte(configCon.UserInfoSessionKey))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var redisOption sessions.Options
	// redisOption.MaxAge = configCon.RedisSessionValidTime
	// redisStore.Options(redisOption)

	// Session中间件
	router.Use(sessions.Sessions(configCon.CookieName, store))
	// router.Use(sessions.Sessions("se-ssion", redisStore))

	// csrf中间件
	router.Use(csrf.Middleware(csrf.Options{
		Secret: configCon.CsrfTokenSecret,
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"resno": 10,
				"msg":   "长时间未响应，请刷新页面后重试！",
			})
			c.Abort()
		},
	}))

	// 静态文件路径
	router.LoadHTMLGlob(filepath.Join(filepath.Join(getCurrentDirectory(), "./views/**/*")))
	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))

	// 请求Handler
	router.GET("/", controllers.IndexGet)

	router.GET("/regist", controllers.RegistGet)
	router.POST("/regist", controllers.RegistPost)

	router.GET("/login", controllers.LoginGet)
	router.POST("/login", controllers.LoginPost)

	router.GET("/logout", controllers.LogoutGet)

	router.GET("/editor", controllers.EditorGet)

	router.Run(":8080")
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
