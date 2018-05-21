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
	wg.Wait()
	//近期文章
	recentList := articleList
	if len(articleList) >= 6{
		recentList = articleList[:6]
	}
	//友情链接
	flink_list := models.FLink_List()
	flinkList := make([]*models.FriendLink,len(flink_list))
	for pos,id := range flink_list{
		wg.Add(1)
		models.MultipleLoadFLink(id,pos,flinkList,&wg)
	}
	wg.Wait()
	context.HTML(http.StatusOK,"index.html",
		gin.H{
		"articleList":articleList,
		"webSet":webSet,
		"categoryList":categoryList,
		"recentList":recentList,
		"flinkList":flinkList,
		})
}

//文章页面
func Article(context *gin.Context){
	id := context.Param("artiId")
	artiId,err := strconv.Atoi(id)
	if err != nil || artiId <= 0{

	}
}