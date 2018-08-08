package controllers

import (
	"blog/models"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"sync"
)

/**
添加一条回复
*/
func AddComment(context *gin.Context) {

	var errcode = -1
	var errinfo = "参数不全，请刷新该页面重试!"
	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()

	//articleid
	aid, ok := context.GetPostForm("aid")
	if !ok {
		//参数错误
		return
	}
	a_id, err := strconv.Atoi(aid)
	if err != nil {
		//参数错误
		return
	}
	tmp_tp, ok := context.GetPostForm("type")
	if !ok {
		return
	}
	if !models.InArray(tmp_tp, []string{"0", "1"}) {
		//参数错误
		return
	}
	tp, _ := strconv.Atoi(tmp_tp) //0:评论;1:回复

	var session *models.Session
	var name string

	tmpSession, ok := context.Get("session")
	if !ok && tp == 1 {
		errcode = -1
		errinfo = "只有站长才可以回复!"
		return
	}
	session = tmpSession.(*models.Session)
	if !session.Has("uid") && tp == 1 {
		errcode = -1
		errinfo = "只有站长才可以回复!"
		return
	}

	if session.Has("uid") {
		name = session.GetSession("nickname").(string)
	} else {
		name, ok = context.GetPostForm("name")
		if !ok {
			//参数错误
			return
		}
		if name == "" {
			errinfo = "姓名不能为空"
			return
		}
	}

	tmp_cid, ok := context.GetPostForm("cid")
	if !ok && tp == 1 {
		errinfo = "参数错误，请刷新后重试!"
		return
	}
	cid, err := strconv.Atoi(tmp_cid)
	if err != nil {
		errinfo = "参数错误，请刷新后重试!"
		return
	}

	content, ok := context.GetPostForm("content")
	if !ok {
		//参数错误
		return
	}
	if content == "" {
		errinfo = "内容不能为空"
		return
	}

	models.AddComment(a_id, tp, cid, name, content)

	errcode = 0
	errinfo = ""

	return
}

/**
后台评论列表
*/
func CommentList(context *gin.Context) {
	var wg sync.WaitGroup

	tmpArteId := context.Param("arteid")

	arteid, _ := strconv.Atoi(tmpArteId)
	//如果没有指定文章id，则读取所有评论，降序排列
	tmpPage := context.Param("page")
	if tmpPage == "" {
		tmpPage = "1"
	}
	curPage, err := strconv.Atoi(tmpPage)
	if err != nil || curPage < 1 {
		curPage = 1
	}

	allCnt := models.ManageCommentCnt(arteid) //评论总数量
	pagesize := 20
	allPage := math.Ceil(float64(allCnt) / float64(pagesize))
	if float64(curPage) > allPage {
		curPage = 1
	}

	offset := (curPage - 1) * pagesize

	args := make(map[string]int)
	args["arteid"] = arteid
	args["pagesize"] = pagesize
	args["offset"] = offset

	//评论页面，不分页
	comment_list := models.ManageCommentList(args)
	commentList := make([]*models.Comment, len(comment_list))
	for pos, id := range comment_list {
		wg.Add(1)
		go models.MultipleLoadComment(id, pos, commentList, &wg)
	}
	wg.Wait()
	pages := make([]int, 0)
	for i := 1; i <= int(allPage); i++ {
		pages = append(pages, i)
	}

	var perNum = 7
	pager := models.NewPage(int(allPage), curPage, perNum, "/manage/commentList/"+strconv.Itoa(arteid))

	context.HTML(http.StatusOK, "commentlist.html", gin.H{
		"commetList": commentList,
		"pager":      pager,
		"aid":        arteid,
	})

}

/*
	评论详情
*/
func CommentInfo(context *gin.Context) {

	cid := context.Param("cid")
	c_id, err := strconv.Atoi(cid)
	if err != nil {
		CommentList(context)
		context.Abort()
	}
	comment := new(models.Comment)
	err = comment.Load(c_id)
	if err != nil {
		CommentList(context)
		context.Abort()
	}

	context.HTML(http.StatusOK, "commentinfo.html", gin.H{
		"comment": comment,
	})
}

/**
删除一条评论
*/
func DelComment(context *gin.Context) {
	var errcode = -1
	var errinfo = "参数不全，请重试"
	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()

	c_id, ok := context.GetPostForm("cid")
	if !ok {
		return
	}
	cid, err := strconv.Atoi(c_id)
	if err != nil {
		return
	}
	a_id, ok := context.GetPostForm("aid")
	if !ok {
		return
	}
	aid, err := strconv.Atoi(a_id)
	if err != nil {
		return
	}

	//执行删除
	code := models.DelComment(aid, cid)

	if code < 0 {
		errcode = -2
		errinfo = "删除失败，请刷新后重试!"
		return
	}

	errcode = 0
	errinfo = ""

	return
}
