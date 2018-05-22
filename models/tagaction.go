package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"fmt"
)

func TagList()[]int{
	list := make([]int,0)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "taglist"
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		sql := "select id from b_tag order by id asc"
		rows,err := db.Query(sql)
		if err != nil{
			log.Error("db.Query() has error:%v",err)
			return list
		}
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id int
		for rows.Next(){
			err = rows.Scan(&id)
			if err != nil{
				continue
			}
			rargs = append(rargs,id,id)
		}
		if len(rargs) > 1{
			rconn.Send("ZADD",rargs...)
		}
	}

	list,err = redis.Ints(rconn.Do("ZRANGE",key,0,-1))
	if err != nil{
		log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
		return list
	}

	return list
}