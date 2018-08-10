package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
)

func Logout(c *gin.Context) {

	//销毁session
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	//跳转到登录页面或者前台首页
	c.Redirect(http.StatusFound, "/")
}
