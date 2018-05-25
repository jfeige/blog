package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/models"
)

func MIndex(context *gin.Context){

	context.HTML(http.StatusOK,"manage/index.html",nil)
}

func Webset(context *gin.Context){

	//网站设置&&个人档案
	webSet := new(models.Webset)
	webSet.Load()

	context.HTML(http.StatusOK,"manage/webset.html",gin.H{
		"webSet":webSet,
	})
}
