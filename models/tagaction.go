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

/**
	添加新标签
 */
func AddTag(tagName string)(errcode int){
	//判断name是否重复
	sql := "select count(*) from b_tag where tag=?"
	db := conn.GetMysqlConn()
	row := db.QueryRow(sql,tagName)
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
	_,err = stmt.Exec(tagName)
	if err != nil{
		log.Error("AddTag has error:%v",err)
		errcode = -3
		return
	}

	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "taglist"

	_,err = rconn.Do("DEL",key)
	if err != nil{
		log.Error("AddTag has error:%v",err)
	}

	return
}


/**
	删除标签
 */
func DelTag(id string)int {
	var aids = make([]string,0)
	//获取有哪些文章引用了该标签
	sql := "select a_id from b_actmapptags where t_id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		log.Error("DelTag has error:%v",err)
		return -2
	}
	defer stmt.Close()
	rows,err := stmt.Query(id)
	if err != nil{
		log.Error("DelTag has error:%v",err)
		return -2
	}
	var aid string
	for rows.Next(){
		err = rows.Scan(&aid)
		if err != nil{
			continue
		}
		aids = append(aids,aid)
	}
	//以下删除表中引用和记录
	sql = "delete from b_tag where id =?"
	tx,err := db.Begin()
	if err != nil{
		log.Error("DelTag has error:%v",err)
		return -2
	}
	stmt_tag,err := tx.Prepare(sql)
	if err != nil{
		log.Error("DelTag has error:%v",err)
		return -2
	}
	_,err = stmt_tag.Exec(id)
	if err != nil{
		log.Error("DelTag has error:%v",err)
		return -2
	}
	stmt_article,err := tx.Prepare(sql)
	if err != nil{
		tx.Rollback()
		log.Error("DelTag has error:%v",err)
		return -2
	}
	_,err = stmt_article.Exec(id)
	if err != nil{
		tx.Rollback()
		log.Error("DelTag has error:%v",err)
		return -2
	}

	tx.Commit()

	rconn := conn.pool.Get()
	defer rconn.Close()

	keys := make([]interface{},0)
	keys = append(keys,"taglist")
	keys = append(keys,"tag:" + id)
	for _,v := range aids{
		keys = append(keys,"tagids:" + v)
	}

	rconn.Do("DEL",keys...)

	return 0

}