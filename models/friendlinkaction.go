package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"fmt"
)

/**
	友链列表
 */
func FLink_List()[]int{
	list := make([]int,0)
	key := "flinkList"
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{
		sql := "select id,index from b_friendlink order by index asc"
		db := conn.GetMysqlConn()
		rows,err := db.Query(sql)
		if err != nil{
			return list
		}
		rargs := make([]interface{},0)
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
	cnt,_ := redis.Int(rconn.Do("ZCARD",key))
	if cnt > 0{
		list,err = redis.Ints(rconn.Do("ZRANGE",key,0,-1))
		if err != nil{
			log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
			return list
		}
	}
	return nil
}