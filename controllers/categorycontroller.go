package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"sync"
	"strconv"
	"math"
	"net/http"
	"blog/models"
)

/**
	类别首页
 */
func CategoryIndex(context *gin.Context){
	var wg sync.WaitGroup

	//类别id.如果tmpCateId 为空,先不考虑这种情况
	tmpCateId := context.Param("cateid")
	cateid,err := strconv.Atoi(tmpCateId)
	if err != nil || cateid < 1{
		cateid = 1
	}

	tmpPage := context.Param("page")
	page,err := strconv.Atoi(tmpPage)
	if err != nil || page < 1{
		page = 1
	}

	allCnt := models.ArticleCnt(cateid)			//文章总数量
	pagesize := models.BlogPageSize
	allPage := math.Ceil(float64(allCnt)/float64(pagesize))
	if float64(page) > allPage{
		page = 1
	}

	offset := (page - 1) * pagesize

	args := make(map[string]int)
	args["isshow"] = -1						//博客的显示控制 -1:全部;1:显示;0:隐藏
	args["pagesize"] = pagesize
	args["offset"] = offset
	args["cateid"] = cateid

	article_ids := models.ArticleList(args)
	articleList := make([]*models.Article,len(article_ids))
	for pos,id := range article_ids{
		wg.Add(1)
		models.MultipleLoadArticle(id,pos,articleList,&wg)
	}
	wg.Wait()

	pages := make([]int,0)
	for i := 1; i <= int(allPage);i++{
		pages = append(pages,i)
	}

	//读取中间件传来的参数
	tmp_gh,_ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["articleList"] = articleList
	gh["allPage"] = allPage
	gh["pages"] = pages
	gh["page"] = page
	if cateid == 0{
		gh["url"] = "/index"
	}else{
		gh["url"] = "/category/" + strconv.Itoa(cateid)
	}

	context.HTML(http.StatusOK,"index.html",gh)
}