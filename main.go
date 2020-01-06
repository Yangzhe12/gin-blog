package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yzBlog/controllers"
	"yzBlog/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("err open databases: ", err)
		return
	}
	defer db.Close()
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(filepath.Join(getCurrentDirectory(), "./views/**/*")))
	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	router.GET("/", controllers.IndexGet)

	router.Run(":8080")
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
