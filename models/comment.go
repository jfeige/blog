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
	Cid int `redis:"cid"`
	Type int `redis:"type"`
	Content string  `redis:"content"`
	Atime int64 `redis:"atime"`
	Ordertime float64 `redis:"ordertime"`
}

const comment_field_cnt = 8

/**
	加载指定的评论
 */

func (this *Comment) Load(id int) error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "comment:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		if len(values) == comment_field_cnt * 2{
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		}else{
			rconn.Do("DEL",key)
		}

	}
	sql := "select id,articleid,name,cid,type,content,atime,ordertime from b_comment where id=? limit 1"
	db := conn.GetMysqlConn()

	row := db.QueryRow(sql,id)
	err = row.Scan(&this.Id,&this.Articleid,&this.Name,&this.Cid,&this.Type,&this.Content,&this.Atime,&this.Ordertime)
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

func FilterNilComment(commentList []*Comment)[]*Comment{
	//过滤空数据
	for k,v := range commentList{
		if v == nil && k  < len(commentList)-1{
			commentList = append(commentList[:k],commentList[k+1:]...)
		}else if k == len(commentList)-1 && v == nil{
			commentList = commentList[:len(commentList)-1]
		}
	}
	return commentList
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


/**
	前台摘要显示
 */
func (this *Comment) FormatContent()string{
	content := []rune(this.Content)
	if len(content) > 25{
		return string(content[:25]) + " ..."
	}
	return this.Content
}