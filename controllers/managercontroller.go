package controllers

import (
	"blog/models"
	log "github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

/**
后台首页
*/
func MIndex(context *gin.Context) {

	context.HTML(http.StatusOK, "manage/index.html", nil)
}

/**
网站设置
*/
func Webset(context *gin.Context) {

	//网站设置&&个人档案
	webSet := new(models.Webset)
	webSet.Load()

	context.HTML(http.StatusOK, "manage/webset.html", gin.H{
		"webSet": webSet,
	})
}

/**
保存网站设置
*/
func UpdateWebSet(context *gin.Context) {
	sid, _ := context.GetPostForm("id")
	id, _ := strconv.Atoi(sid)
	sitename, _ := context.GetPostForm("sitename")
	sitedesc, _ := context.GetPostForm("sitedesc")
	siteurl, _ := context.GetPostForm("siteurl")
	keywords, _ := context.GetPostForm("keywords")
	descri, _ := context.GetPostForm("descri")
	name, _ := context.GetPostForm("name")
	phone, _ := context.GetPostForm("phone")
	qq, _ := context.GetPostForm("qq")
	email, _ := context.GetPostForm("email")
	place, _ := context.GetPostForm("place")
	github, _ := context.GetPostForm("github")

	webSet := new(models.Webset)
	webSet.Id = id
	webSet.Sitename = sitename
	webSet.Sitedesc = sitedesc
	webSet.Siteurl = siteurl
	webSet.Keywords = keywords
	webSet.Descri = descri
	webSet.Name = name
	webSet.Phone = phone
	webSet.Qq = qq
	webSet.Email = email
	webSet.Place = place
	webSet.Github = github

	err := webSet.UpdateWebSet()
	var errcode int
	var errinfo string
	if err != nil {
		log.Error("updateWetSet has error:%v", err)
		errcode = -1
	}
	context.JSON(http.StatusOK, gin.H{
		"errcode": errcode,
		"errinfo": errinfo,
	})
}

/**
标签管理
*/
func Tag(context *gin.Context) {
	var wg sync.WaitGroup
	//标签
	tag_list := models.TagList()
	tagList := make([]*models.Tag, len(tag_list))
	for pos, id := range tag_list {
		wg.Add(1)
		go models.MultipleLoadTag(id, pos, tagList, &wg)
	}

	wg.Wait()
	tagList = models.FilterNilTag(tagList)

	context.HTML(http.StatusOK, "manage/tag.html", gin.H{
		"tagList": tagList,
	})
}
