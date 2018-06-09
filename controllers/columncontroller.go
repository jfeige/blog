package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"sync"
	"net/http"
	"blog/models"
	"strconv"
)



/**
	后台类别首页
 */
func ColumnManage(context *gin.Context){
	var wg sync.WaitGroup
	//分类列表
	column_ids := models.ColumnList()
	columnList := make([]*models.Column,len(column_ids))
	for pos,id := range column_ids{
		wg.Add(1)
		go models.MultipleLoadColumn(id,pos,columnList,&wg)
	}
	wg.Wait()
	columnList = models.FilterNilColumn(columnList)

	context.HTML(http.StatusOK,"column.html",gin.H{
		"columnList":columnList,
	})
}


/**
	添加一个类别
 */
func AddColumn(context *gin.Context){
	if context.Request.Method == "POST"{
		var errcode int
		var errinfo string

		defer func(){
			context.JSON(http.StatusOK,gin.H{
				"errcode":errcode,
				"errinfo":errinfo,
			})

		}()
		title,ok := context.GetPostForm("title")
		if !ok || title == ""{
			errcode = -1
			errinfo = "参数错误，请刷新后重试!"
			return
		}
		url,ok := context.GetPostForm("url")
		if !ok || url == ""{
			errcode = -1
			errinfo = "参数错误，请刷新后重试!"
			return
		}
		tmpTp,ok := context.GetPostForm("tp")
		if !ok || tmpTp == ""{
			errcode = -1
			errinfo = "参数错误，请刷新后重试!"
			return
		}
		tp,err := strconv.Atoi(tmpTp)
		if err != nil{
			errcode = -1
			errinfo = "参数错误，请刷新后重试!"
			return
		}
		if !models.InArray(tp,[]int{0,1}){
			errcode = -1
			errinfo = "参数错误，请刷新后重试!"
			return
		}
		code := models.AddColumn(title,url,tp)
		if code < 0{
			errcode = code
			if code == -2{
				errinfo = "已存相同栏目，不能重复添加!"
			}else{
				errinfo = "添加失败，请稍后重试!"
			}
		}
	}else{
		var wg sync.WaitGroup

		categroy_list := models.CategoryList()
		categoryList := make([]*models.Category,len(categroy_list))
		for pos,id := range categroy_list{
			wg.Add(1)
			go models.MultipleLoadCategory(id,pos,categoryList,&wg)
		}

		wg.Wait()


		context.HTML(http.StatusOK,"addcolumn.html",gin.H{
			"categoryList":categoryList,
		})

	}
}


/**
	跳转到修改类别页面
 */
func UpdateColumn(context *gin.Context){

	cid := context.Param("cid")
	id,_ := strconv.Atoi(cid)
	column := new(models.Column)
	column.Load(id)

	context.HTML(http.StatusOK,"updatecolumn.html",gin.H{
		"column":column,
	})
}



/**
	提交修改类别
 */
func UpColumn(context *gin.Context){
	var errcode int
	var errinfo string

	defer func(){
		context.JSON(http.StatusOK,gin.H{
			"errcode":errcode,
			"errinfo":errinfo,
		})

	}()
	cid,ok := context.GetPostForm("id");
	if !ok{
		errcode = -1
		errinfo = "参数错误，请刷新后重试"
	}
	column_sort,ok := context.GetPostForm("sort");
	if !ok{
		errcode = -1
		errinfo = "参数错误，请刷新后重试"
	}
	title,ok := context.GetPostForm("title");
	if !ok{
		errcode = -1
		errinfo = "参数错误，请刷新后重试"
	}
	url,ok := context.GetPostForm("url");
	if !ok{
		errcode = -1
		errinfo = "参数错误，请刷新后重试"
	}

	sort,err := strconv.Atoi(column_sort)
	if err != nil{
		sort = 1
	}
	//处理sort，如果sort大于了当前类别数量，则sort＝类别数量
	column_list := models.ColumnList()
	if sort > len(column_list){
		sort = len(column_list)
	}

	//执行更新入库
	code := models.UpColumn(cid,title,url,sort)
	if code < 0{
		errcode = -1
		errinfo = "数据库异常，请稍后重试"
	}
	return
}


/**
	删除一个栏目
 */
func DelColumn(context *gin.Context){
	var errcode int
	var errinfo string

	defer func(){
		context.JSON(http.StatusOK,gin.H{
			"errcode":errcode,
			"errinfo":errinfo,
		})
	}()

	cid,ok := context.GetPostForm("id")
	if !ok{
		errcode = -1
		errinfo = "参数错误，请重试"
		return
	}
	id,err := strconv.Atoi(cid)
	if err != nil{
		errcode = -1
		errinfo = "参数错误，请重试"
		return
	}

	code := models.DelColomn(id)

	if code < 0{
		errcode = code
		errinfo = "删除失败，请刷新后重试"
		return
	}

	return

}