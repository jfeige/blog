package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
	"sync"
	"strconv"
)


//首页
func Index(context *gin.Context){

	var wg sync.WaitGroup
	//文章列表
	args := make(map[string]interface{})
	args["page"] = 1
	args["isshow"] = -1
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
	//读取中间件传来的参数
	tmp_gh,_ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["articleList"] = articleList
	gh["recentList"] = recentList

	context.HTML(http.StatusOK,"index.html",gh)
}

//文章页面
func Article(context *gin.Context){
	id := context.Param("artiId")
	artiId,err := strconv.Atoi(id)
	if err != nil || artiId <= 0{
		//参数错误，跳转到首页
		context.Redirect(0,"/")
	}

	//根据文章id读取
	article := new(models.Article)
	err = article.Load(artiId)
	if err != nil{
		//数据错误或者id不正确

	}

	var wg sync.WaitGroup

	//评论页面，不分页
	comment_list := models.CommentList(artiId)
	commentList := make([]*models.Comment,len(comment_list))
	for pos,id := range comment_list{
		wg.Add(1)
		models.MultipleLoadComment(id,pos,commentList,&wg)
	}
	wg.Wait()

	tmp_gh,_ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["commentList"] = commentList
	gh["article"] = article

	context.HTML(http.StatusOK,"article.html",gh)
}