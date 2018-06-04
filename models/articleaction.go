package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	log "github.com/alecthomas/log4go"
	"fmt"
	"time"
	"strings"
)


/**
	文章总数
 */
func ArticleCnt(cateid ...int)int{
	var cid,cnt int
	if len(cateid) > 0{
		cid = cateid[0]
	}
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "article_cnt:" + strconv.Itoa(cid)
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select count(1) from b_article "
		if cid > 0{
			sql = "select count(1) from b_article where categoryid=? "
			pargs = append(pargs,cid)
		}
		stmt,err := db.Prepare(sql)
		if err != nil{
			log.Error(fmt.Sprintf("db.Prepare has error:",err))
			return	0
		}
		defer stmt.Close()
		row := stmt.QueryRow(pargs...)
		err = row.Scan(&cnt)
		if err != nil{
			log.Error(fmt.Sprintf("row.Scan has error:",err))
			return	0
		}
		err = rconn.Send("set",key,cnt)
		if err != nil{
			log.Error(fmt.Sprintf("rconn.Send has error:",err))
			return	cnt
		}
	}

	cnt,err := redis.Int(rconn.Do("get",key))
	if err != nil{
		log.Error(fmt.Sprintf("redis.Int has error:%v",err))
		return 0
	}
	return cnt

}


/**
	文章列表
 */
func ArticleList(args map[string]int)[]int{
	list := make([]int,0)
	pagesize,ok := args["pagesize"]
	if !ok{
		pagesize = 10
	}
	isshow,ok := args["isshow"]
	if !ok{
		isshow = -1
	}
	offset,ok := args["offset"]
	if !ok{
		offset = 0
	}
	cateid,ok := args["cateid"]		//分类id
	if !ok || cateid < 0{
		cateid = 0
	}

	if !InArray(isshow,[]int{-1,0,1}){
		isshow = -1
	}

	key := "articleList:"  + strconv.Itoa(cateid)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		pargs := make([]interface{},0)
		sql := "select id,publish_time from b_article order by publish_time desc "
		if cateid > 0{
			sql = "select id,publish_time from b_article where categoryid=? order by publish_time desc "
			pargs = append(pargs,cateid)
		}
		stmt,err := db.Prepare(sql)
		if err != nil{
			log.Error(fmt.Sprintf("db.Prepare has error:",err))
			return	list
		}
		defer stmt.Close()

		rows,err := stmt.Query(pargs...)
		if err != nil{
			log.Error(fmt.Sprintf("stmt.Query has error:",err))
			return	list
		}
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,publish_time int
		for rows.Next(){
			rows.Scan(&id,&publish_time)
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
		log.Error(fmt.Sprintf("redis.Ints has error:%v",err))
		return list
	}
	return list
}


/**
	修改文章内容
 */
func UpdateArticleInfo(a_id int,title,content,tagids string)int{

	db := conn.GetMysqlConn()
	sql := "update b_article set title=?,content=? where id=?"
	tx,err := db.Begin()
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	defer tx.Rollback()
	stmt_a,err := tx.Prepare(sql)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	defer stmt_a.Close()
	_,err = stmt_a.Exec(title,content,a_id)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	//首先移除所有的标签，然后再插入
	sql = "delete from b_actmapptags where a_id=?";
	stmt_d,err := tx.Prepare(sql)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}
	defer stmt_d.Close()
	_,err = stmt_d.Exec(a_id)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
		return -2
	}

	tags := strings.Split(tagids,",")
	for _,id := range tags{
		sql = "insert into b_actmapptags(a_id,t_id) values(?,?)"
		stmt_t,err := tx.Prepare(sql)
		if err != nil{
			log.Error("UpdateArticleInfo has error:%v",err)
			return -2
		}
		defer stmt_t.Close()
		_,err = stmt_t.Exec(a_id,id)
		if err != nil{
			log.Error("UpdateArticleInfo has error:%v",err)
			return -2
		}
	}

	tx.Commit()

	//更新缓存
	rconn := conn.pool.Get()
	defer rconn.Close()

	keys := make([]interface{},0)
	keys = append(keys,"article:" + strconv.Itoa(a_id))
	keys = append(keys,"tagids:" + strconv.Itoa(a_id))

	_,err = rconn.Do("DEL",keys...)
	if err != nil{
		log.Error("UpdateArticleInfo has error:%v",err)
	}
	return 0
}

