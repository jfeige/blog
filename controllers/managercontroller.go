package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/models"
	"sync"
)

/**
	后台首页
 */
func MIndex(context *gin.Context){

	context.HTML(http.StatusOK,"manage/index.html",nil)
}

/**
	网站设置
 */
func Webset(context *gin.Context){

	//网站设置&&个人档案
	webSet := new(models.Webset)
	webSet.Load()

	context.HTML(http.StatusOK,"manage/webset.html",gin.H{
		"webSet":webSet,
	})
}

/**
	标签管理
 */
func Tag(context *gin.Context){
	if context.Request.Method == "GET"{
		var wg sync.WaitGroup
		//标签
		tag_list := models.TagList()
		tagList := make([]*models.Tag,len(tag_list))
		for pos,id := range tag_list{
			wg.Add(1)
			models.MultipleLoadTag(id,pos,tagList,&wg)
		}

		wg.Wait()

		context.HTML(http.StatusOK,"manage/tag.html",gin.H{
			"tagList":tagList,
		})
	}else{

	}

}