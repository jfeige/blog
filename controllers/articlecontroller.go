package controllers

import (
	"blog/models"
	log "github.com/alecthomas/log4go"
	"gopkg.in/gin-gonic/gin.v1"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
)

/**
前台文章详情
*/
func Article(context *gin.Context) {
	id := context.Param("arteid")
	a_id, err := strconv.Atoi(id)
	if err != nil || a_id <= 0 {
		//参数错误，跳转到首页
		ErrArgs(context)
		return
	}
	//根据文章id读取
	article := new(models.Article)
	err = article.Load(a_id)
	if err != nil {
		//数据错误或者id不正确
		ToError(context, 500, "服务器开小差了，稍后再试吧!")
		return
	}
	//累计浏览量
	go AddReadCnt(a_id)

	var wg sync.WaitGroup

	//评论页面，不分页
	comment_list := models.CommentList(a_id)
	commentList := make([]*models.Comment, len(comment_list))
	for pos, id := range comment_list {
		wg.Add(1)
		models.MultipleLoadComment(id, pos, commentList, &wg)
	}
	wg.Wait()

	tmp_gh, _ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["commList"] = commentList
	gh["cns"] = len(commentList)
	gh["article"] = article

	tmpSession, ok := context.Get("session")
	if ok {
		session := tmpSession.(*models.Session)
		gh["session"] = session
	}

	context.HTML(http.StatusOK, "article.html", gh)
}

/**
后台文章详情
*/
func ArticleInfo(context *gin.Context) {
	id := context.Param("arteid")
	artiId, err := strconv.Atoi(id)
	if err != nil || artiId <= 0 {
		//参数错误，跳转到首页
		ErrArgs(context)
		return
	}
	//根据文章id读取
	article := new(models.Article)
	err = article.Load(artiId)
	if err != nil {
		//数据错误或者id不正确
		ToError(context, 500, "服务器开小差了，稍后再试吧!")
		return
	}
	var wg sync.WaitGroup
	categroy_list := models.CategoryList()
	categoryList := make([]*models.Category, len(categroy_list))
	for pos, id := range categroy_list {
		wg.Add(1)
		go models.MultipleLoadCategory(id, pos, categoryList, &wg)
	}

	//标签
	tag_list := models.TagList()
	tagList := make([]*models.Tag, len(tag_list))
	for pos, id := range tag_list {
		wg.Add(1)
		go models.MultipleLoadTag(id, pos, tagList, &wg)
	}

	wg.Wait()

	context.HTML(http.StatusOK, "articleinfo.html", gin.H{
		"article":      article,
		"categoryList": categoryList,
		"tagList":      tagList,
	})
}

/**
后台文章列表
*/
func ArticleList(context *gin.Context) {
	var wg sync.WaitGroup

	c_id := context.Param("cateid")

	tmpPage := context.Param("page")
	if tmpPage == "" {
		tmpPage = "1"
	}
	curPage, err := strconv.Atoi(tmpPage)
	if err != nil || curPage < 1 {
		curPage = 1
	}
	cateid, err := strconv.Atoi(c_id)
	if err != nil || cateid <= 0 {
		cateid = 0
	}

	allCnt := models.ArticleCnt(cateid) //文章总数量
	pagesize := 20
	allPage := math.Ceil(float64(allCnt) / float64(pagesize))
	if float64(curPage) > allPage {
		curPage = 1
	}

	offset := (curPage - 1) * pagesize

	args := make(map[string]int)
	args["cateid"] = cateid
	args["pagesize"] = pagesize
	args["offset"] = offset

	article_ids := models.ArticleList(args)
	articleList := make([]*models.Article, len(article_ids))
	for pos, id := range article_ids {
		wg.Add(1)
		go models.MultipleLoadArticle(id, pos, articleList, &wg)
	}
	wg.Wait()

	var perNum = 7
	pager := models.NewPage(int(allPage), curPage, perNum, "/manage/articleList/"+strconv.Itoa(cateid))

	context.HTML(http.StatusOK, "manage/articlelist.html", gin.H{
		"articleList": articleList,
		"pager":       pager,
		"cateid":      cateid,
	})
}

