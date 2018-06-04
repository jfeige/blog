package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"net/http"
	"strconv"
	"sync"
	"math"
)


/**
	添加一条回复	ajax请求
 */
func AddComment(context *gin.Context){

	var errcode = -1
	var errinfo = "参数不全，请刷新该页面重试"
	defer func(){
		context.JSON(http.StatusOK,gin.H{
			"errcode":errcode,
			"errinfo":errinfo,
		})
	}()

	//articleid
	aid,ok := context.GetPostForm("aid")
	if !ok{
		//参数错误
		return
	}
	a_id,err := strconv.Atoi(aid)
	if err != nil{
		//参数错误
		return
	}
	name,ok := context.GetPostForm("name")
	if !ok{
		//参数错误
		return
	}
	if name == ""{
		errinfo = "姓名不能为空"
		return
	}
	content,ok := context.GetPostForm("content")
	if !ok{
		//参数错误
		return
	}
	if content == ""{
		errinfo = "内容不能为空"
		return
	}
	models.AddComment(a_id,name,content)

	errcode = 0
	errinfo = ""

	return
}



/**
	后台评论列表
 */
func CommentList(context *gin.Context){
	var wg sync.WaitGroup

	tmpArteId := context.Param("arteid")
	if tmpArteId != ""{
		tmpArteId = tmpArteId[1:]
	}else{
		tmpArteId = "0"
	}

	arteid,_ := strconv.Atoi(tmpArteId)
	//如果没有指定文章id，则读取所有评论，降序排列

	tmpPage,ok := context.GetQuery("page")
	if !ok{
		tmpPage = "1"
	}
	page,err := strconv.Atoi(tmpPage)
	if err != nil || page < 1{
		page = 1
	}

	allCnt := models.ManageCommentCnt(arteid)			//评论总数量
	pagesize := 20
	allPage := math.Ceil(float64(allCnt)/float64(pagesize))
	if float64(page) > allPage{
		page = 1
	}

	offset := (page - 1) * pagesize

	args := make(map[string]int)
	args["arteid"] = arteid
	args["pagesize"] = pagesize
	args["offset"] = offset

	//评论页面，不分页
	comment_list := models.ManageCommentList(args)
	commentList := make([]*models.Comment,len(comment_list))
	for pos,id := range comment_list{
		wg.Add(1)
		models.MultipleLoadComment(id,pos,commentList,&wg)
	}
	wg.Wait()
	pages := make([]int,0)
	for i := 1; i <= int(allPage);i++{
		pages = append(pages,i)
	}


	context.HTML(http.StatusOK,"commentlist.html",gin.H{
		"commetList":commentList,
		"allPage" : int(allPage),
		"pages": pages,
		"page": page,
		"prevPage":page-1,
		"nextPage":page+1,
		"aid":arteid,
		"url": "/manage/comment/"+tmpArteId,

	})

}

/*
	评论详情
 */
func CommentInfo(context *gin.Context){

	cid := context.Param("cid")
	c_id,err := strconv.Atoi(cid)
	if err != nil{
		CommentList(context)
		context.Abort()
	}
	comment := new(models.Comment)
	err = comment.Load(c_id)
	if err != nil{
		CommentList(context)
		context.Abort()
	}

	context.HTML(http.StatusOK,"commentinfo.html",gin.H{
		"comment":comment,
	})
}

/**
	删除一条评论
 */
func DelComment(context *gin.Context){
	var errcode int
	var errinfo string
	defer func(){
		context.JSON(http.StatusOK,gin.H{
			"errcode":errcode,
			"errinfo":errinfo,
		})
	}()

	c_id,ok := context.GetPostForm("cid")
	if !ok{
		errcode = -1
		errinfo = "参数不全，请重试"
		return
	}
	cid,err := strconv.Atoi(c_id)
	if err != nil{
		errcode = -1
		errinfo = "参数错误，请重试"
		return
	}
	a_id,ok := context.GetPostForm("aid")
	if !ok{
		errcode = -1
		errinfo = "参数不全，请重试"
		return
	}
	aid,err := strconv.Atoi(a_id)
	if err != nil{
		errcode = -1
		errinfo = "参数错误，请重试"
		return
	}

	//执行删除
	code := models.DelComment(aid,cid)

	if code < 0{
		errcode = -2
		errinfo = "删除失败，请刷新后重试!"
		return
	}

	return
}