package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
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
		defer rows.Close()
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,sort int
		for rows.Next(){
			err := rows.Scan(&id,&sort)
			if err != nil{
				log.Error("rows.Scan has error:%v",err)
				continue
			}
			rargs = append(rargs,sort,id)
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
	添加一个类别
 */
func AddCategory(name string)int{
	db := conn.GetMysqlConn()
	sql := "select count(*) from b_category where name=?"
	row := db.QueryRow(sql,name)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil{
		log.Error("AddCategory has error:%v",err)
		return -2
	}
	if cnt > 0 {
		return -3			//已有该类别
	}

	categroy_list := CategoryList()
	var sort = len(categroy_list)+1
	sql = "insert into b_category(name,sort) values(?,?)"
	result,err := db.Exec(sql,name,sort)
	if err != nil{
		log.Error("AddCategory has error:%v",err)
		return -4
	}
	id,err := result.LastInsertId()
	if err != nil{
		log.Error("AddCategory has error:%v",err)
		return -4
	}

	AddZsetData("categroyList",sort,id)
	
	return 0
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

	go DelKey("category:" + id)
	DelZsetData("categroyList",id)

	return 1
}

/**
	更新类别
 */
func UpCatetory(id,name string,sort int)int{
	db := conn.GetMysqlConn()
	sql := "update b_category set name=?,sort=? where id=?"
	_,err := db.Exec(sql,name,sort,id)
	if err != nil{
		log.Error("UpCatetory has error:%v",err)
		return -2
	}

	AddZsetData("categroyList",sort,id)
	args := make([]interface{},0)
	args = append(args,"category:" + id)
	args = append(args,"name")
	args = append(args,name)
	args = append(args,"sort")
	args = append(args,sort)
	HMset(args)

	return 0
}