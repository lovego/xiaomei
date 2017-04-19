package redis

import (
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/lovego/xiaomei/config"
)

var pools = struct {
	sync.RWMutex
	m map[string]*redis.Pool
}{m: make(map[string]*redis.Pool)}

func Pool(name string) *redis.Pool {
	pools.RLock()
	p := pools.m[name]
	pools.RUnlock()
	if p == nil {
		p = &redis.Pool{
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
		pools.Lock()
		pools.m[name] = p
		pools.Unlock()
	}
	return p
}

func Do(name string, work func(redis.Conn)) {
	conn := Pool(name).Get()
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
