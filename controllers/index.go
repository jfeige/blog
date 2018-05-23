package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
	"sync"
	"strconv"
	"math"
)


//首页
func Index(context *gin.Context){

	var wg sync.WaitGroup
	//文章列表
	tmpPage := context.Param("page")
	page,err := strconv.Atoi(tmpPage)
	if err != nil || page < 1{
		page = 1
	}

	allCnt := models.ArticleCnt()			//文章总数量
	pagesize := models.BlogPageSize
	allPage := math.Ceil(float64(allCnt)/float64(pagesize))
	if float64(page) > allPage{
		page = 1
	}

	offset := (page - 1) * pagesize

	args := make(map[string]int)
	args["isshow"] = -1						//博客的显示控制 -1:全部;1:显示;0:隐藏
	args["pagesize"] = pagesize
	args["offset"] = offset

	article_ids := models.ArticleList(args)
	articleList := make([]*models.Article,len(article_ids))
	for pos,id := range article_ids{
		wg.Add(1)
		models.MultipleLoadArticle(id,pos,articleList,&wg)
	}
	wg.Wait()

	pages := make([]int,0)
	for i := 1; i <= int(allPage);i++{
		pages = append(pages,i)
	}

	//读取中间件传来的参数
	tmp_gh,_ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["articleList"] = articleList
	gh["allPage"] = allPage
	gh["pages"] = pages
	gh["page"] = page
	gh["url"] = "/index"

	context.HTML(http.StatusOK,"index.html",gh)
}

//文章页面
func Article(context *gin.Context){
	id := context.Param("arteid")
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