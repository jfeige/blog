package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
)

//标记
type Tag struct {
	Id int `redis:"id"`
	Tag string `redis:"tag"`
}

const tag_field_cnt = 2

func (this *Tag) Load(id int)error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "tag:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		if len(values) == tag_field_cnt * 2{
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		}else{
			rconn.Do("DEL",key)
		}

	}
	sql := "select id,tag from b_tag where id=? limit 1"
	db := conn.GetMysqlConn()

	row := db.QueryRow(sql,id)
	err = row.Scan(&this.Id,&this.Tag)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}



/**
	多线程加载Article对象
 */
func MultipleLoadTag(id int,position int,tag_list []*Tag,wg *sync.WaitGroup){
	defer wg.Done()
	tag := new(Tag)
	err := tag.Load(id)
	if err == nil{
		tag_list[position] = tag
	}
	return
}


func FilterNilTag(tagList []*Tag)[]*Tag{
	//过滤空数据
	for k,v := range tagList{
		if v == nil && k  < len(tagList)-1{
			tagList = append(tagList[:k],tagList[k+1:]...)
		}else if k == len(tagList)-1 && v == nil{
			tagList = tagList[:len(tagList)-1]
		}
	}
	return tagList
}
