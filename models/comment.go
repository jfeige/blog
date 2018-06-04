package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"sync"
	"time"
)
//评论
type Comment struct {
	Id int `redis:"id"`
	Articleid int `redis:"articleid"`
	Name string `redis:"name"`
	Content string  `redis:"content"`
	Atime int64 `redis:"atime"`
}

/**
	加载指定的评论
 */

func (this *Comment) Load(id int) error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "comment:" + strconv.Itoa(id)
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
		log.Error("db.Prepare() has error:%v",err)
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Articleid,&this.Name,&this.Content,&this.Atime)
	if err != nil{
		log.Error("row.Scan() has error:%v",err)
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



func (this *Comment) FormatAtime(format string)string{

	return time.Unix(int64(this.Atime),0).Format(format)

}


/**
	返回评论所在文章信息
 */
func (this *Comment) ArticleInfo()*Article{
	article := new(Article)
	article.Load(this.Articleid)
	return article
}