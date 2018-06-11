package models

import (
	"strconv"
	"github.com/garyburd/redigo/redis"
	"sync"
)
type Column struct {
	Id int `redis:"id"`
	Title string `redis:"title"`
	Url string `redis:"url"'`
	Sort string `redis:"sort"`
	Tp int `redis:"tp"`
}


const column_field_cnt = 2

func (this *Column) Load(id int)error{
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "column:" + strconv.Itoa(id)
	values,err := redis.Values(rconn.Do("HGETALL",key))
	if err == nil && len(values) > 0{
		if len(values) == column_field_cnt * 2{
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		}else{
			rconn.Do("DEL",key)
		}

	}
	sql := "select id,title,url,sort,tp from b_column where id=? limit 1"
	db := conn.GetMysqlConn()

	row := db.QueryRow(sql,id)
	err = row.Scan(&this.Id,&this.Title,&this.Url,&this.Sort,&this.Tp)
	if err != nil{
		return err
	}
	rconn.Send("HMSET",redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}


/**
	多线程加载Article对象
 */
func MultipleLoadColumn(id int,position int,column_list []*Column,wg *sync.WaitGroup){
	defer wg.Done()
	column := new(Column)
	err := column.Load(id)
	if err == nil{
		column_list[position] = column
	}
	return
}


func FilterNilColumn(columnList []*Column)[]*Column{
	//过滤空数据
	for k,v := range columnList{
		if v == nil && k  < len(columnList)-1{
			columnList = append(columnList[:k],columnList[k+1:]...)
		}else if k == len(columnList)-1 && v == nil{
			columnList = columnList[:len(columnList)-1]
		}
	}
	return columnList
}