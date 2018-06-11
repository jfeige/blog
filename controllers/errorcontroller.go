package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"sync"
	"blog/models"
	"strings"
	"fmt"
)

/**
	没有找到路由  /manage/index2
 */
func ToError(context *gin.Context,gh map[string]interface{}){

	requestUri := strings.ToUpper(context.Request.RequestURI)
	if strings.HasPrefix(requestUri,"/MANAGE"){
		context.HTML(http.StatusOK,"manage/error.html",gh)
	}else{
		var wg sync.WaitGroup

		webSet := new(models.Webset)
		webSet.Load()

		//推荐阅读
		args := make(map[string]int)
		args["page"] = 1
		args["isshow"] = -1
		args["pagesize"] = 10
		args["offset"] = 0
		article_ids := models.ArticleList(args)
		articleList := make([]*models.Article,len(article_ids))
		for pos,id := range article_ids{
			wg.Add(1)
			go models.MultipleLoadArticle(id,pos,articleList,&wg)
		}

		//首页菜单
		column_ids := models.ColumnList()
		columnList := make([]*models.Column,len(column_ids))
		for pos,id := range column_ids{
			wg.Add(1)
			go models.MultipleLoadColumn(id,pos,columnList,&wg)
		}

		articleList = models.FilterNilArticle(articleList)
		columnList = models.FilterNilColumn(columnList)

		wg.Wait()

		gh["articleList"] = articleList
		gh["columnList"] = columnList
		gh["webSet"] = webSet

		context.HTML(http.StatusOK,"front/error.html",gh)
	}
}


func NoRouter(context *gin.Context){
	fmt.Println("---------NoRouter")
	gh := make(map[string]interface{})
	gh["errcode"] = "404"
	gh["errinfo"] = "页面找不到了!"

	ToError(context,gh)
	//context.HTML(http.StatusOK,"error.html",gh)
}

/**
	缺少参数
 */
func ErrLackArgs(context *gin.Context){

	gh := make(map[string]interface{})
	gh["errcode"] = "400"
	gh["errinfo"] = "参数不全，请重试!"


	context.HTML(http.StatusOK,"error.html",gh)

}