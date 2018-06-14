package controllers

import (
	"blog/models"
	"gopkg.in/gin-gonic/gin.v1"
	"math"
	"net/http"
	"strconv"
	"sync"
)

//首页
func Index(context *gin.Context) {
	var wg sync.WaitGroup
	//文章列表
	tmpPage := context.Param("page")
	if tmpPage != "" {
		tmpPage = tmpPage[1:]
	}

	curPage, err := strconv.Atoi(tmpPage)
	if err != nil || curPage < 1 {
		curPage = 1
	}

	allCnt := models.ArticleCnt() //文章总数量
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
	//order 0:publish_time;1:read_count阅读量

	article_ids := models.ArticleList(args)
	articleList := make([]*models.Article, len(article_ids))
	for pos, id := range article_ids {
		wg.Add(1)
		models.MultipleLoadArticle(id, pos, articleList, &wg)
	}
	wg.Wait()

	//过滤空数据
	articleList = models.FilterNilArticle(articleList)

	var perNum = 7
	pager := models.NewPage(int(allPage), curPage, perNum, "/index")

	//读取中间件传来的参数
	tmp_gh, _ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})
	gh["articleList"] = articleList
	gh["pager"] = pager

	context.HTML(http.StatusOK, "front/index.html", gh)
}

func MLogin(context *gin.Context) {
	var url = "/manage/index"
	loginname, _ := context.GetPostForm("loginname")
	password, _ := context.GetPostForm("password")

	user, login_ret := models.Login(loginname, password)
	if login_ret {
		//登录成功，写入session
		tmpSession, _ := context.Get("session")
		session := tmpSession.(*models.Session)

		session.SetSession("uid", user.Id)
		session.SetSession("name", user.Name)
		session.SetSession("nickname", user.Nickname)

	} else {
		login_ret = false
	}

	context.JSON(http.StatusOK, gin.H{
		"ret":  login_ret,
		"purl": url,
	})
}

/**
关于我
*/
func Myinfo(context *gin.Context) {
	tmp_gh, _ := context.Get("gh")
	gh := tmp_gh.(map[string]interface{})

	context.HTML(http.StatusOK, "myinfo.html", gh)
}

/**
登录
*/
func Login(context *gin.Context) {
	//跳转到登录页面
	context.HTML(http.StatusOK, "login.html", nil)

}
