package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"fmt"
)

/**
	分类列表
 */
func CategoryList()[]int{
	list := make([]int,0)
	key := "categroyList"
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{
		sql := "select id,sort from b_category order by sort asc"
		db := conn.GetMysqlConn()
		rows,err := db.Query(sql)
		if err != nil{
			log.Error("db.Query has error:%v",err)
			return list
		}
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,index int
		for rows.Next(){
			err := rows.Scan(&id,&index)
			if err != nil{
				log.Error(fmt.Sprintf("rows.Scan has error:%v",err))
				continue
			}
			rargs = append(rargs,index,id)
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


/**
	删除一个类别
 */
func DelCatetory(id string)int{
	db := conn.GetMysqlConn()
	sql := "delete from b_category where id=?"
	_,err := db.Exec(sql,id)
	if err != nil{
		log.Error("DelCategory has error:%v",err)
		return -2
	}

	key := "category:" + id

	rconn := conn.pool.Get()
	defer rconn.Close()

	_,err = rconn.Do("DEL",key)
	if err != nil{
		log.Error("DelCategory has error:%v",err)
	}
	return 0
}