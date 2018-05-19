package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
)

//标记
type Tag struct {
	Id int `redis:"id"`
	Tag string `redis:"tag"`
}



func (this *Tag) Load(id int)error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "tag:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		err = redis.ScanStruct(values, this)
		if err == nil {
			return nil
		}
	}
	sql := "select id,tag from b_tag where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Tag)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}