package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
)
type User struct {
	Id int `redis:"id"`
	Name string `redis:"name"`
	Passwd string `redis:"passwd"`
	Nickname string `redis:"nickname"`
}


func (this *User) Load(id int)error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "user:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		err = redis.ScanStruct(values, this)
		if err == nil {
			return nil
		}
	}
	sql := "select id,name,passwd,nickname from b_user where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Name,&this.Passwd,&this.Nickname)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}