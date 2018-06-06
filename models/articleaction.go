package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	log "github.com/alecthomas/log4go"
	"fmt"
	"time"
	"strings"
)


/**
	文章总数
 */
func ArticleCnt(cateid ...int)int{
	var cid,cnt int
	if len(cateid) > 0{
		cid = cateid[0]
	}
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "articleList:cnt:" + strconv.Itoa(cid)
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select count(1) from b_article "
		if cid > 0{
			sql = "select count(1) from b_article where categoryid=? "
			pargs = append(pargs,cid)
		}
		row := db.QueryRow(sql,pargs...)
		err = row.Scan(&cnt)
		if err != nil{
			log.Error(fmt.Sprintf("row.Scan has error:",err))
			return	0
		}
		err = rconn.Send("set",key,cnt)
		if err != nil{
			log.Error(fmt.Sprintf("rconn.Send has error:",err))
			return	cnt
		}
	}

	cnt,err := redis.Int(rconn.Do("get",key))
	if err != nil{
		log.Error(fmt.Sprintf("redis.Int has error:%v",err))
		return 0
	}
	return cnt

}


/**
	文章列表
 */
func ArticleList(args map[string]int)[]int{
	list := make([]int,0)
	pagesize,ok := args["pagesize"]
	if !ok{
		pagesize = 10
	}
	isshow,ok := args["isshow"]
	if !ok{
		isshow = -1
	}
	offset,ok := args["offset"]
	if !ok{
		offset = 0
	}
	cateid,ok := args["cateid"]		//分类id
	if !ok || cateid < 0{
		cateid = 0
	}
	order,ok := args["order"]		//order 0:publish_time;1:read_count阅读量
	if !ok{
		order = 0
	}

	if !InArray(isshow,[]int{-1,0,1}){
		isshow = -1
	}
	if !InArray(order,[]int{0,1}){
		order = 0
	}

	key := "articleList:"  + strconv.Itoa(cateid) + "|" + strconv.Itoa(order)
	fmt.Println(key)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select id,publish_time,read_count from b_article order by publish_time desc "
		if cateid > 0{
			sql = "select id,publish_time,read_count from b_article where categoryid=? order by publish_time desc "
			pargs = append(pargs,cateid)
		}

		rows,err := db.Query(sql,pargs...)
		if err != nil{
			log.Error(fmt.Sprintf("stmt.Query has error:",err))
			return	list
		}
		defer rows.Close()
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,publish_time,read_count int
		for rows.Next(){
			rows.Scan(&id,&publish_time,&read_count)
			if order == 0{
				rargs = append(rargs,publish_time,id)
			}else{
				rargs = append(rargs,read_count,id)
			}
		}
		if len(rargs) > 1{
			rconn.Send("ZADD",rargs...)
		}
	}

	limit := offset + pagesize - 1
	var params = []interface{}{key, offset, limit}
	list,err = redis.Ints(rconn.Do("ZREVRANGE",params...))
	if err != nil{
		log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
		return list
	}
	return list
}


/**
	修改文章内容
 */
func UpdateArticleInfo(a_id int,title,content,tagids string)int{

	db := conn.GetMysqlConn()
	sql := "update b_article set title=?,content=? where id=?"
	tx,err := db.Begin()
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	defer tx.Rollback()
	_,err = tx.Exec(sql,title,content,a_id)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	//首先移除所有的标签，然后再插入
	sql = "delete from b_actmapptags where a_id=?";
	_,err = tx.Exec(sql,a_id)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	tags := strings.Split(tagids,",")
	for _,id := range tags{
		sql = "insert into b_actmapptags(a_id,t_id) values(?,?)"
		_,err = tx.Exec(sql,a_id,id)
		if err != nil{
			log.Error("UpdateArticleInfo has error:%v",err)
			return -2
		}
	}

	tx.Commit()

	//更新缓存
	rconn := conn.pool.Get()
	defer rconn.Close()


	keys := make([]interface{},0)
	keys = append(keys,"article:" + strconv.Itoa(a_id))
	keys = append(keys,"tagids:" + strconv.Itoa(a_id))

	_,err = rconn.Do("DEL",keys...)

	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
	}
	return 0

}

/**
	添加一篇文章
 */
func AddArticle(cateid int,title,content,tagids string) int{
	db := conn.GetMysqlConn()
	sql := "insert into b_article(title,content,categoryid,publish_time,publish_date) values(?,?,?,?,?)"

	tx,err := db.Begin()
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	defer tx.Rollback()

	result,err := tx.Exec(sql,title,content,cateid,time.Now().Unix(),time.Now().Format("20060102"))
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	a_id,err := result.LastInsertId()
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}

	tags := strings.Split(tagids,",")
	for _,id := range tags{
		sql = "insert into b_actmapptags(a_id,t_id) values(?,?)"
		_,err = tx.Exec(sql,a_id,id)
		if err != nil{
			log.Error("AddArticle has error:%v",err)
			return -2
		}
	}


	sql = "update b_category set article_cnt=article_cnt+1 where id=?"
	_,err = tx.Exec(sql,cateid)
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	tx.Commit()
	//更新缓存
	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "articleList:*"
	DelKeys(key)

	keys := make([]interface{},0)
	keys = append(keys,"category:" + strconv.Itoa(cateid))
	keys = append(keys,"category:0")

	_,err = rconn.Do("DEL",keys...)
	if err != nil{
		log.Error("AddArticle has error:%v",err)
	}
	return 0
}


/**
	更新文章浏览数
 */
func UpdateReadCnt(a_id,cnt int){
	sql := "update b_article set read_count = read_count + ? where id=?"
	db := conn.GetMysqlConn()

	_,err = db.Exec(sql,cnt,a_id)
	if err != nil{
		log.Error("UpdateReadCnt has error:%v",err)
		return
	}

	key := "article:" + strconv.Itoa(a_id)

	rconn := conn.pool.Get()
	defer rconn.Close()

	rconn.Do("HINCRBY",key,"read_count",cnt)
}


/**
	删除文章
 */
func DelArticle(aid,cateid int)int{
	db := conn.GetMysqlConn()
	sql := "call delArticle(?,?)"
	row := db.QueryRow(sql,aid,cateid)

	var errcode int
	err = row.Scan(&errcode)
	if err != nil{
		log.Error("DelArticle has error:%v",err)
		return -2
	}
	if errcode < 0{
		return errcode
	}
	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "articleList:*"
	DelKeys(key)
	key = "commentList:*"
	DelKeys(key)

	keys := make([]interface{},0)
	keys = append(keys,"article:" + strconv.Itoa(aid))
	keys = append(keys,"category:" + strconv.Itoa(cateid))
	keys = append(keys,"tagids:" + strconv.Itoa(aid))

	_,err = rconn.Do("DEL",keys...)
	if err != nil{
		log.Error("DelArticle has error:%v",err)
	}
	return 0
}