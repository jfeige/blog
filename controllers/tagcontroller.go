package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/models"
	"strconv"
	"sync"
	"math"
)


func TagIndex(context *gin.Context){
	tagid := context.Param("tagid")

	t_id,err := strconv.Atoi(tagid)
	if err != nil || t_id == 0{
		context.Redirect(http.StatusFound,"/")
	}

	tmpPage := context.Param("page")
	if tmpPage != ""{
		tmpPage = tmpPage[1:]
	}

	page,err := strconv.Atoi(tmpPage)
	if err != nil || page < 1{
		page = 1
	}

	allCnt := models.ArticleByTagCnt(t_id)			//该标签下文章总数量
	pagesize := models.BlogPageSize
	allPage := math.Ceil(float64(allCnt)/float64(pagesize))
	if float64(page) > allPage{
		page = 1
	}

	offset := (page - 1) * pagesize

	args := make(map[string]int)
	args["pagesize"] = pagesize
	args["offset"] = offset
	args["tid"] = t_id

	//根据标签id查询所有符合条件的文章
	var wg sync.WaitGroup

	article_ids := models.ArticleListByTag(args)
	articleList := make([]*models.Article,len(article_ids))
	for pos,id := range article_ids{
		wg.Add(1)
		go models.MultipleLoadArticle(id,pos,articleList,&wg)
	}
	wg.Wait()
	if len(articleList) == 0{
		errinfo := make(map[string]interface{})
		errinfo["errcode"] = ""
		errinfo["errinfo"] = "暂无文章,换个标签试试吧"
		ToError(context,errinfo)
		context.Abort()
		return
	}
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

	gh["url"] = "/tag/" + tagid

	context.HTML(http.StatusOK,"front/index.html",gh)


}


/**
	添加一个标签
 */
func AddTag(context *gin.Context){
	gh := make(map[string]interface{})
	defer func(){
		context.JSON(http.StatusOK,gh)
	}()

	tagName,ok := context.GetPostForm("name")
	if !ok{
		gh["errcode"] = -1
		gh["errinfo"] = "参数错误，请重试!"
		return
	}
	errcode := models.AddTag(tagName)
	gh["errcode"] = errcode
	if errcode < 0{
		gh["errinfo"] = "添加失败，请刷新后重试"
		if errcode == -2{
			gh["errinfo"] = "已存在该标签，不能重复添加"
			return
		}
	}
}

/**
	删除一个标签
 */
func DelTag(context *gin.Context){
	gh := make(map[string]interface{})
	defer func(){
		context.JSON(http.StatusOK,gh)
	}()

	tagid,ok := context.GetPostForm("id")
	if !ok{
		gh["errcode"] = -1
		gh["errinfo"] = "参数不全，请重试"
		return
	}
	errcode := models.DelTag(tagid)
	if errcode < 0{
		gh["errcode"] = -2
		gh["errinfo"] = "数据库异常，请稍后重试"
		return
	}

	gh["errcode"] = 0
}