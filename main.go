package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/controllers"
	"net/http"
	"blog/models"
	"fmt"
	"sync"
)


func main() {

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode

	err := models.InitBaseConfig("./conf/blog.ini")
	if err != nil{
		fmt.Println(err)
		return
	}

	//前台文章浏览量入库
	go controllers.ProcessReadData()

	router := initRouter()

	http.ListenAndServe(":8080", router)

}




func initRouter()*gin.Engine{

	router := gin.Default()

	//模版文件和静态资源文件
	router.LoadHTMLGlob("views/**/*")
	//静态文件加载
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico","./static/img/favicon.ico")


	//前台页面-----------------------------------------------------------
	//首页
	router.GET("/",FrontWare(),SessionWare(),controllers.Index)
	router.GET("/index/*page",SessionWare(),FrontWare(),controllers.Index)
	//文章详情页面
	router.GET("/article/:arteid",FrontWare(), controllers.Article)
	//类别页面
	router.GET("/category/:cateid/*page",FrontWare(), controllers.CategoryFront)
	//标签页面
	router.GET("/tag/*tagid",FrontWare(), controllers.TagIndex)
	//添加一条回复
	router.POST("/comment/addComment",SessionWare(),controllers.AddComment)
	//留言板

	//跳转到登录页面
	router.GET("/login",controllers.Login)
	//登录
	router.POST("/mlogin",SessionWare(),controllers.MLogin)

	//后台页面-----------------------------------------------------------
	manageRouter := router.Group("/manage")
	manageRouter.Use(NoSessionWare())
	//后台首页
	manageRouter.GET("",controllers.MIndex)
	manageRouter.GET("/index",controllers.MIndex)
	manageRouter.GET("/webset",controllers.Webset)
	//更新网站设置
	manageRouter.POST("/updateWebSet",controllers.UpdateWebSet)
	//标签首页
	manageRouter.GET("/tag",controllers.Tag)
	//添加标签
	manageRouter.POST("/addTag",controllers.AddTag )
	//删除标签
	manageRouter.POST("/delTag",controllers.DelTag)
	//类别首页
	manageRouter.GET("/category",controllers.CategoryManage)
	//添加类别
	manageRouter.POST("/addcatetory",controllers.AddCategory)
	//删除类别
	manageRouter.POST("/delcategory",controllers.DelCategory)
	//修改类别，跳转到修改页面
	manageRouter.GET("/updateCatetory/:cateid",controllers.UpdateCatetory)
	//提交修改类别
	manageRouter.POST("/upCatetory",controllers.UpCatetory)
	//查看文章列表
	manageRouter.GET("/articleList/*cateid",controllers.ArticleList)
	//评论列表
	manageRouter.GET("/commentList/*arteid",controllers.CommentList)
	//文章详情
	manageRouter.GET("/articleinfo/:arteid",controllers.ArticleInfo)
	//文章修改
	manageRouter.POST("/updateArticle",controllers.UpfateArticleInfo)
	//添加文章
	manageRouter.Any("/addArticle",controllers.AddArticle)
	//删除文章
	manageRouter.GET("/delArticle",controllers.DelArticle)
	//查看评论详情
	manageRouter.GET("/viewComment/:cid",controllers.CommentInfo)
	//删除评论
	manageRouter.POST("/delComment",controllers.DelComment)

	//退出登录
	manageRouter.GET("/logout",controllers.Logout)


	manageRouter.GET("/test",func(context *gin.Context){
		context.HTML(http.StatusOK,"test.html",nil)
	})

	//404处理
	router.NoRoute(controllers.ToError)


	return router
}




/**
	Session中间件
 */
func SessionWare()gin.HandlerFunc{
	return func(c *gin.Context){
		var session *models.Session
		sessid,_ := c.Cookie("session_id")
		if sessid == ""{
			session = models.NewSession()
		}else{
			session = models.NewSession(sessid)
		}
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    session.SessionID(),
			Path:     "/",
			HttpOnly: true,
		}
		if session.Has("uid"){
			session.Expire()
		}
		http.SetCookie(c.Writer, cookie)

		c.Set("session", session)
		c.Next()
	}
}

/**
	没有session
 */
func NoSessionWare()gin.HandlerFunc{
	return func(c *gin.Context){
		var session *models.Session
		sessid,_ := c.Cookie("session_id")
		if sessid == ""{
			session = models.NewSession()
		}else{
			session = models.NewSession(sessid)
		}
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    session.SessionID(),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		if !session.Has("uid"){
			c.Redirect(http.StatusFound,"/login")
			c.Abort()
			return
		}
		session.Expire()
		c.Set("session", session)
		c.Next()
	}
}


/**
	前台页面专用中间件，用于读取页面右侧数据
 */
func FrontWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wg sync.WaitGroup

		//网站设置&&个人档案
		webSet := new(models.Webset)
		webSet.Load()

		//分类列表
		categroy_list := models.CategoryList()
		categoryList := make([]*models.Category,len(categroy_list))
		for pos,id := range categroy_list{
			wg.Add(1)
			models.MultipleLoadCategory(id,pos,categoryList,&wg)
		}

		//友情链接
		flink_list := models.FLink_List()
		flinkList := make([]*models.FriendLink,len(flink_list))
		for pos,id := range flink_list{
			wg.Add(1)
			models.MultipleLoadFLink(id,pos,flinkList,&wg)
		}
		//标签
		tag_list := models.TagList()
		tagList := make([]*models.Tag,len(tag_list))
		for pos,id := range tag_list{
			wg.Add(1)
			models.MultipleLoadTag(id,pos,tagList,&wg)
		}

		//文章列表
		args := make(map[string]int)
		args["page"] = 1
		args["pagesize"] = 10
		args["offset"] = 0
		article_ids := models.ArticleList(args)
		articleList := make([]*models.Article,len(article_ids))
		for pos,id := range article_ids{
			wg.Add(1)
			models.MultipleLoadArticle(id,pos,articleList,&wg)
		}
		wg.Wait()

		//近期文章
		recentList := articleList
		if len(articleList) >= 6{
			recentList = articleList[:6]
		}

		wg.Wait()

		gh := make(map[string]interface{})
		gh["webSet"] = webSet
		gh["categoryList"] = categoryList
		gh["flinkList"] = flinkList
		gh["tagList"] = tagList
		gh["recentList"] = recentList

		c.Set("gh", gh)
		c.Next()
	}

}
