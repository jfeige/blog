package main

import (
	"blog/controllers"
	"blog/models"
	log "github.com/alecthomas/log4go"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"runtime"
	"sync"
	"flag"
	"fmt"
)

var(

	logfile = flag.String("log","./conf/blog-log.xml","log4go file path!")
	configfile = flag.String("config","./conf/blog.ini","config file path")

)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	log.LoadConfiguration(*logfile)
	defer log.Close()

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode

	err := models.InitBaseConfig(*configfile)
	if err != nil {
		log.Error("InitBaseConfig has error:%v", err)
		return
	}

	//前台文章浏览量入库
	go controllers.ProcessReadData()

	router := initRouter()

	err = http.ListenAndServe(models.AppPort, router)
	if err != nil {
		fmt.Printf("http.ListenAndServe has error:%v\n", err)
		log.Error("http.ListenAndServe has error:%v\n", err)
		return
	}

}

func initRouter() *gin.Engine {

	router := gin.Default()

	//模版文件和静态资源文件
	router.LoadHTMLGlob("views/**/*")
	//静态文件加载
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	//前台页面-----------------------------------------------------------
	//首页
	router.GET("/", FrontWare(), SessionWare(), controllers.Index)
	router.GET("/index/*page", FrontWare(), SessionWare(), controllers.Index)
	//文章详情页面
	router.GET("/article/:arteid", FrontWare(), SessionWare(), controllers.Article)
	//类别页面
	router.GET("/category/:cateid/*page", FrontWare(), SessionWare(), controllers.CategoryFront)
	//标签页面
	router.GET("/tag/:tagid/*page", FrontWare(), SessionWare(), controllers.TagIndex)
	//添加一条回复
	router.POST("/addComment", SessionWare(), controllers.AddComment)
	//留言板
	router.GET("/msg/*page", FrontWare(), SessionWare(), controllers.MessageBorad)
	//提交留言
	router.POST("/addMsg", SessionWare(), controllers.AddMsg)
	//关于我
	router.GET("/my", FrontWare(), SessionWare(), controllers.Myinfo)

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

/**
Session已经存在
*/
func ExistSessionWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session *models.Session
		sessid, _ := c.Cookie(models.SessionName)
		if sessid != "" {
			session = models.NewSession(sessid)
			if session.Has("uid") {
				//跳转到后台首页
				c.Redirect(http.StatusFound, "/manage/index")
				c.Abort()
				return
			} else {
				cookie := &http.Cookie{
					Name:     models.SessionName,
					Value:    session.SessionID(),
					Path:     "/",
					HttpOnly: true,
				}
				http.SetCookie(c.Writer, cookie)
				c.Set("session", session)
				c.Next()
			}
		}
	}
}

/**
Session中间件
*/
func SessionWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session *models.Session
		sessid, _ := c.Cookie(models.SessionName)
		if sessid == "" {
			session = models.NewSession()
		} else {
			session = models.NewSession(sessid)
		}
		cookie := &http.Cookie{
			Name:     models.SessionName,
			Value:    session.SessionID(),
			Path:     "/",
			HttpOnly: true,
		}
		if session.Has("uid") {
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
func NoSessionWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session *models.Session
		sessid, _ := c.Cookie(models.SessionName)
		if sessid == "" {
			session = models.NewSession()
		} else {
			session = models.NewSession(sessid)
		}
		cookie := &http.Cookie{
			Name:     models.SessionName,
			Value:    session.SessionID(),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		if !session.Has("uid") {
			c.Redirect(http.StatusFound, "/login")
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
		categoryList := make([]*models.Category, len(categroy_list))
		for pos, id := range categroy_list {
			wg.Add(1)
			go models.MultipleLoadCategory(id, pos, categoryList, &wg)
		}
		//友情链接
		flink_list := models.FLink_List()
		flinkList := make([]*models.FriendLink, len(flink_list))
		for pos, id := range flink_list {
			wg.Add(1)
			go models.MultipleLoadFLink(id, pos, flinkList, &wg)
		}
		//标签
		tag_list := models.TagList()
		tagList := make([]*models.Tag, len(tag_list))
		for pos, id := range tag_list {
			wg.Add(1)
			go models.MultipleLoadTag(id, pos, tagList, &wg)
		}
		//最新评论,现实最近的6条评论(纯评论)
		args := make(map[string]int)
		args["pagesize"] = 6
		args["type"] = 0 //0:评论;1:回复;-1:全部

		comment_list := models.ManageCommentList(args)
		commentList := make([]*models.Comment, len(comment_list))
		for pos, id := range comment_list {
			wg.Add(1)
			go models.MultipleLoadComment(id, pos, commentList, &wg)
		}
		//热门文章(按照阅读量降序)
		args = make(map[string]int)
		args["order"] = 1
		args["pagesize"] = 6
		hotArticle_list := models.ArticleList(args)
		hotArticleList := make([]*models.Article, len(hotArticle_list))
		for pos, id := range hotArticle_list {
			wg.Add(1)
			go models.MultipleLoadArticle(id, pos, hotArticleList, &wg)
		}
		//文章列表
		args = make(map[string]int)
		article_ids := models.ArticleList(args)
		articleList := make([]*models.Article, len(article_ids))
		for pos, id := range article_ids {
			wg.Add(1)
			go models.MultipleLoadArticle(id, pos, articleList, &wg)
		}
		//首页菜单
		column_ids := models.ColumnList()
		columnList := make([]*models.Column, len(column_ids))
		for pos, id := range column_ids {
			wg.Add(1)
			go models.MultipleLoadColumn(id, pos, columnList, &wg)
		}

		wg.Wait()

		//过滤空数据
		articleList = models.FilterNilArticle(articleList)
		hotArticleList = models.FilterNilArticle(hotArticleList)
		categoryList = models.FilterNilCategory(categoryList)
		flinkList = models.FilterNilFriendLink(flinkList)
		tagList = models.FilterNilTag(tagList)
		columnList = models.FilterNilColumn(columnList)
		commentList = models.FilterNilComment(commentList)

		//近期文章
		recentList := articleList
		if len(articleList) >= 6 {
			recentList = articleList[:6]
		}

		gh := make(map[string]interface{})
		gh["webSet"] = webSet
		gh["categoryList"] = categoryList
		gh["flinkList"] = flinkList
		gh["tagList"] = tagList
		gh["recentList"] = recentList
		gh["commentList"] = commentList
		gh["hotArticleList"] = hotArticleList
		gh["columnList"] = columnList
		c.Set("gh", gh)
		c.Next()
	}

}
