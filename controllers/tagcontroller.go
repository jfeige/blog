package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"

	"fmt"
	"net/http"
	"blog/models"
)


func TagIndex(context *gin.Context){

	tmpTagid := context.Param("tagid")
	if tmpTagid == ""{
		gh := make(map[string]interface{})
		gh["errinfo"] = ""
		context.HTML(http.StatusOK,"index.html",gh)
	}else{
		fmt.Println("tagid is null",tmpTagid)
	}


}


/**
	删除一个标签
 */
func DelTag(context *gin.Context){
	gh := make(map[string]interface{})
	defer context.JSON(http.StatusOK,gh)

	tagid,ok := context.GetPostForm("id")
	if !ok{
		gh["errcode"] = -1
		gh["errinfo"] = "参数不全，请重试"
		return
	}
	errcode := models.DelTag(tagid)
	if errcode < 0{
		gh["errcode"] = -2
		gh["errinfo"] = "数据库异常，请稍后重试"
		return
	}

	gh["errcode"] = 0
}