package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
)
//评论
type Comment struct {
	Id int `redis:"id"`
	Articleid string `redis:"articleid"`
	Name string `redis:"name"`
	Content int  `redis:"content"`
	Atime int `redis:"atime"`
}

/**
	加载指定的评论
 */

func (this *Comment) Load(id int) error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "article:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		err = redis.ScanStruct(values, this)
		if err == nil {
			return nil
		}
	}
	sql := "select id,articleid,name,content,atime from b_comment where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Articleid,&this.Name,&this.Content,&this.Atime)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}

/**
	多线程加载Comment对象
 */
func MultipleLoadComment(id int,position int,comment_list []*Comment,wg *sync.WaitGroup){
	defer wg.Done()
	comment := new(Comment)
	err := comment.Load(id)
	if err == nil{
		comment_list[position] = comment
	}
	return
}