/**
	添加一篇文章
 */
func AddArticle(cateid int,title,content,tagids string) int{
	db := conn.GetMysqlConn()
	sql := "insert into b_article(title,content,categoryid,publish_time,publish_date) values(?,?,?,?,?)"

	tx,err := db.Begin()
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	defer tx.Rollback()
	stmt_a,err := tx.Prepare(sql)
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	defer stmt_a.Close()
	result,err := stmt_a.Exec(title,content,cateid,time.Now().Unix(),time.Now().Format("20060102"))
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	a_id,err := result.LastInsertId()
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}

	tags := strings.Split(tagids,",")
	for _,id := range tags{
		sql = "insert into b_actmapptags(a_id,t_id) values(?,?)"
		stmt_t,err := tx.Prepare(sql)
		if err != nil{
			log.Error("AddArticle has error:%v",err)
			return -2
		}
		defer stmt_t.Close()
		_,err = stmt_t.Exec(a_id,id)
		if err != nil{
			log.Error("AddArticle has error:%v",err)
			return -2
		}
	}


	sql = "update b_category set article_cnt=article_cnt+1 where id=?"
	stmt_c,err := tx.Prepare(sql)
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	defer stmt_c.Close()
	_,err = stmt_c.Exec(cateid)
	if err != nil{
		log.Error("AddArticle has error:%v",err)
		return -2
	}
	tx.Commit()
	//更新缓存
	rconn := conn.pool.Get()
	defer rconn.Close()

	keys := make([]interface{},0)
	keys = append(keys,"article_cnt:" + strconv.Itoa(cateid))
	keys = append(keys,"article_cnt:0")
	keys = append(keys,"articleList:"  + strconv.Itoa(cateid))
	keys = append(keys,"articleList:0")
	keys = append(keys,"category:" + strconv.Itoa(cateid))

	_,err = rconn.Do("DEL",keys...)
	if err != nil{
		log.Error("AddArticle has error:%v",err)
	}
	return 0
}


/**
	更新文章浏览数
 */
func UpdateReadCnt(a_id,cnt int){
	sql := "update b_article set read_count = read_count + ? where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		log.Error("UpdateReadCnt has error:%v",err)
		return
	}
	defer stmt.Close()
	_,err = stmt.Exec(cnt,a_id)
	if err != nil{
		log.Error("UpdateReadCnt has error:%v",err)
		return
	}

	key := "article:" + strconv.Itoa(a_id)

	rconn := conn.pool.Get()
	defer rconn.Close()

	rconn.Do("HINCRBY",key,"read_count",cnt)
}


/**
	删除文章
 */
func DelArticle(aid,cateid int)int{
	db := conn.GetMysqlConn()
	sql := "delArticle(?,?)"

	stmt,err := db.Prepare(sql)
	if err != nil{
		log.Error("DelArticle has error:%v",err)
		return -2
	}
	defer stmt.Close()
	row := stmt.QueryRow(aid,cateid)

	var errcode int

	err = row.Scan(&errcode)
	if err != nil{
		log.Error("DelArticle has error:%v",err)
		return -2
	}
	if errcode < 0{
		return errcode
	}
	//key := "article_cnt:" + strconv.Itoa(cid)
	//key := "articleList:"  + strconv.Itoa(cateid)
	//key := "category:" + strconv.Itoa(id)
	//key := "commentcnt:" + strconv.Itoa(arteid)
	//key := "tagids:" + strconv.Itoa(this.Id)
	rconn := conn.pool.Get()
	defer rconn.Close()

	keys := make([]interface{},0)
	keys = append(keys,"article:" + strconv.Itoa(aid))
	keys = append(keys,"article_cnt:" + strconv.Itoa(cateid))
	keys = append(keys,"article_cnt:0")
	keys = append(keys,"articleList:" + strconv.Itoa(aid))
	keys = append(keys,"articleList:0")
	keys = append(keys,"category:" + strconv.Itoa(cateid))
	keys = append(keys,"tagids:" + strconv.Itoa(aid))

	_,err = rconn.Do("DEL",keys...)
	if err != nil{
		log.Error("DelArticle has error:%v",err)
	}
	return 0
}