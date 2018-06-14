package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
)

//类别
type Category struct {
	Id          int    `redis:"id"`
	Name        string `redis:"name"`
	Article_cnt int    `redis:"article_cnt"`
	Sort        string `redis:"sort"`
}

const category_field_cnt = 4

func (this *Category) Load(id int) error {
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "category:" + strconv.Itoa(id)
	values, err := redis.Values(rconn.Do("HGETALL", key))
	if err == nil && len(values) > 0 {
		if len(values) == category_field_cnt*2 {
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		} else {
			rconn.Do("DEL", key)
		}

	}
	sql := "select id,name,article_cnt,sort from b_category where id=? limit 1"
	db := conn.GetMysqlConn()

	row := db.QueryRow(sql, id)
	err = row.Scan(&this.Id, &this.Name, &this.Article_cnt, &this.Sort)
	if err != nil {
		return err
	}
	rconn.Send("HMSET", redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}

/**
多线程加载Category对象
*/
func MultipleLoadCategory(id int, position int, category_list []*Category, wg *sync.WaitGroup) {
	defer wg.Done()
	category := new(Category)
	err := category.Load(id)
	if err == nil {
		category_list[position] = category
	}
	return
}

func FilterNilCategory(categoryList []*Category) []*Category {
	//过滤空数据
	for k, v := range categoryList {
		if v == nil && k < len(categoryList)-1 {
			categoryList = append(categoryList[:k], categoryList[k+1:]...)
		} else if k == len(categoryList)-1 && v == nil {
			categoryList = categoryList[:len(categoryList)-1]
		}
	}
	return categoryList
}
