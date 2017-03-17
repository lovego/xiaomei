package redis

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

func Do(name string, work func(redis.Conn)) {
	redisConns.RLock()
	redisPool := redisConns.m[name]
	redisConns.RUnlock()
	if redisPool == nil {
		redisPool = &redis.Pool{
			MaxIdle:     32,
			MaxActive:   32,
			IdleTimeout: 600 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					config.DataSource(`redis`, name),
					redis.DialConnectTimeout(time.Second),
					redis.DialReadTimeout(time.Second),
					redis.DialWriteTimeout(time.Second),
				)
			},
		}
		redisConns.Lock()
		redisConns.m[name] = redisPool
		redisConns.Unlock()
	}
	conn := redisPool.Get()
	defer conn.Close()
	work(conn)
}

func SubscribeConn(name string) (redis.Conn, error) {
	return redis.DialURL(
		config.DataSource(`redis`, name),
		redis.DialConnectTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
	)
}
