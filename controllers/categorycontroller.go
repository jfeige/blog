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
类别首页(前台)
*/
func CategoryFront(context *gin.Context) {
	var wg sync.WaitGroup

	//类别id.如果tmpCateId 为空,先不考虑这种情况
	tmpCateId := context.Param("cateid")
	cateid, err := strconv.Atoi(tmpCateId)
	if err != nil || cateid < 1 {
		cateid = 1
	}

	tmpPage := context.Param("page")
	if tmpPage != "" {
		tmpPage = tmpPage[1:]
	}

	curPage, err := strconv.Atoi(tmpPage)
	if err != nil || curPage < 1 {
		curPage = 1
	}

	allCnt := models.ArticleCnt(cateid) //文章总数量
	pagesize := models.BlogPageSize
	allPage := math.Ceil(float64(allCnt) / float64(pagesize))
	if float64(curPage) > allPage {
		curPage = 1
	}

	offset := (curPage - 1) * pagesize

	args := make(map[string]int)
	args["isshow"] = -1 //博客的显示控制 -1:全部;1:显示;0:隐藏
	args["pagesize"] = pagesize
	args["offset"] = offset
	args["cateid"] = cateid

	article_ids := models.ArticleList(args)
	articleList := make([]*models.Article, len(article_ids))
	for pos, id := range article_ids {
		wg.Add(1)
		go models.MultipleLoadArticle(id, pos, articleList, &wg)
	}
	wg.Wait()

	var url string
	if cateid == 0 {
		url = "/index"
	} else {
		url = "/category/" + strconv.Itoa(cateid)
	}
	var perNum = 7
	pager := models.NewPage(int(allPage), curPage, perNum, url)

	//读取中间件传来的参数
	tmp_gh, _ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["articleList"] = articleList
	gh["pager"] = pager

	context.HTML(http.StatusOK, "front/index.html", gh)
}

/**
后台类别首页
*/
func CategoryManage(context *gin.Context) {
	var wg sync.WaitGroup
	//分类列表
	categroy_list := models.CategoryList()
	categoryList := make([]*models.Category, len(categroy_list))
	for pos, id := range categroy_list {
		wg.Add(1)
		go models.MultipleLoadCategory(id, pos, categoryList, &wg)
	}

	wg.Wait()

	context.HTML(http.StatusOK, "category.html", gin.H{
		"categoryList": categoryList,
	})
}

/**
删除一个类别
*/
func DelCategory(context *gin.Context) {
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

	code := models.DelCatetory(id)

	if code < 0 {
		errcode = code
		errinfo = "删除失败，请刷新后重试"
		return
	}

	return

}

/**
添加一个类别
*/
func AddCategory(context *gin.Context) {

	var errcode int
	var errinfo string

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})

	}()
	cate_name, ok := context.GetPostForm("name")
	if !ok || cate_name == "" {
		errcode = -1
		errinfo = "参数错误，请刷新后重试!"
		return
	}
	code := models.AddCategory(cate_name)
	if code < 0 {
		errcode = code
		if code == -3 {
			errinfo = "已存在该标签，不能重复添加!"
		} else {
			errinfo = "添加失败，请稍后重试!"
		}
	}
}

/**
跳转到修改类别页面
*/
func UpdateCatetory(context *gin.Context) {

	cateid := context.Param("cateid")
	id, _ := strconv.Atoi(cateid)
	category := new(models.Category)
	category.Load(id)

	context.HTML(http.StatusOK, "updatecategory.html", gin.H{
		"cateotry": category,
	})
}

/**
提交修改类别
*/
func UpCatetory(context *gin.Context) {
	var errcode = -1
	var errinfo = "参数错误，请刷新后重试"

	defer func() {
		context.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"errinfo": errinfo,
		})

	}()
	cateid, ok := context.GetPostForm("id")
	if !ok {
		return
	}
	catesort, ok := context.GetPostForm("sort")
	if !ok {
		return
	}
	catename, ok := context.GetPostForm("name")
	if !ok {
		return
	}

	sort, err := strconv.Atoi(catesort)
	if err != nil {
		sort = 1
	}
	//处理sort，如果sort大于了当前类别数量，则sort＝类别数量
	categroy_list := models.CategoryList()
	if sort > len(categroy_list) {
		sort = len(categroy_list)
	}

	//执行更新入库
	code := models.UpCatetory(cateid, catename, sort)
	if code < 0 {
		errcode = -2
		errinfo = "数据库异常，请稍后重试"
		return
	}

	errcode = 0
	errinfo = ""
	return
}
