package models

import
(
	"github.com/garyburd/redigo/redis"
)
type Webset struct {
	Id int `redis:"id"`
	Banner string `redis:"banner"`
	Sitename string `redis:"sitename"`
	Sitedesc string `redis:"sitedesc"`
	Name string `redis:"name"`
	Nickname string `redis:"nickname"`
	Weburl string `redis:"weburl"`
	Place string `redis:"place"`
	Vocation string `redis:"vocation"`
}


func (this *Webset) Load()error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "webset"
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		err = redis.ScanStruct(values, this)
		if err == nil {
			return nil
		}
	}
	sql := "select id,banner,sitename,sitedesc,name,nickname,weburl,place,vocation from b_webset order by id desc limit 1"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	err = row.Scan(&this.Id,&this.Banner,&this.Sitename,&this.Sitedesc,&this.Name,&this.Nickname,&this.Weburl,&this.Place,&this.Vocation)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}