package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
)

func (c *Consume) getPartitionOffset(n int32) int64 {
	offset := c.getPartitionRedisOffset(n)

	// NOTE: 由于kafka容量限制可能丢弃部分消息导致redis缓存的offset失效
	//       我们要检查并修正redis缓存offset在[oldest,newest]有效区间
	if oldestOffset, err := c.Client.GetOffset(c.Topic, n, sarama.OffsetOldest); err == nil {
		if offset <= oldestOffset {
			offset = oldestOffset
		}
	}
	if newestOffset, err := c.Client.GetOffset(c.Topic, n, sarama.OffsetNewest); err == nil {
		if offset >= newestOffset {
			offset = newestOffset
		}
	}
	return offset
}

func (c *Consume) getPartitionRedisOffset(n int32) int64 {
	conn := c.RedisPool.Get()
	defer conn.Close()
	v, err := conn.Do(`HGET`, c.OffsetsKey, n)
	if err != nil && err != redis.ErrNil {
		log.Printf("partition %d get offset error: %s", n, err)
	}
	if v == nil {
		return sarama.OffsetOldest
	}
	offset, _ := redis.Int64(v, err)
	return offset + 1
}

func (c *Consume) setPartitionOffset(n int32, offset int64) {
	conn := c.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do(`HSET`, c.OffsetsKey, n, offset)
	if err != nil {
		log.Printf("partition %d set offset(%d) error: %s", n, offset, err)
	}
}
