package route

import (
	"blog/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)


func aa(){
	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "182.92.158.94:6379", "lifei", []byte("secret"))
	r.Use(sessions.Sessions("session_id", store))

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8000")
}


func InitRouter() *gin.Engine {

	router := gin.Default()

	store, _ := redis.NewStore(10, "tcp", "182.92.158.94:6379", "lifei", []byte("secret"))
	router.Use(sessions.Sessions("session_id", store))
	router.Use(SessionWare())
	router.Use(FrontWare())

	//模版文件和静态资源文件
	router.LoadHTMLGlob("views/**/*")
	//静态文件加载
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	//前台页面-----------------------------------------------------------
	//首页
	router.GET("/", FrontWare(),  controllers.Index)
	router.GET("/index/*page", FrontWare(),  controllers.Index)
	//文章详情页面
	router.GET("/article/:arteid", FrontWare(), controllers.Article)
	//类别页面
	router.GET("/category/:cateid/*page", FrontWare(),  controllers.CategoryFront)
	//标签页面
	router.GET("/tag/:tagid/*page", FrontWare(),  controllers.TagIndex)
	//添加一条回复
	router.POST("/addComment",  controllers.AddComment)
	//留言板
	router.GET("/msg/*page", FrontWare(),  controllers.MessageBorad)
	//提交留言
	router.POST("/addMsg",  controllers.AddMsg)
	//关于我
	router.GET("/my", FrontWare(), controllers.Myinfo)

	//跳转到登录页面
	router.GET("/login", ExistSessionWare(), controllers.Login)
	//登录
	router.POST("/mlogin", ExistSessionWare(), controllers.MLogin)

	//后台页面-----------------------------------------------------------
	manageRouter := router.Group("/manage")
	manageRouter.Use(NoSessionWare())
	//后台首页
	manageRouter.GET("", controllers.MIndex)
	manageRouter.GET("/index", controllers.MIndex)
	//网站基本资料设置页面
	manageRouter.GET("/webset", controllers.Webset)
	//更新网站设置
	manageRouter.POST("/updateWebSet", controllers.UpdateWebSet)
	//标签首页
	manageRouter.GET("/tag", controllers.Tag)
	//添加标签
	manageRouter.POST("/addTag", controllers.AddTag)
	//删除标签
	manageRouter.POST("/delTag", controllers.DelTag)
	//类别首页
	manageRouter.GET("/category", controllers.CategoryManage)
	//添加类别
	manageRouter.POST("/addcatetory", controllers.AddCategory)
	//删除类别
	manageRouter.POST("/delcategory", controllers.DelCategory)
	//修改类别，跳转到修改页面
	manageRouter.GET("/updateCatetory/:cateid", controllers.UpdateCatetory)
	//提交修改类别
	manageRouter.POST("/upCatetory", controllers.UpCatetory)
	//查看文章列表
	manageRouter.GET("/articleList/:cateid/*page", controllers.ArticleList)
	//评论列表
	manageRouter.GET("/commentList/:arteid/*page", controllers.CommentList)
	//文章详情
	manageRouter.GET("/articleinfo/:arteid", controllers.ArticleInfo)
	//文章修改
	manageRouter.POST("/updateArticle", controllers.UpfateArticleInfo)
	//添加文章
	manageRouter.Any("/addArticle", controllers.AddArticle)
	//添加文章或者修改文章时，上传图片
	manageRouter.POST("/upimg", controllers.UpImg)
	//删除文章
	manageRouter.POST("/delArticle", controllers.DelArticle)
	//查看评论详情
	manageRouter.GET("/viewComment/:cid", controllers.CommentInfo)
	//删除评论
	manageRouter.POST("/delComment", controllers.DelComment)
	//留言管理
	manageRouter.GET("/msgList/*page", controllers.MsgList)
	//查看留言详情
	manageRouter.GET("/viewMsg/:mid", controllers.MsgInfo)
	//删除留言
	manageRouter.POST("/delMessage", controllers.DelMessage)
	//批量删除留言
	manageRouter.POST("/delMultiMessage", controllers.DelMultiMessage)
	//栏目管理
	manageRouter.GET("/column", controllers.ColumnManage)
	//添加栏目
	manageRouter.GET("/addColumn", controllers.AddColumn)
	//添加一个栏目
	manageRouter.POST("/addColumn", controllers.AddColumn)
	//修改栏目
	manageRouter.GET("/updateColumn/:cid", controllers.UpdateColumn)
	//提交栏目修改
	manageRouter.POST("/upColumn", controllers.UpColumn)
	//删除一个栏目
	manageRouter.POST("/delColumn", controllers.DelColumn)
	//友链列表
	manageRouter.GET("/flink", controllers.FlinkList)
	//添加友链
	manageRouter.POST("/addflink", controllers.AddFlink)
	//删除友链
	manageRouter.POST("/delflink", controllers.DelFlink)
	//修改友链，跳转到修改页面
	manageRouter.GET("/updateFlink/:fid", controllers.UpdateFlink)
	//提交修改友链
	manageRouter.POST("/upFlink", controllers.UpFlink)
	//退出登录
	manageRouter.GET("/logout", controllers.Logout)

	//manageRouter.GET("/test",func(context *gin.Context){
	//	context.HTML(http.StatusOK,"test.html",nil)
	//})

	//404处理(这里还有问题，前台和后台的区分问题)
	router.NoRoute(controllers.NoRouter)

	return router
}