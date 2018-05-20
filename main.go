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
	fmt.Println(err)
	router := gin.Default()

	//t, err := template.ParseFiles("views/index.tmpl")

	router.LoadHTMLGlob("views/*")
	//router.Static("/static/","/Users/lifei/Documents/golang/src/blog/static/")
	router.Static("/static", "./static")
	router.GET("/index", controllers.Index)
	router.GET("/article",controllers.Article)

	http.ListenAndServe(":8080", router)
}

