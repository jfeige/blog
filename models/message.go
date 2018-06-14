package models

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"time"
)

//留言
type Message struct {
	Id      int    `redis:"id"`
	User    string `redis:"user"`
	Content string `redis:"content"`
	Atime   int64  `redis:"atime"`
	Mdate   int    `redis:"mdate"`
}

const message_field_cnt = 5

func (this *Message) Load(id int) error {
	rconn := conn.GetRedisConn()
	defer rconn.Close()

	key := "message:" + strconv.Itoa(id)
	values, err := redis.Values(rconn.Do("HGETALL", key))
	if err == nil && len(values) > 0 {
		if len(values) == message_field_cnt*2 {
			err = redis.ScanStruct(values, this)
			if err == nil {
				return nil
			}
		} else {
			rconn.Do("DEL", key)
		}

	}
	sql := "select id,user,content,atime,mdate from b_messageboard where id=? limit 1"
	db := conn.GetMysqlConn()

	row := db.QueryRow(sql, id)
	err = row.Scan(&this.Id, &this.User, &this.Content, &this.Atime, &this.Mdate)
	if err != nil {
		return err
	}
	rconn.Send("HMSET", redis.Args{}.Add(key).AddFlat(this)...)
	return nil
}

/**
多线程加载Article对象
*/
func MultipleLoadMessage(id int, position int, msg_list []*Message, wg *sync.WaitGroup) {
	defer wg.Done()
	msg := new(Message)
	err := msg.Load(id)
	if err == nil {
		msg_list[position] = msg
	}
	return
}

/**
过滤空数据
*/
func FilterNilMessage(msgList []*Message) []*Message {
	//过滤空数据
	for k, v := range msgList {
		if v == nil && k < len(msgList)-1 {
			msgList = append(msgList[:k], msgList[k+1:]...)
		} else if k == len(msgList)-1 && v == nil {
			msgList = msgList[:len(msgList)-1]
		}
	}
	return msgList
}

/**
按照给定的格式格式化日期和时间
*/
func (this *Message) FormatPublishTime(format string) string {

	return time.Unix(this.Atime, 0).Format(format)
}

/**
前台摘要显示
*/
func (this *Message) FormatContent() string {
	content := []rune(this.Content)
	if len(content) > 25 {
		return string(content[:25]) + " ..."
	}
	return this.Content
}
