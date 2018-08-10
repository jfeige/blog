package route

import (
	//"gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"net/http"
	"sync"
	"blog/models"
	"fmt"
)

/**
Session已经存在
*/
func ExistSessionWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if uid := session.Get("uid");uid != nil{
			fmt.Println("---------")
			c.Redirect(http.StatusFound, "/manage/index")
			c.Abort()
			return
		}
	}
}

/**
Session中间件
*/
func SessionWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		c.Set("session", session)
		session.Save()
		c.Next()
	}
}

/**
没有session
*/
func NoSessionWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println("=========",session)
		if uid := session.Get("uid");uid == nil{
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
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
