package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
	log "github.com/alecthomas/log4go"
	"strings"
	"time"
)
//文章
type Article struct {
	Id int `redis:"id"`
	Title string `redis:"title"`
	Content string `redis:"content"`
	Userid int  `redis:"userid"`
	Categoryid int `redis:"categoryid"`
	Read_count int `redis:"read_count"`
	Comment_count int `redis:"comment_count"`
	Publish_time int64 `redis:"publish_time"`
	Publish_date int64 `redis:"publish_date"`
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
	sql := "select id,title,content,userid,categoryid,read_count,comment_count,publish_time,publish_date,isshow from b_article where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Title,&this.Content,&this.Userid,&this.Categoryid,&this.Read_count,&this.Comment_count,&this.Publish_time,&this.Publish_date,&this.Isshow)
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
/**
	获取作者
 */
func (this *Article) User()string{
	return "张三"
}

/**
	获取所属类别
 */
func (this *Article) Category()string{
	catetory := new(Category)
	err := catetory.Load(this.Categoryid)
	if err != nil{
		log.Error("catetory.Load() has error:%v",err)
		return ""
	}
	return catetory.Name
}

/**
	得到该文章对应的标签id
 */
func (this *Article) GetTagIds() string{
	var ids string
	key := "tagids:" + strconv.Itoa(this.Id)
	rconn := conn.pool.Get()
	defer rconn.Close()

	tmpIds,err:= redis.String(rconn.Do("GET",key))
	if tmpIds != "" && err == nil{
		return tmpIds
	}
	sql := "select t_id from b_actmapptags where a_id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		log.Error("GetTagIds has error:%v",err)
		return ids
	}
	defer stmt.Close()
	rows,err := stmt.Query(this.Id)
	if err != nil{
		log.Error("GetTagIds has error:%v",err)
		return ids
	}
	tmp_Ids := make([]string,0)
	var tid string
	for rows.Next(){
		err = rows.Scan(&tid)
		if err != nil{
			continue
		}
		tmp_Ids = append(tmp_Ids,tid)
	}
	ids = strings.Join(tmp_Ids,",")
	rconn.Do("SET",key,ids)
	return ids
}

/**
	解析所属标签
 */
func (this *Article) Tag()map[int]string{
	tags := make(map[int]string)

	tagId := this.GetTagIds()
	if tagId != ""{
		ids := strings.Split(tagId,",")
		for _,tmpId := range ids{
			id,_ := strconv.Atoi(tmpId)
			tag := new(Tag)
			err := tag.Load(id)
			if err != nil{
				log.Error("tag.Load() has error:%v",err)
				continue
			}
			tags[id] = tag.Tag
		}
	}

	return tags
}

/**
	判断指定的tagID是否是
 */
func (this *Article) IsTag(tid int)bool{
	tags := this.Tag()

	_,ok := tags[tid]

	return ok
}

/**
	格式化日期和时间
 */
func (this *Article) PublishTime(flags ...int)string{
	var flag = 0		//0:publish_time;1:publish_date
	var ftime string
	if len(flags) > 0{
		flag = flags[0]
	}
	if flag == 0{
		ftime = time.Unix(this.Publish_time,0).Format("2006-01-02 15:04:05")
	}else{
		ftime = time.Unix(this.Publish_time,0).Format("2006-01-02")
	}

	return ftime
}

/**
	按照给定的格式格式化日期和时间
 */
func (this *Article) FormatPublishTime(format string)string{

	return time.Unix(this.Publish_time,0).Format(format)
}


/**
	前台摘要显示
 */
func (this *Article) FormatContent()string{
	if len(this.Content) > 500{
		return this.Content[:500]
	}
	return this.Content
}