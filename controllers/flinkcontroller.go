package controllers

import (
	"blog/models"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"strconv"
	"sync"
)

/**
后台友链列表
*/
func FlinkList(context *gin.Context) {
	var wg sync.WaitGroup
	//友情连接列表
	flink_list := models.FLink_List()
	flinkList := make([]*models.FriendLink, len(flink_list))
	for pos, id := range flink_list {
		wg.Add(1)
		go models.MultipleLoadFLink(id, pos, flinkList, &wg)
	}
	wg.Wait()

	flinkList = models.FilterNilFriendLink(flinkList)

	context.HTML(http.StatusOK, "flinklist.html", gin.H{
		"flinkList": flinkList,
	})
}

/**
添加一个类别
*/
func AddFlink(context *gin.Context) {

	var errcode int
	var errinfo string

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()
	webname, ok := context.GetPostForm("webname")
	if !ok || webname == "" {
		errcode = -1
		errinfo = "参数错误，请刷新后重试!"
		return
	}
	weburl, ok := context.GetPostForm("weburl")
	if !ok || weburl == "" {
		errcode = -1
		errinfo = "参数错误，请刷新后重试!"
		return
	}
	code := models.AddFlink(webname, weburl)
	if code < 0 {
		errcode = code
		if code == -3 {
			errinfo = "已存在该友链，不能重复添加!"
		} else {
			errinfo = "添加失败，请稍后重试!"
		}
	}
}

/**
删除一个友链
*/
func DelFlink(context *gin.Context) {
	var errcode int
	var errinfo string

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()

	id, ok := context.GetPostForm("id")
	if !ok {
		errcode = -1
		errinfo = "参数错误，请重试"
		return
	}

	code := models.DelFlink(id)

	if code < 0 {
		errcode = code
		errinfo = "删除失败，请刷新后重试"
		return
	}

	return

}

/**
跳转到修改友链页面
*/
func UpdateFlink(context *gin.Context) {

	fid := context.Param("fid")
	id, _ := strconv.Atoi(fid)
	flink := new(models.FriendLink)
	flink.Load(id)
	//这里应该加上异常处理

	context.HTML(http.StatusOK, "updateflink.html", gin.H{
		"flink": flink,
	})
}

/**
提交修改友链
*/
func UpFlink(context *gin.Context) {
	var errcode = -1
	var errinfo = "参数错误，请刷新后重试"

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})

	}()
	fid, ok := context.GetPostForm("id")
	if !ok {
		return
	}
	fsort, ok := context.GetPostForm("sort")
	if !ok {
		return
	}
	webname, ok := context.GetPostForm("webname")
	if !ok {
		return
	}

	weburl, ok := context.GetPostForm("weburl")
	if !ok {
		return
	}

	sort, err := strconv.Atoi(fsort)
	if err != nil {
		sort = 1
	}
	//处理sort，如果sort大于了当前类别数量，则sort＝类别数量
	flink_list := models.FLink_List()
	if sort > len(flink_list) {
		sort = len(flink_list)
	}

	//执行更新入库
	code := models.UpFlink(fid, webname, weburl, sort)
	if code < 0 {
		errcode = -2
		errinfo = "数据库异常，请稍后重试"
		return
	}
	errcode = 0
	errinfo = ""

	return
}
