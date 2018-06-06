package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"strconv"
)

func TagList()[]int{
	list := make([]int,0)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "tagList"
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		sql := "select id from b_tag order by id asc"
		rows,err := db.Query(sql)
		if err != nil{
			log.Error("db.Query() has error:%v",err)
			return list
		}
		defer rows.Close()
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
		log.Error("redis.Ints has error:%v",err)
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

	_,err = db.Exec(sql,tagName)
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

	rows,err := db.Query(sql,id)
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
	defer tx.Rollback()
	_,err = tx.Exec(sql,id)
	if err != nil{
		log.Error("DelTag has error:%v",err)
		return -2
	}
	sql = "delete from b_actmapptags where t_id=?"

	_,err = tx.Exec(sql,id)
	if err != nil{
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



/**
	指定标签下文章总数
 */
func ArticleByTagCnt(tid int)int{
	var cnt int

	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "tagList:cnt:" + strconv.Itoa(tid)
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		sql := "select count(1) from b_actmapptags where t_id=? "

		row := db.QueryRow(sql,tid)
		err = row.Scan(&cnt)
		if err != nil{
			log.Error("row.Scan has error:%v",err)
			return	0
		}
		err = rconn.Send("set",key,cnt)
		if err != nil{
			log.Error("rconn.Send has error:%v",err)
			return	cnt
		}
	}

	cnt,err := redis.Int(rconn.Do("get",key))
	if err != nil{
		log.Error("redis.Int has error:%v",err)
		return 0
	}
	return cnt

}


/**
	文章列表
 */
func ArticleListByTag(args map[string]int)[]int{
	list := make([]int,0)
	pagesize,ok := args["pagesize"]
	if !ok{
		pagesize = 10
	}
	offset,ok := args["offset"]
	if !ok{
		offset = 0
	}
	tid := args["tid"]		//标签id

	key := "tagList:"  + strconv.Itoa(tid)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		sql := "SELECT a.id,a.publish_time FROM lifei.b_actmapptags m,lifei.b_article a where m.t_id=? and m.a_id=a.id order by a.id desc;"

		rows,err := db.Query(sql,tid)
		if err != nil{
			log.Error("stmt.Query has error:%v",err)
			return	list
		}
		defer rows.Close()
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,publish_time int
		for rows.Next(){
			err = rows.Scan(&id,&publish_time)
			if err != nil{
				log.Error("ArticleListByTagCnt has error:%v",err)
				continue
			}
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
		log.Error("redis.Ints has error:%v",err)
		return list
	}
	return list
}