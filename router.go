package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/controllers"
)

func initRouter()*gin.Engine{
	router := gin.Default()


	//模版文件和静态资源文件
	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")


	//首页
	router.GET("/",FrontWare(),SessionWare(),controllers.Index)
	router.GET("/index/*page",SessionWare(),FrontWare(),controllers.Index)
	//文章详情页面
	router.GET("/article/:arteid",FrontWare(), controllers.Article)
	//类别页面
	router.GET("/category/:cateid/*page",FrontWare(), controllers.CategoryIndex)
	//标签页面
	router.GET("/tag/*tagid",FrontWare(), controllers.TagIndex)
	//添加一条回复
	router.POST("/comment/add",controllers.AddComment)
	//留言板

	//跳转到登录页面
	router.GET("/login",controllers.Login)

	router.POST("/login",SessionWare(),controllers.Login)
	//404处理
	router.NoRoute(NoRouteWare(),controllers.ErrNoRoute)


	return router
}
