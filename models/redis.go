package models

import (
	"errors"
	log "github.com/alecthomas/log4go"
	"time"
	"github.com/garyburd/redigo/redis"
)

var(
	raddress string
	rpasswd string
	rmaxidle int
	rmaxactive int
	rtimeout int
)

func initRedisConfig()error{

	rpasswd = lcf.String("redis::rpasswd")
	if muser == "" {
		return errors.New("Can't not find redis parameters:rpasswd")
	}
	raddress = lcf.String("redis::raddress")
	if raddress == "" {
		return errors.New("Can't not find redis parameters:raddress")
	}

	rmaxidle,err = lcf.Int("redis::rmaxidle")
	if rmaxidle == 0 {
		return errors.New("Can't not find redis parameters:rmaxidle")
	}
	rmaxactive,err = lcf.Int("redis::rmaxactive")
	if rmaxidle == 0 {
		return errors.New("Can't not find redis parameters:rmaxactive")
	}

	rtimeout,err = lcf.Int("redis::rtimeout")
	if rtimeout == 0 {
		return errors.New("Can't not find redis parameters:rtimeout")
	}
	return nil
}

func initRedisPool() *redis.Pool {
	Pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", raddress)
			if err != nil {
				log.Exit("Con't init Redis Pool.Error:", err)
				return nil, err
			}
			err = conn.Send("AUTH", rpasswd)
			if err != nil {
				log.Exit("Con't Auth Redis.Error:", err)
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			err := conn.Send("PING")
			return err
		},
	}
	return Pool
}

/**
	模糊匹配，删除多个key
 */
func BlurDelKeys(key string){
	rconn := conn.pool.Get()
	defer rconn.Close()

	keys,err := redis.Values(rconn.Do("keys",key))

	if err != nil{
		log.Error("DelKeys has error!key:%s,error:%v",key,err)
		return
	}
	if len(keys) > 0{
		_,err = rconn.Do("DEL",keys...)
		if err != nil{
			log.Error("DelKeys has error!key:%s,error:%v",key,err)
			return
		}
	}
}

/**
	删除单个key
 */
func DelKey(key string){
	rconn := conn.pool.Get()
	defer rconn.Close()

	_,err = rconn.Do("DEL",key)
	if err != nil{
		log.Error("DelKey has error!key:%s,error:%v",key,err)
		return
	}
}

/**
	删除多个key
 */
func DelKeys(key []interface{}){
	rconn := conn.pool.Get()
	defer rconn.Close()

	_,err = rconn.Do("DEL",key...)
	if err != nil{
		log.Error("DelKeys has error!key:%s,error:%v",key,err)
		return
	}
}

/**
	从zset中删除单个数据
 */

func DelZsetData(key string,value interface{}){
	rconn := conn.pool.Get()
	defer rconn.Close()

	rconn.Do("ZREM",key,value)

}

/**
	zset添加单个数据
 */

func AddZsetData(key string,score,value interface{}){
	rconn := conn.pool.Get()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if exists{
		rconn.Do("ZADD",key,score,value)
		return
	}
}


/**
	数字减1
 */
func Decr(key string){
	rconn  := conn.pool.Get()
	defer rconn.Close()

	ret,_ := redis.Int(rconn.Do("DECR",key))
	if ret < 0 {
		rconn.Do("DEL",key)
		return
	}
}

/**
	数字加1
 */
func Incr(key string){
	rconn  := conn.pool.Get()
	defer rconn.Close()

	exists,_ := redis.Bool(rconn.Do("EXISTS",key))
	if exists {
		rconn.Do("INCR", key)
	}
}

/**
	HMset
 */
func HMset(args []interface{}){
	rconn := conn.pool.Get()
	defer rconn.Close()

	rconn.Do("HMSET",args...)
}