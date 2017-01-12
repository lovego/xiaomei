package db

import (
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/garyburd/redigo/redis"
)

var redisConns = struct {
	sync.RWMutex
	m map[string]*redis.Pool
}{m: make(map[string]*redis.Pool)}

func RedisDo(name string, work func(redis.Conn)) {
	redisConns.RLock()
	redisPool := redisConns[name]
	redisConns.RUnlock()
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
		redisConns.Lock()
		redisConns[name] = redisPool
		redisConns.Unlock()
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
