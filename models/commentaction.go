package models

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"time"
	"strconv"
)

/**
	根据文章id获取评论id列表
 */
func CommentList(articleid int)[]int{
	list := make([]int,0)
	key := "commentList:" + strconv.Itoa(articleid)
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{
		sql := "select id,atime from b_comment where articleid=? order by atime desc"
		db := conn.GetMysqlConn()
		rows,err := db.Query(sql,articleid)
		if err != nil{
			log.Error("CommentList has error:%v",err)
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
	list,err = redis.Ints(rconn.Do("ZREVRANGE",key,0,-1))
	if err != nil{
		log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
		return list
	}
	return list
}

/**
	评论总数量
 */
func ManageCommentCnt(articleid ...int)int{
	var arteid int
	if len(articleid) > 0{
		arteid = articleid[0]
	}

	key := "commentcnt:" + strconv.Itoa(arteid)
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{

		sql := "select count(1) from b_comment"
		if arteid > 0{
			sql = "select count(1) from b_comment where articleid=?"
		}
		db := conn.GetMysqlConn()
		var cnt int
		if arteid > 0{
			row := db.QueryRow(sql,arteid)
			err = row.Scan(&cnt)
		}else{
			row := db.QueryRow(sql)
			err = row.Scan(&cnt)
		}

		if err != nil{
			log.Error("ManageCommentCnt has error:%v",err)
		}else{
			rconn.Do("SET",key,cnt)
		}

		return cnt

	}
	cnt,err := redis.Int(rconn.Do("GET",key))
	if err != nil{
		log.Error("ManageCommentCnt has error:%v",err)
	}
	return cnt
}

/**
	后台读取评论列表
 */
func ManageCommentList(args map[string]int)[]int{

	arteid,ok := args["arteid"]
	if !ok{
		arteid = 0
	}
	pagesize := args["pagesize"]
	offset := args["offset"]

	list := make([]int,0)
	key := "commentList:" + strconv.Itoa(arteid)
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{
		pargs := make([]interface{},0)
		sql := "select id,atime from b_comment order by atime desc"
		if arteid > 0{
			sql = "select id,atime from b_comment where articleid=? order by atime desc"
			pargs = append(pargs,arteid)
		}
		db := conn.GetMysqlConn()
		var err error
		rows,err := db.Query(sql,pargs...)

		if err != nil{
			log.Error("ManageCommentList has error:%v",err)
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
	limit := offset + pagesize - 1
	list,err = redis.Ints(rconn.Do("ZREVRANGE",key,offset,limit))
	if err != nil{
		log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
		return list
	}
	return list
}


/**
	添加一条回复
 */

func AddComment(aid int,name interface{},content string){
	sql := "insert into b_comment(articleid,name,content,atime) values(?,?,?,?)"
	db := conn.GetMysqlConn()

	tx,err := db.Begin()
	if err != nil{
		log.Error("AddComment has error.aid:%d,name:%v,content:%s,error:%v",aid,name,content,err)
		return
	}
	stmt,err := tx.Prepare(sql)
	if err != nil{
		log.Error("AddComment has error. aid:%d,name:%v,content;%s,error:%v",aid,name,content,err)
		return
	}
	stmt.Exec(aid,name,content,time.Now().Unix())

	sql = "update b_article set comment_count=comment_count+1 where id=?"
	up_stmt,err := tx.Prepare(sql)
	if err != nil{
		log.Error("AddComment has error. aid:%d,name:%v,content;%s,error:%v",aid,name,content,err)
		tx.Rollback()
		return
	}
	up_stmt.Exec(aid)

	tx.Commit()

	rconn := conn.pool.Get()
	defer rconn.Close()

	keys := make([]interface{},0)
	keys = append(keys,"commentList:" + strconv.Itoa(aid))
	keys = append(keys,"commentList:0")
	keys = append(keys,"article:" + strconv.Itoa(aid))
	rconn.Do("DEL",keys...)

}