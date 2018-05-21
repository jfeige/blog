package models

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
)

/**
	根据文章id获取评论id列表
 */
func CommentList(articleid int)[]int{
	list := make([]int,0)
	key := "commentList"
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{
		sql := "select id,atime from b_comment order by atime desc"
		db := conn.GetMysqlConn()
		rows,err := db.Query(sql)
		if err != nil{
			log.Error("db.Query has error:%v",err)
			return list
		}
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,atime int
		for rows.Next(){
			err := rows.Scan(&id,&atime)
			if err != nil{
				log.Error(fmt.Sprintf("rows.Scan has error:%v",err))
				continue
			}
			rargs = append(rargs,atime,id)
		}
		if len(rargs) > 1{
			rconn.Send("ZADD",rargs...)
		}
	}
	cnt,_ := redis.Int(rconn.Do("ZCARD",key))
	if cnt > 0{
		list,err = redis.Ints(rconn.Do("ZREVRANGE",key,0,-1))
		if err != nil{
			log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
			return list
		}
	}
	return list
}
