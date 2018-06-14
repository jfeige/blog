package controllers

import (
	"blog/models"
	"gopkg.in/gin-gonic/gin.v1"
	"math"
	"net/http"
	"strconv"
	"sync"
)

/**
前台标签首页
*/
func TagIndex(context *gin.Context) {
	tagid := context.Param("tagid")

	t_id, err := strconv.Atoi(tagid)
	if err != nil || t_id == 0 {
		context.Redirect(http.StatusFound, "/")
	}

	tmpPage := context.Param("page")
	if tmpPage != "" {
		tmpPage = tmpPage[1:]
	}

	curPage, err := strconv.Atoi(tmpPage)
	if err != nil || curPage < 1 {
		curPage = 1
	}

	allCnt := models.ArticleByTagCnt(t_id) //该标签下文章总数量
	pagesize := models.BlogPageSize
	allPage := math.Ceil(float64(allCnt) / float64(pagesize))
	if float64(curPage) > allPage {
		curPage = 1
	}

	offset := (curPage - 1) * pagesize

	args := make(map[string]int)
	args["pagesize"] = pagesize
	args["offset"] = offset
	args["tid"] = t_id

	//根据标签id查询所有符合条件的文章
	var wg sync.WaitGroup

	article_ids := models.ArticleListByTag(args)
	articleList := make([]*models.Article, len(article_ids))
	for pos, id := range article_ids {
		wg.Add(1)
		go models.MultipleLoadArticle(id, pos, articleList, &wg)
	}
	wg.Wait()
	if len(articleList) == 0 {
		ToError(context, 404, "暂无文章,换个标签试试吧")

		context.Abort()
		return
	}

	var perNum = 7
	var url = "/tag/" + tagid
	pager := models.NewPage(int(allPage), curPage, perNum, url)

	//读取中间件传来的参数
	tmp_gh, _ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["articleList"] = articleList
	gh["pager"] = pager

	context.HTML(http.StatusOK, "front/index.html", gh)

}

/**
添加一个标签
*/
func AddTag(context *gin.Context) {
	var errcode int
	var errinfo string
	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()

	tagName, ok := context.GetPostForm("name")
	if !ok {
		errcode = -1
		errinfo = "参数错误，请重试!"
		return
	}
	code := models.AddTag(tagName)
	errcode = code

	if code < 0 {
		errinfo = "添加失败，请刷新后重试"
		if code == -2 {
			errinfo = "已存在该标签，不能重复添加"
			return
		}
	}
}

/**
删除一个标签
*/
func DelTag(context *gin.Context) {
	var errcode int
	var errinfo string
	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})
	}()

	tagid, ok := context.GetPostForm("id")
	if !ok {
		errcode = -1
		errinfo = "参数不全，请重试"
		return
	}
	code := models.DelTag(tagid)
	errcode = code
	if errcode < 0 {
		errinfo = "数据库异常，请稍后重试"
		return
	}

}
