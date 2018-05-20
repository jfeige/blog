package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
)

type FriendLink struct {
	Id int `redis:"id"`
	Webname string `redis:"webname"`
	Weburl string `redis:"weburl"`
	Index 	int `redis:"index"`
}

func (this *FriendLink) Load(id int) error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "flink:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		err = redis.ScanStruct(values, this)
		if err == nil {
			return nil
		}
	}
	sql := "select id,webname,weburl,index from b_friendlink where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Webname,&this.Weburl,&this.Index)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}
