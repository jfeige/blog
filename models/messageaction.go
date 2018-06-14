package models

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

/**
留言总数
*/
func MsgCnt() int {
	var cnt int

	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "msgList:cnt"
	exists, _ := redis.Bool(rconn.Do("EXISTS", key))
	if !exists {
		db := conn.GetMysqlConn()
		sql := "select count(1) from b_messageboard "

		row := db.QueryRow(sql)
		err = row.Scan(&cnt)
		if err != nil {
			log.Error(fmt.Sprintf("row.Scan has error:", err))
			return 0
		}
		err = rconn.Send("set", key, cnt)
		if err != nil {
			log.Error(fmt.Sprintf("rconn.Send has error:", err))
			return cnt
		}
	}

	cnt, err := redis.Int(rconn.Do("get", key))
	if err != nil {
		log.Error(fmt.Sprintf("redis.Int has error:%v", err))
		return 0
	}
	return cnt

}

/**
留言列表
*/
func MsgList(args map[string]int) []int {
	list := make([]int, 0)
	pagesize, ok := args["pagesize"]
	if !ok {
		pagesize = 10
	}
	offset, ok := args["offset"]
	if !ok {
		offset = 0
	}

	key := "msgList:"
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	exists, _ := redis.Bool(rconn.Do("EXISTS", key))
	if !exists {
		db := conn.GetMysqlConn()
		sql := "select id,atime from b_messageboard order by atime desc "

		rows, err := db.Query(sql)
		if err != nil {
			log.Error("MsgList has error:%v", err)
			return list
		}
		defer rows.Close()
		rargs := make([]interface{}, 0)
		rargs = append(rargs, key)
		var id, atime int
		for rows.Next() {
			err = rows.Scan(&id, &atime)
			if err != nil {
				log.Error("MsgList has error:%v", err)
				continue
			}
			rargs = append(rargs, atime, id)

		}
		if len(rargs) > 1 {
			rconn.Send("ZADD", rargs...)
		}
	}

	limit := offset + pagesize - 1
	var params = []interface{}{key, offset, limit}
	list, err = redis.Ints(rconn.Do("ZRANGE", params...))
	if err != nil {
		log.Error(fmt.Sprintf("MsgList has error:%v", err))
	}
	return list
}

/**
添加一条留言
*/
func AddMsg(name, content string) int {
	sql := "insert into b_messageboard(user,content,atime,mdate) values(?,?,?,?)"
	db := conn.GetMysqlConn()

	atime := time.Now().Unix()
	mdate := time.Now().Format("20060102")
	result, err := db.Exec(sql, name, content, atime, mdate)
	if err != nil {
		log.Error("AddMsg has error. name:%v,content:%s,error:%v", name, content, err)
		return -1
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Error("AddMsg has error. name:%v,content:%s,error:%v", name, content, err)
		return -1
	}

	AddZsetData("msgList:", atime, id)
	Incr("msgList:cnt")

	return 1

}

/**
删除一条评论
*/
func DelMessage(mid int) int {
	db := conn.GetMysqlConn()
	sql := "delete from b_messageboard where id=?"

	var errcode int
	_, err := db.Exec(sql, mid)
	if err != nil {
		return -2
	}

	DelZsetData("msgList:*", mid)
	Decr("msgList:cnt")

	go DelKey("message:" + strconv.Itoa(mid))
	return errcode
}

/**
删除多条评论
*/
func DelMultiMessage(ids string, mids []string) int {

	db := conn.GetMysqlConn()
	sql := "delete from b_messageboard where id in (" + ids + ")"

	var errcode int
	_, err := db.Exec(sql)
	if err != nil {
		log.Error("DelMultiMessage has error:%v", err)
		return -2
	}

	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "msgList:*"
	BlurDelKeys(key)

	keys := make([]interface{}, 0)
	for _, mid := range mids {
		keys = append(keys, "message:"+mid)
	}
	go DelKeys(keys)

	return errcode
}
