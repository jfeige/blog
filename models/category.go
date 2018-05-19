package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
)
//类别
type Category struct {
	Id int `redis:"id"`
	Name string `redis:"name"`
	Index string `redis:"index"`
}


func (this *Category) Load(id int)error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "category:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		err = redis.ScanStruct(values, this)
		if err == nil {
			return nil
		}
	}
	sql := "select id,name,index from b_category where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Name,&this.Index)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}