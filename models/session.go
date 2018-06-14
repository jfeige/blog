package models

import (
	log "github.com/alecthomas/log4go"
	"github.com/garyburd/redigo/redis"
)

var(
	SessionName string
	SessionTime int
)

type Session struct {
	sessid string
}

func NewSession(sessid ...string) *Session {
	if len(sessid) > 0 {
		return &Session{
			sessid: sessid[0],
		}
	} else {
		return &Session{
			sessid: generateSessid(),
		}
	}
}

/**
设置session过期时间
*/
func (this *Session) Expire() {
	rconn := conn.pool.Get()
	defer rconn.Close()

	_,err := rconn.Do("EXPIRE", this.sessid, SessionTime)
	if err != nil{
		log.Error("Session EXPIRE has error:%v", err)
	}
}

/**

 */
func (this *Session) SessionID() string {
	return this.sessid
}

/**
写入session，以哈希的形式写入redis
*/
func (this *Session) SetSession(key string, value interface{}) {
	rconn := conn.pool.Get()
	defer rconn.Close()

	_,err := rconn.Do("HMSET", this.sessid, key, value)
	if err != nil{
		log.Error("SetSession has error:%v", err)
	}
}

/**

 */
func (this *Session) GetSession(key string) interface{} {
	rconn := conn.pool.Get()
	defer rconn.Close()
	values, err := redis.StringMap(rconn.Do("HGETALL", this.sessid))
	if err != nil {
		log.Error("GetSession has error:%v", err)
		return nil
	}
	return values[key]
}

/**
判断session是否存在
*/
func (this *Session) Has(key string) bool {
	rconn := conn.pool.Get()
	defer rconn.Close()
	values, err := redis.StringMap(rconn.Do("HGETALL", this.sessid))
	if err != nil {
		log.Error("Session Has has error:%v", err)
		return false
	}
	_, ok := values[key]
	return ok
}

/**
删除session或者session中指定的字段
*/

func (this *Session) Del(keys ...string) {
	rconn := conn.pool.Get()
	defer rconn.Close()

	var err error
	if len(keys) > 0 {
		_,err = rconn.Do("HDEL", redis.Args{}.Add(this.sessid).AddFlat(keys)...)
	} else {
		_,err = rconn.Do("DEL", this.sessid)
	}
	if err != nil{
		log.Error("Session Del has error:%v", err)
	}
}

func generateSessid() string {
	rstr := RandStr(6)
	return ToMd5(rstr)
}
