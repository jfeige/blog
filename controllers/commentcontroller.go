package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
	"strconv"
	"fmt"
)


/**
	添加一条回复	ajax请求
 */
func AddComment(context *gin.Context){

	var name interface{}
	tmpSession,exists := context.Get("session")
	if exists{
		session := tmpSession.(*models.Session)
		if session.Has("uid"){
			name = session.GetSession("name")
		}else{
			name = "guest"
		}
	}
	fmt.Println("---name---",name)
	//articleid
	articleid,ok := context.GetPostForm("aid")
	if !ok{
		//参数错误
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"参数不全，请刷新该页面重试",
		})
	}
	aid,err := strconv.Atoi(articleid)
	if err != nil{
		//参数错误
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"参数错误，请刷新该页面重试",
		})
	}
	fmt.Println("---aid---",aid)
	content,ok := context.GetPostForm("content")
	if !ok{
		//参数错误
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"参数不全，请刷新该页面重试",
		})
	}
	if content == ""{
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"内容不能为空",
		})
	}
	fmt.Println("---content---",content)
	models.AddComment(aid,name,content)

	context.JSON(http.StatusOK,gin.H{
		"errcode":0,
		"errinfo":"",
	})
}