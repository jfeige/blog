package models

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/garyburd/redigo/redis"
)

/**
友链列表
*/
func FLink_List() []int {
	list := make([]int, 0)
	key := "flinkList"
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists, _ := redis.Bool(rconn.Do("EXISTS", key))

	if !exists {
		sql := "select id,sort from b_friendlink order by sort asc"
		db := conn.GetMysqlConn()
		rows, err := db.Query(sql)
		if err != nil {
			fmt.Println(err)
			return list
		}
		defer rows.Close()
		rargs := make([]interface{}, 0)
		rargs = append(rargs, key)
		var id, sort int
		for rows.Next() {
			err := rows.Scan(&id, &sort)
			if err != nil {
				log.Error("rows.Scan has error:%v", err)
				continue
			}
			rargs = append(rargs, sort, id)
		}
		if len(rargs) > 1 {
			rconn.Send("ZADD", rargs...)
		}
	}
	list, err = redis.Ints(rconn.Do("ZRANGE", key, 0, -1))
	if err != nil {
		log.Error("redis.Ints has error:%v", err)
		return list
	}
	return list
}

/**
添加一个友链
*/
func AddFlink(webname, weburl string) int {
	db := conn.GetMysqlConn()
	sql := "select count(*) from b_friendlink where webname=?"
	row := db.QueryRow(sql, webname)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		log.Error("AddFlink has error:%v", err)
		return -2
	}
	if cnt > 0 {
		return -3 //已有该类别
	}

	categroy_list := CategoryList()
	var sort = len(categroy_list) + 1
	sql = "insert into b_friendlink(webname,weburl,sort) values(?,?,?)"
	result, err := db.Exec(sql, webname, weburl, sort)
	if err != nil {
		log.Error("AddFlink has error:%v", err)
		return -4
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Error("AddFlink has error:%v", err)
		return -4
	}

	AddZsetData("flinkList", sort, id)

	return 0
}

/**
删除一个友链
*/
func DelFlink(id string) int {
	db := conn.GetMysqlConn()
	sql := "delete from b_friendlink where id=?"
	_, err := db.Exec(sql, id)
	if err != nil {
		log.Error("DelFlink has error:%v", err)
		return -2
	}

	go DelKey("flink:" + id)
	DelZsetData("flinkList", id)

	return 1
}

/**
更新类别
*/
func UpFlink(id, webname, weburl string, sort int) int {
	db := conn.GetMysqlConn()
	sql := "update b_friendlink set webname=?,weburl=?,sort=? where id=?"
	_, err := db.Exec(sql, webname, weburl, sort, id)
	if err != nil {
		log.Error("UpFlink has error:%v", err)
		return -2
	}

	AddZsetData("flinkList", sort, id)
	args := make([]interface{}, 0)
	args = append(args, "flink:"+id)
	args = append(args, "webname")
	args = append(args, webname)
	args = append(args, "weburl")
	args = append(args, weburl)
	args = append(args, "sort")
	args = append(args, sort)
	HMset(args)

	return 0
}
