package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
)
//类别
type Category struct {
	Id int `redis:"id"`
	Name string `redis:"name"`
	Article_cnt int `redis:"article_cnt"`
	Sort string `redis:"sort"`
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
	sql := "select id,name,article_cnt,sort from b_category where id=?"
	db := conn.GetMysqlConn()
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&this.Id,&this.Name,&this.Article_cnt,&this.Sort)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}


/**
	多线程加载Category对象
 */
func MultipleLoadCategory(id int,position int,category_list []*Category,wg *sync.WaitGroup){
	defer wg.Done()
	category := new(Category)
	err := category.Load(id)
	if err == nil{
		category_list[position] = category
	}
	return
}