package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/controllers"
	"blog/models"
	"fmt"
	"sync"
)


func main() {
	err := models.InitBaseConfig("./conf/blog.ini")
	if err != nil{
		fmt.Println(err)
		return
	}
	router := gin.Default()

	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")
	//首页
	router.GET("/",FrontWare(),controllers.Index)
	router.GET("/index/:page",FrontWare(),controllers.Index)
	//文章详情页面
	router.GET("/article/:arteid",FrontWare(), controllers.Article)
	//类别页面
	router.GET("/category/:cateid/:page",FrontWare(), controllers.CategoryIndex)
	//标签页面
	router.GET("/tag/:tagId",FrontWare(), controllers.Index)
	//留言板
	http.ListenAndServe(":8080", router)
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
		args["isshow"] = -1
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
