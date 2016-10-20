package config

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisPool *redis.Pool

func RedisDo(work func(redis.Conn)) {
	if redisPool == nil {
		redisPool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 600 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					Data.Redis,
					redis.DialConnectTimeout(time.Second),
					redis.DialReadTimeout(time.Second),
					redis.DialWriteTimeout(time.Second),
				)
			},
		}
	}
	conn := redisPool.Get()
	defer conn.Close()
	work(conn)
}

func RedisSubscribeConn() (redis.Conn, error) {
	return redis.DialURL(
		Data.Redis,
		redis.DialConnectTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
	)
}
