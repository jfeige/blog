package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	log "github.com/alecthomas/log4go"
	"fmt"
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
	key := "article_cnt:" + strconv.Itoa(cid)
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select count(1) from b_article "
		if cid > 0{
			sql = "select count(1) from b_article where categoryid=? "
			pargs = append(pargs,cid)
		}
		stmt,err := db.Prepare(sql)
		if err != nil{
			log.Error(fmt.Sprintf("db.Prepare has error:",err))
			return	0
		}
		defer stmt.Close()
		row := stmt.QueryRow(pargs...)
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

	if !InArray(isshow,[]int{-1,0,1}){
		isshow = -1
	}

	key := "articleList:"  + strconv.Itoa(cateid)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select id,publish_time from b_article order by publish_time desc "
		if cateid > 0{
			sql = "select id,publish_time from b_article where categoryid=? order by publish_time desc "
			pargs = append(pargs,cateid)
		}
		stmt,err := db.Prepare(sql)
		if err != nil{
			log.Error(fmt.Sprintf("db.Prepare has error:",err))
			return	list
		}
		defer stmt.Close()

		rows,err := stmt.Query(pargs...)
		if err != nil{
			log.Error(fmt.Sprintf("stmt.Query has error:",err))
			return	list
		}
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,publish_time int
		for rows.Next(){
			rows.Scan(&id,&publish_time)
			rargs = append(rargs,publish_time,id)
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
func UpdateArticleInfo(id int,title,content string)int{

	db := conn.GetMysqlConn()
	sql := "update b_article set title=?,content=? where id=?"
	stmt,err := db.Prepare(sql)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	defer stmt.Close()
	_,err = stmt.Exec(title,content,id)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}

	//更新缓存
	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "article:" + strconv.Itoa(id)

	_,err = rconn.Do("DEL",key)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
	}
	return 0
}