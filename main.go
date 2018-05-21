package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"blog/controllers"
	"blog/models"
	"fmt"
	"sync"
)


func main() {
	err := models.InitBaseConfig("./conf/blog.ini")
	if err != nil{
		fmt.Println(err)
		return
	}
	router := gin.Default()

	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")
	//首页
	router.GET("/",FrongWare(),controllers.Index)
	//文章页面
	router.GET("/article/:artiId",FrongWare(), controllers.Article)
	//类别页面
	router.GET("/category/:cateId",FrongWare(), controllers.Index)
	//留言板

	http.ListenAndServe(":8080", router)
}


/**
前台页面专用中间件，用于读取页面右侧数据
 */
func FrongWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wg sync.WaitGroup

		//网站设置&&个人档案
		webSet := new(models.Webset)
		webSet.Load()

		//分类列表
		categroy_list := models.CategoryList()
		categoryList := make([]*models.Category,len(categroy_list))
		for pos,id := range categroy_list{
			wg.Add(1)
			models.MultipleLoadCategory(id,pos,categoryList,&wg)
		}
		wg.Wait()
		//近期文章
		//recentList := articleList
		//if len(articleList) >= 6{
		//	recentList = articleList[:6]
		//}

		//友情链接
		flink_list := models.FLink_List()
		flinkList := make([]*models.FriendLink,len(flink_list))
		for pos,id := range flink_list{
			wg.Add(1)
			models.MultipleLoadFLink(id,pos,flinkList,&wg)
		}
		wg.Wait()

		gh := make(map[string]interface{})
		gh["webSet"] = webSet
		gh["categoryList"] = categoryList
		gh["flinkList"] = flinkList

		c.Set("gh", gh)
		c.Next()
	}
}