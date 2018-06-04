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

	var name interface{}
	tmpSession,exists := context.Get("session")
	if exists{
		session := tmpSession.(*models.Session)
		if session.Has("uid"){
			name = session.GetSession("name")
		}else{
			name = "guest"
		}
	}
	//articleid
	articleid,ok := context.GetPostForm("aid")
	if !ok{
		//参数错误
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"参数不全，请刷新该页面重试",
		})
	}
	aid,err := strconv.Atoi(articleid)
	if err != nil{
		//参数错误
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"参数错误，请刷新该页面重试",
		})
	}
	content,ok := context.GetPostForm("content")
	if !ok{
		//参数错误
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"参数不全，请刷新该页面重试",
		})
	}
	if content == ""{
		context.JSON(http.StatusOK,gin.H{
			"errcode":-1,
			"errinfo":"内容不能为空",
		})
	}
	models.AddComment(aid,name,content)

	context.JSON(http.StatusOK,gin.H{
		"errcode":0,
		"errinfo":"",
	})
}



/**
	后台评论列表
 */
func CommentList(context *gin.Context){
	var wg sync.WaitGroup

	tmpArteId := context.Param("arteid")
	if tmpArteId != ""{
		tmpArteId = tmpArteId[1:]
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
		"cateid":arteid,
		"url": "/manage/comment/"+tmpArteId,

	})


}