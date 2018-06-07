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
		sql := "select id,ordertime from b_comment where articleid=? order by atime asc"
		db := conn.GetMysqlConn()
		rows,err := db.Query(sql,articleid)
		if err != nil{
			log.Error("CommentList has error:%v",err)
			return list
		}
		defer rows.Close()
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,ordertime,score float64
		for rows.Next(){
			err := rows.Scan(&id,&ordertime)
			if err != nil{
				log.Error(fmt.Sprintf("rows.Scan has error:%v",err))
				continue
			}
			score = ordertime * 1000000 + id
			rargs = append(rargs,score,id)
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
	评论总数量
 */
func ManageCommentCnt(articleid ...int)int{
	var arteid int
	if len(articleid) > 0{
		arteid = articleid[0]
	}

	key := "commentList:cnt:" + strconv.Itoa(arteid) + "|-1"
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
	pagesize,ok := args["pagesize"]
	if !ok{
		pagesize = 10
	}
	offset,ok := args["offset"]
	if !ok{
		offset = 0
	}
	tp,ok := args["type"]
	if !ok{
		tp = -1
	}

	list := make([]int,0)
	key := "commentList:" + strconv.Itoa(arteid) + "|" + strconv.Itoa(tp)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	fmt.Println(key)
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))

	if !exists{
		pargs := make([]interface{},0)
		sql := "select id,atime from b_comment order by atime desc"
		if tp != -1{
			if arteid > 0{
				sql = "select id,atime from b_comment where articleid=? and type=? order by atime desc"
				pargs = append(pargs,arteid,tp)
			}
		}else{
			if arteid > 0{
				sql = "select id,atime from b_comment where articleid=? order by atime desc"
				pargs = append(pargs,arteid)
			}
		}

		db := conn.GetMysqlConn()
		var err error
		rows,err := db.Query(sql,pargs...)

		if err != nil{
			log.Error("ManageCommentList has error:%v",err)
			return list
		}
		defer rows.Close()
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
	添加一条评论或回复
	tp 0:评论 1:回复
	cid :回复的评论id
 */

func AddComment(aid,tp,cid int,name interface{},content string){
	sql := "insert into b_comment(articleid,type,cid,name,content,atime) values(?,?,?,?,?,?)"
	db := conn.GetMysqlConn()

	tx,err := db.Begin()
	if err != nil{
		log.Error("AddComment has error. aid:%d,tp:%d,cid:%d,name:%v,content;%s,error:%v",aid,tp,cid,name,content,err)
		return
	}
	defer tx.Rollback()

	result,err := tx.Exec(sql,aid,tp,cid,name,content,time.Now().Unix())
	if err != nil{
		log.Error("AddComment has error. aid:%d,tp:%d,cid:%d,name:%v,content;%s,error:%v",aid,tp,cid,name,content,err)
		return
	}
	commentId,_ := result.LastInsertId()
	var ordertime float64
	if tp == 0{
		ordertime = float64(commentId)
	}else{
		ordertime = float64(cid) + 0.1
	}
	sql = "update b_comment set ordertime=?,cid=? where id=?"
	tx.Exec(sql,ordertime,commentId,commentId)

	if tp == 0{
		sql = "update b_article set comment_count=comment_count+1 where id=?"
		tx.Exec(sql,aid)
	}


	tx.Commit()

	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "commentList:*"
	DelKeys(key)

	keys := make([]interface{},0)
	keys = append(keys,"article:" + strconv.Itoa(aid))
	rconn.Do("DEL",keys...)

}


/**
	删除一条评论
 */
func DelComment(aid,cid int)int{
	db := conn.GetMysqlConn()
	sql := "call delComment(?,?)"


	var errcode int
	row := db.QueryRow(sql,aid,cid)
	err = row.Scan(&errcode)
	if err != nil{
		return -2
	}

	rconn := conn.pool.Get()
	defer rconn.Close()

	key := "commentList:*"
	DelKeys(key)

	keys := make([]interface{},0)
	keys = append(keys,"comment:" + strconv.Itoa(cid))
	keys = append(keys,"article:" + strconv.Itoa(aid))

	_,err = rconn.Do("DEL",keys...)
	if err != nil{
		log.Error("DelComment has error:%v",err)
	}
	return errcode
}