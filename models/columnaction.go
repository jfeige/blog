package models

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
)

func ColumnList()[]int{
	list := make([]int,0)
	rconn := conn.GetRedisConn()
	defer rconn.Close()
	key := "columnList"
	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if !exists{
		db := conn.GetMysqlConn()
		sql := "select id,sort from b_column order by sort asc"
		rows,err := db.Query(sql)
		if err != nil{
			log.Error("db.Query() has error:%v",err)
			return list
		}
		defer rows.Close()
		rargs := make([]interface{},0)
		rargs = append(rargs,key)
		var id,sort int
		for rows.Next(){
			err = rows.Scan(&id,&sort)
			if err != nil{
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
	添加一个栏目
 */
func AddColumn(title,url string,tp int)int{
	sql := "select count(1) from b_column where title=? or url =?"
	db := conn.GetMysqlConn()
	row := db.QueryRow(sql,title,url)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil{
		log.Error("AddColumn has error:%v")
		return -1
	}
	if cnt > 0{
		return -2;	//已有相同标题或者url的栏目，不能继续添加
	}

	columnList := ColumnList()
	var sort = len(columnList)+1

	sql = "insert into b_column(title,url,sort,tp) values(?,?,?,?)"
	result,err := db.Exec(sql,title,url,sort,tp)
	if err != nil{
		return -1
	}
	id,err := result.LastInsertId()
	if err != nil{
		log.Error("AddColumn has error:%v")
		return -1
	}

	go AddZsetData("columnList",sort,id)

	return 1
}

/**
	更新栏目
 */
func UpColumn(id,title,url string,sort int)int{
	db := conn.GetMysqlConn()
	sql := "update b_column set title=?,url=?,sort=? where id=?"
	_,err := db.Exec(sql,title,url,sort,id)
	if err != nil{
		log.Error("UpCatetory has error:%v",err)
		return -2
	}

	AddZsetData("columnList",sort,id)
	args := make([]interface{},0)
	args = append(args,"column:" + id)
	args = append(args,"title")
	args = append(args,title)
	args = append(args,"sort")
	args = append(args,sort)
	HMset(args)

	return 0
}

/**
	移除一个栏目
 */
func DelColomn(id int)int{
	sql := "delete from b_column where id=?"
	db := conn.GetMysqlConn()

	_,err := db.Exec(sql,id)

	if err != nil{
		log.Error("RemoveColomn has error:%v,id:%d",err,id)
		return -1
	}

	go DelZsetData("columnList",id)

	return 1
}