package models

import
(
	"github.com/garyburd/redigo/redis"
	"fmt"
)
type Webset struct {
	Id int `redis:"id"`
	Sitename string `redis:"sitename"`
	Sitedesc string `redis:"sitedesc"`
	Siteurl  string `redis:"siteurl"`
	Keywords string `redis:"keywords"`
	Descri string `redis:"descri"`
	Name string `redis:"name"`
	Phone string `redis:"phone"`
	Qq string `redis:"qq"`
	Email string `redis:"email"`
	Place string `redis:"place"`
	Github string `redis:"github"`
}

const webset_fieldd_count = 12

func (this *Webset) Load()error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "webset"
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		if len(values) == webset_fieldd_count * 2{
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		}else{
			rconn.Do("DEL",key)
		}
	}
	sql := "select id,sitename,sitedesc,siteurl,keywords,descri,name,phone,qq,email,place,github from b_webset order by id desc limit 1"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	err = row.Scan(&this.Id,&this.Sitename,&this.Sitedesc,&this.Siteurl,&this.Keywords,&this.Descri,&this.Name,&this.Phone,&this.Qq,&this.Email,&this.Place,&this.Github)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}



func (this *Webset)UpdateWebSet() error{
	db := conn.GetMysqlConn()
	sql := "update b_webset set sitename=?,sitedesc=?,siteurl=?,keywords=?,descri=?,name=?,phone=?,qq=?,email=?,place=?,github=? where id=?"
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	fmt.Println(this)
	defer stmt.Close()
	_,err = stmt.Exec(this.Sitename,this.Sitedesc,this.Siteurl,this.Keywords,this.Descri,this.Name,this.Phone,this.Qq,this.Email,this.Place,this.Github,this.Id)
	if err != nil{
		return err
	}

	rconn := conn.pool.Get()
	defer rconn.Close()


	key := "webset"
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)

	return nil
}