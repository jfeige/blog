package controllers

import (
	"blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

/**
错误处理
*/
func error(context *gin.Context, gh map[string]interface{}) {

	requestUri := strings.ToUpper(context.Request.RequestURI)

	if strings.HasPrefix(requestUri, "/MANAGE") {

		context.HTML(http.StatusOK, "manage/error.html", gh)

	} else {
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

		articleList = models.FilterNilArticle(articleList)
		columnList = models.FilterNilColumn(columnList)

		wg.Wait()

		gh["articleList"] = articleList
		gh["columnList"] = columnList
		gh["webSet"] = webSet

		context.HTML(http.StatusOK, "front/error.html", gh)
	}
}

/**
没有找到路由
*/

func NoRouter(context *gin.Context) {
	gh := make(map[string]interface{})
	gh["errcode"] = "404"
	gh["errinfo"] = "页面找不到了!"

	error(context, gh)
}

/**
缺少参数
*/
func ErrArgs(context *gin.Context) {

	gh := make(map[string]interface{})
	gh["errcode"] = "400"
	gh["errinfo"] = "参数错误，再试一次吧!"

	error(context, gh)
}

/**
统一错误处理
*/
func ToError(context *gin.Context, errcode int, errinfo string) {
	gh := make(map[string]interface{})
	gh["errcode"] = errcode
	gh["errinfo"] = errinfo

	error(context, gh)
}
