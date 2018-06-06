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
func DelKeys(key string){
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