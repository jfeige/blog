package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
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

/**
	添加新标签
 */
func AddTag(context *gin.Context){
	var errcode int
	defer context.JSON(http.StatusOK,gin.H{
		"errcode":errcode,
	})
	name,_ := context.GetPostForm("name")

	if name == ""{
		//判断name是否为空
		errcode = -1
		return
	}
	//判断name是否重复
	sql := "select count(*) from b_tag where tag=?"
	db := conn.GetMysqlConn()
	row := db.QueryRow(sql)
	var cnt int
	row.Scan(&cnt)
	if cnt > 0{
		errcode = -2
		return
	}
	sql = "insert into b_tag(tag) values(?)"
	stmt,err := db.Prepare(sql)
	if err != nil{
		log.Error("AddTag has error:%v",err)
		errcode = -3
		return
	}
	defer stmt.Close()
	_,err = stmt.Exec(name)
	if err != nil{
		log.Error("AddTag has error:%v",err)
		errcode = -4
		return
	}

	return
}