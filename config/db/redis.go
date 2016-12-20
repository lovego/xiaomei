package db

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisConns map[string]*redis.Pool

func init() {
	if redisConns == nil {
		redisConns = make(map[string]*redis.Pool)
	}
}

func RedisDo(name string, work func(redis.Conn)) {
	redisPool := redisConns[name]
	if redisPool == nil {
		redisPool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 600 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					config.Redis()[name],
					redis.DialConnectTimeout(time.Second),
					redis.DialReadTimeout(time.Second),
					redis.DialWriteTimeout(time.Second),
				)
			},
		}
		redisConns[name] = redisPool
	}
	conn := redisPool.Get()
	defer conn.Close()
	work(conn)
}

func RedisSubscribeConn(name string) (redis.Conn, error) {
	return redis.DialURL(
		config.Redis()[name],
		redis.DialConnectTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
	)
}
