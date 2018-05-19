package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
)


func Index(context *gin.Context){

	result := make([]*models.Article,0)

	args := make(map[string]interface{})
	args["page"] = 1
	args["isshow"] = -1
	article_ids := models.ArticleList(args)

	article := new(models.Article)
	for _,id := range article_ids{
		err := article.Load(id)
		if err != nil{
			continue
		}
		result = append(result,article)
	}


	context.HTML(http.StatusOK,"index.tmpl",gin.H{
		"result":result,
	})
}