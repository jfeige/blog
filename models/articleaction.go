package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	log "github.com/alecthomas/log4go"
	"fmt"
)

/**
	博客列表
 */
func ArticleList(args map[string]interface{})[]int{
	list := make([]int,0)
	pagesize := blog_pagesize
	page := args["page"].(int)
	isshow := args["isshow"].(int)	//博客的显示控制 -1:全部;1:显示;0:隐藏

	if !InArray(isshow,[]int{-1,0,1}){
		isshow = -1
	}

	key := "articleList:" + strconv.Itoa(isshow)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select id,publish_time from b_article order by publish_time desc "
		if isshow > -1{
			sql = "select id,publish_time from b_article order by publish_time where isshow=? where desc "
			pargs = append(pargs,isshow)
		}
		stmt,err := db.Prepare(sql)
		if err != nil{
			log.Error(fmt.Sprintf("db.Prepare has error:",err))
			return	list
		}
		defer stmt.Close()

		rows,err := stmt.Query()
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

	cnt,_ := redis.Int(rconn.Do("ZCARD",key))
	if cnt > 0{
		offset := (page - 1) * pagesize
		limit := offset + pagesize - 1
		var args = []interface{}{key, offset, limit}
		list,err = redis.Ints(rconn.Do("ZREVRANGE",args...))
		if err != nil{
			log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
			return list
		}
	}
	return list
}

