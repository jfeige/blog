package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"blog/models"
	"sync"
	"strconv"
	"net/http"
	"math"
)




//文章详情
func ArticleInfo(context *gin.Context){
	id := context.Param("arteid")
	artiId,err := strconv.Atoi(id)
	if err != nil || artiId <= 0{
		//参数错误，跳转到首页
		context.Redirect(0,"/")
	}
	//根据文章id读取
	article := new(models.Article)
	err = article.Load(artiId)
	if err != nil{
		//数据错误或者id不正确

	}


	context.HTML(http.StatusOK,"articleinfo.html",gin.H{
		"article":article,
	})
}
/**
	后台文章列表
 */
func ArticleList(context *gin.Context){
	var wg sync.WaitGroup

	id := context.Param("cateid")
	if id != ""{
		id = id[1:]
	}
	tmpPage,ok := context.GetQuery("page")
	if !ok{
		tmpPage = "1"
	}
	page,err := strconv.Atoi(tmpPage)
	if err != nil || page < 1{
		page = 1
	}
	cateid,err := strconv.Atoi(id)
	if err != nil || cateid <= 0{
		cateid = 0
	}

	allCnt := models.ArticleCnt(cateid)			//文章总数量
	pagesize := 20
	allPage := math.Ceil(float64(allCnt)/float64(pagesize))
	if float64(page) > allPage{
		page = 1
	}

	offset := (page - 1) * pagesize

	args := make(map[string]int)
	args["cateid"] = cateid
	args["pagesize"] = pagesize
	args["offset"] = offset

	article_ids := models.ArticleList(args)
	articleList := make([]*models.Article,len(article_ids))
	for pos,id := range article_ids {
		wg.Add(1)
		models.MultipleLoadArticle(id,pos,articleList,&wg)
	}
	wg.Wait()
	//暂不考虑分页显示

	pages := make([]int,0)
	for i := 1; i <= int(allPage);i++{
		pages = append(pages,i)
	}

	context.HTML(http.StatusOK,"manage/articlelist.html",gin.H{
		"articleList":articleList,
		"allPage" : int(allPage),
		"pages": pages,
		"page": page,
		"prevPage":page-1,
		"nextPage":page+1,
		"cateid":id,
		"url": "/manage/article/"+id,
	})
}


/**
	提交文章修改
 */
func UpfateArticleInfo(context *gin.Context){

	var errcode int
	var errinfo string

	defer func(){
		context.JSON(http.StatusOK,gin.H{
			"errcode":errcode,
			"errinfo":errinfo,
		})
	}()

	id,ok := context.GetPostForm("id")
	if !ok{
		errcode = -1
		errinfo = "参数错误，请重试!"
	}
	arteid,err := strconv.Atoi(id)
	if err != nil{
		errcode = -1
		errinfo = "参数错误，请重试!"
	}

	title,ok := context.GetPostForm("title")
	if !ok || title == ""{
		errcode = -1
		errinfo = "参数错误，请重试!"
	}

	content,ok := context.GetPostForm("content")
	if !ok || content == ""{
		errcode = -1
		errinfo = "参数错误，请重试!"
	}

	code := models.UpdateArticleInfo(arteid,title,content)
	if code < 0{
		errcode = -2
		errinfo = "修改失败，请稍后重试!"
	}

	return
}