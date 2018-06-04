package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
)

type FriendLink struct {
	Id int `redis:"id"`
	Webname string `redis:"webname"`
	Weburl string `redis:"weburl"`
	Sort 	int `redis:"sort"`
}


const friendlink_field_cnt = 4

func (this *FriendLink) Load(id int) error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "flink:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		if len(values) == friendlink_field_cnt * 2 {
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		}else{
			rconn.Do("DEL",key)
		}
	}
	sql := "select id,webname,weburl,sort from b_friendlink where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Webname,&this.Weburl,&this.Sort)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}



/**
	多线程加载Category对象
 */
func MultipleLoadFLink(id int,position int,flink_list []*FriendLink,wg *sync.WaitGroup){
	defer wg.Done()
	flink := new(FriendLink)
	err := flink.Load(id)
	if err == nil{
		flink_list[position] = flink
	}
	return
}