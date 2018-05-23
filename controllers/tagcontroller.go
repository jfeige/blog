package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"

	"fmt"
	"net/http"
)
//o(╥﹏╥)o

func TagIndex(context *gin.Context){

	tmpTagid := context.Param("tagid")
	if tmpTagid == ""{
		gh := make(map[string]interface{})
		gh["errinfo"] = ""
		context.HTML(http.StatusOK,"index.html",gh)
		//context.Redirect(http.StatusMovedPermanently,"http://www.baidu.com")
	}else{
		fmt.Println("tagid is null",tmpTagid)
	}


}
