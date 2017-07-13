package kafka

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
)

func TestConsume(t *testing.T) {
	var addrs = []string{
		"10.13.3.25:9092",
		"10.13.3.26:9092",
		"10.13.3.30:9092",
	}
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(addrs, nil)
	if err != nil {
		panic(err)
	}
	(&Consume{
		Client:   client,
		Consumer: consumer,
		Topic:    "dmall_stock_test",
		Group:    "test-group",
		LogPath:  "tmp/topic.log",
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
		Handler: func(m *sarama.ConsumerMessage, logMap map[string]interface{}) {
			logMap[`size`] = 9
			// fmt.Println("%s\n", m.Value)
		},
	}).Start()
}
