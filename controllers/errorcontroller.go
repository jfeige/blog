package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

/**
	没有找到路由
 */
func ErrNoRoute(context *gin.Context){

	//读取中间件传来的参数
	tmp_gh,_ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["errcode"] = "404"
	gh["errinfo"] = "页面找不到了!"

	context.HTML(http.StatusOK,"error.html",gh)

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