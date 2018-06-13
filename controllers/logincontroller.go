package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
	"fmt"
)


func Logout(context *gin.Context){

	//销毁session
	tmpSession,exists := context.Get("session")
	if exists{
		fmt.Println("session will be deleted")
		session := tmpSession.(*models.Session)
		session.Del()
	}

	//跳转到登录页面或者前台首页
	context.Redirect(http.StatusFound,"/")
}