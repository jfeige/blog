package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
)
//文章
type Article struct {
	Id int `redis:"id"`
	Title string `redis:"title"`
	Content string `redis:"content"`
	Userid int  `redis:"userid"`
	Categoryid int `redis:"categoryid"`
	Tagid int `redis:"tagid"`
	Read_count int `redis:"read_count"`
	Comment_count int `redis:"comment_count"`
	Publish_time int `redis:"publish_time"`
	Publish_date int `redis:"publish_date"`
	Isshow int `redis:"isshow"`
}

/**
	加载指定的文章
 */

func (this *Article) Load(id int) error{
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
	sql := "select id,title,content,userid,categoryid,tagid,read_count,comment_count,publish_time,publish_date,isshow from b_article where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Title,&this.Content,&this.Userid,&this.Categoryid,&this.Tagid,&this.Read_count,&this.Comment_count,&this.Publish_time,&this.Publish_date,&this.Isshow)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}

/**
	多线程加载Article对象
 */
func MultipleLoadArticle(id int,position int,article_list []*Article,wg *sync.WaitGroup){
	defer wg.Done()
	article := new(Article)
	err := article.Load(id)
	if err == nil{
		article_list[position] = article
	}
	return
}