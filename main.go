package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/controllers"
	"blog/models"
	"fmt"
)


func main() {

	err := models.InitBaseConfig("./conf/blog.ini")
	if err != nil{
		fmt.Println(err)
		return
	}
	router := gin.Default()

	//t, err := template.ParseFiles("views/index.tmpl")

	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")
	//首页
	router.GET("/index", controllers.Index)
	//文章页面
	router.GET("/article/:artiId", controllers.Index)
	//类别页面
	router.GET("/category/:cateId", controllers.Index)
	//留言板

	http.ListenAndServe(":8080", router)
}

