package controllers

import (
	"blog/models"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Logout(context *gin.Context) {

	//销毁session
	tmpSession, exists := context.Get("session")
	if exists {
		session := tmpSession.(*models.Session)
		session.Del()
	}

	//跳转到登录页面或者前台首页
	context.Redirect(http.StatusFound, "/")
}
