package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
	"fmt"
)

func Logout(context *gin.Context) {

	//销毁session
	tmpSession, exists := context.Get("session")
	if exists {
		fmt.Println("-----------------")
		session := tmpSession.(sessions.Session)
		session.Clear()
	}

	//跳转到登录页面或者前台首页
	context.Redirect(http.StatusFound, "/")
}
