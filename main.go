package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/controllers"
	"blog/models"
	//"html/template"
	//"fmt"
	"fmt"
)


func main() {

	err := models.InitBaseConfig("./conf/blog.ini")
	fmt.Println(err)
	router := gin.Default()

	//t, err := template.ParseFiles("views/index.tmpl")

	//fmt.Println(t,err)
	router.LoadHTMLGlob("views/*")
	router.GET("/index", controllers.Index)

	http.ListenAndServe(":8080", router)
}

func lifei(context *gin.Context) {



	context.String(http.StatusOK, "hahahah")
}
