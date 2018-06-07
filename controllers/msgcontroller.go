package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/models"
	"strconv"
	"sync"
	"math"
)


func MessageBorad(context *gin.Context){
	var wg sync.WaitGroup
	//文章列表
	tmpPage := context.Param("page")
	if tmpPage != ""{
		tmpPage = tmpPage[1:]
	}

	page,err := strconv.Atoi(tmpPage)
	if err != nil || page < 1{
		page = 1
	}

	allCnt := models.MsgCnt()
	pagesize := 10
	allPage := math.Ceil(float64(allCnt)/float64(pagesize))
	if float64(page) > allPage{
		page = 1
	}
	offset := (page - 1) * pagesize

	args := make(map[string]int)
	args["isshow"] = -1						//博客的显示控制 -1:全部;1:显示;0:隐藏
	args["pagesize"] = pagesize
	args["offset"] = offset
	//order 0:publish_time;1:read_count阅读量

	msg_ids := models.MsgList(args)
	msgList := make([]*models.Message,len(msg_ids))
	for pos,id := range msg_ids{
		wg.Add(1)
		models.MultipleLoadMessage(id,pos,msgList,&wg)
	}
	wg.Wait()

	//过滤空数据
	msgList = models.FilterNilMessage(msgList)

	pages := make([]int,0)
	for i := 1; i <= int(allPage);i++{
		pages = append(pages,i)
	}

	//读取中间件传来的参数
	tmp_gh,_ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["msgList"] = msgList
	gh["allPage"] = allPage
	gh["pages"] = pages
	gh["page"] = page
	gh["allCnt"] = allCnt
	gh["url"] = "/msg"

	context.HTML(http.StatusOK,"msg.html",gh)
}

/**
	添加一条留言
 */
func AddMsg(context *gin.Context){
	var errcode = -1
	var errinfo = "参数不全，请刷新该页面重试!"
	defer func(){
		context.JSON(http.StatusOK,gin.H{
			"errcode":errcode,
			"errinfo":errinfo,
		})
	}()

	name,ok := context.GetPostForm("name")
	if !ok{
		//参数错误
		return
	}
	if name == ""{
		errinfo = "姓名不能为空！"
		return
	}

	content,ok := context.GetPostForm("content")
	if !ok{
		//参数错误
		return
	}
	if content == ""{
		errinfo = "内容不能为空！"
		return
	}

	ret := models.AddMsg(name,content)
	if ret < 0{
		errinfo = "数据库异常，请刷新后重试！"
		return
	}
	errcode = 0
	errinfo = ""

	return
}


/**
	删除一个留言
 */