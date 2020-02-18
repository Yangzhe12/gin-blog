package main

import (
	"flag"
	"fmt"
	"gin-demo/csrf-demo/csrf"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	config "gin-blog/conf"
	"gin-blog/controllers"
	"gin-blog/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	configFilePath := flag.String("C", "conf/conf.yaml", "config file path")
	flag.Parse()
	if err := config.LoadConfiguration(*configFilePath); err != nil {
		fmt.Println("err parsing config log file", err)
		return
	}

	conf := config.GetConfiguration()

	// 数据库连接初始化
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("err open databases: ", err)
		return
	}
	defer db.Close()

	// redis连接池初始化
	redisPool := utils.InitRedisPool()
	defer redisPool.Close()

	router := gin.Default()
	setTemplate(router)
	// 使用redis存储Session
	redisStroe, err := redis.NewStore(10, "tcp", conf.RedisAddress, "", []byte(conf.CookieSecret))
	if err != nil {
		fmt.Println(err)
	}
	redisStroe.Options(sessions.Options{
		MaxAge: conf.CsrfTokenValidTime,
	})

	// CsrfToken Session中间件
	router.Use(sessions.Sessions(conf.CsrfCookieName, redisStroe))

	// 用户信息 Session中间件
	utils.SessionKey = conf.UserInfoSessionKey
	router.Use(utils.Sessions(conf.UserInfoCookieKey, redisStroe))

	// 用户数和文章数，存入redis

	// 静态文件路径
	router.LoadHTMLGlob(filepath.Join(filepath.Join(getCurrentDirectory(), "./views/**/*")))
	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))

	// 路由
	// 请求Handler
	router.GET("/", controllers.IndexGet)

	// 查看文章
	router.GET("/article/:articleID", controllers.ArticleGet)

	// 用户文章列表
	router.GET("/blog/:username", controllers.UserBlogGet)

	// 注册
	router.GET("/regist", csrfTokenFunc(), controllers.RegistGet)
	router.POST("/regist", csrfTokenFunc(), controllers.RegistPost)

	// 登陆
	router.GET("/login", csrfTokenFunc(), controllers.LoginGet)
	router.POST("/login", csrfTokenFunc(), controllers.LoginPost)

	// 退出登陆
	router.GET("/logout", controllers.LogoutGet)

	// 写文章
	router.GET("/md", controllers.MdGet)
	router.POST("/md", controllers.MdPost)

	// 点赞
	router.POST("/liked", controllers.LikePost)

	router.Run(":8080")
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// csrf中间件
func csrfTokenFunc() gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: config.GetConfiguration().CookieSecret,
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"resno": 10,
				"msg":   "长时间未响应，请刷新页面后重试！",
			})
			c.Abort()
		},
	})
}

func setTemplate(engine *gin.Engine) {

	funcMap := template.FuncMap{
		"add":   utils.Add,
		"minus": utils.Minus,
	}

	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
}
