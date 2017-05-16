package kafka

import (
	"os"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
)

func TestConsume(t *testing.T) {
	c, err := sarama.NewConsumer([]string{
		"10.13.3.25:9092", "10.13.3.26:9092", "10.13.3.30:9092",
	}, nil)
	if err != nil {
		panic(err)
	}
	(&Consume{
		Consumer: c,
		Topic:    "dmall_stock_test",
		LogDir:   ".",
		RedisPool: &redis.Pool{
			MaxIdle:     5,
			MaxActive:   5,
			IdleTimeout: 600 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					`redis://:@192.168.56.56:6379/1`,
					redis.DialConnectTimeout(time.Second),
					redis.DialReadTimeout(time.Second),
					redis.DialWriteTimeout(time.Second),
				)
			},
		},
		Handler: func(m *sarama.ConsumerMessage, f *os.File) {
			// fmt.Println("%s\n", m.Value)
		},
	}).Start()
}
