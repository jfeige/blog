package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
)


func Index(context *gin.Context){

	//文章列表
	article_list := make([]*models.Article,0)
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
		article_list = append(article_list,article)
	}
	//网站设置&&个人档案
	webSet := new(models.Webset)
	webSet.Load()

	//分类列表
	categoryList := make([]*models.Category,0)
	categroy_list := models.CategoryList()
	category := new(models.Category)
	for _,id := range categroy_list{
		err := category.Load(id)
		if err != nil{
			continue
		}
		categoryList = append(categoryList,category)
	}
	//近期文章
	recent_list := article_list
	if len(article_list) >= 6{
		recent_list = article_list[:6]
	}
	//友情链接
	flinkList := make([]*models.FriendLink,0)
	flink_list := models.FLink_List()
	flink := new(models.FriendLink)
	for _,id := range flink_list{
		err := flink.Load(id)
		if err != nil{
			continue
		}
		flinkList = append(flinkList,flink)
	}


	context.HTML(http.StatusOK,"index.tmpl",
		gin.H{
		"articleList":article_list,
		"webSet":webSet,
		"categoryList":categoryList,
		"recent_List":recent_list,
		"flink_List":flinkList,
		})
}