/**
提交文章修改
*/
func UpfateArticleInfo(context *gin.Context) {

	var errcode = -1
	var errinfo = "参数错误，请重试!"

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()

	id, ok := context.GetPostForm("id")
	if !ok {
		return
	}
	arteid, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	title, ok := context.GetPostForm("title")
	if !ok || title == "" {
		return
	}

	content, ok := context.GetPostForm("content")
	if !ok || content == "" {
		return
	}

	tagids, ok := context.GetPostForm("checkID")
	if !ok || tagids == "" {
		return
	}

	code := models.UpdateArticleInfo(arteid, title, content, tagids)
	if code < 0 {
		errcode = -2
		errinfo = "修改失败，请稍后重试!"
		return
	}
	errcode = 0
	errinfo = ""
	return
}

/**
添加文章
*/
func AddArticle(context *gin.Context) {
	if context.Request.Method == "POST" {
		var errcode = -1
		var errinfo = "参数错误，请重试!"

		defer func() {
			context.JSON(http.StatusOK, gin.H{
				"errcode": errcode,
				"errinfo": errinfo,
			})
		}()
		id, ok := context.GetPostForm("cateid")
		if !ok {
			return
		}
		cateid, err := strconv.Atoi(id)
		if err != nil {
			return
		}

		title, ok := context.GetPostForm("title")
		if !ok || title == "" {
			return
		}

		content, ok := context.GetPostForm("content")
		if !ok || content == "" {
			return
		}

		//
		tagids, ok := context.GetPostForm("checkID")
		if !ok || tagids == "" {
			return
		}
		tmpSession, _ := context.Get("session")
		session := tmpSession.(*models.Session)
		user := session.GetSession("nickname").(string)

		code := models.AddArticle(cateid, title, user, content, tagids)
		if code < 0 {
			errcode = -2
			errinfo = "添加失败，请刷新后重试!"
			return
		}
		errcode = 0
		errinfo = ""
	} else {
		var wg sync.WaitGroup
		//获取类别列表
		categroy_list := models.CategoryList()
		categoryList := make([]*models.Category, len(categroy_list))
		for pos, id := range categroy_list {
			wg.Add(1)
			go models.MultipleLoadCategory(id, pos, categoryList, &wg)
		}

		//标签
		tag_list := models.TagList()
		tagList := make([]*models.Tag, len(tag_list))
		for pos, id := range tag_list {
			wg.Add(1)
			go models.MultipleLoadTag(id, pos, tagList, &wg)
		}

		wg.Wait()

		context.HTML(http.StatusOK, "addarticle.html", gin.H{
			"categoryList": categoryList,
			"tagList":      tagList,
		})
	}
}

/**
删除文章
*/
func DelArticle(context *gin.Context) {
	var errcode = -1
	var errinfo = "参数错误，请重试!"

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()
	a_id, ok := context.GetPostForm("arteid")
	if !ok {
		return
	}
	aid, err := strconv.Atoi(a_id)
	if err != nil {
		return
	}
	cate_id, ok := context.GetPostForm("cateid")
	if !ok {
		return
	}
	cateid, err := strconv.Atoi(cate_id)
	if err != nil {
		return
	}
	code := models.DelArticle(aid, cateid)
	if code < 0 {
		errcode = -2
		errinfo = "删除失败，请刷新后重试!"
		return
	}
	errcode = 0
	errinfo = ""
	return
}

func UpImg(context *gin.Context) {
	var success = 0
	var message = "上传失败,请刷新页面后重试"
	var url string
	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"success": success,
			"message": message,
			"url":     url,
		})
	}()
	img, header, err := context.Request.FormFile("editormd-image-file")
	if err != nil {
		log.Error("UpImg has error:%v", err)
		return
	}
	filename := header.Filename
	out, err := os.Create("static/images/" + filename)
	if err != nil {
		log.Error("Test has error:%v", err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, img)
	if err != nil {
		log.Error("Test has error:%v", err)
		return
	}
	success = 1
	message = "上传成功"
	url = "/static/images/" + filename
	return
}
