package kafka

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type Consume struct {
	Client     sarama.Client
	Consumer   sarama.Consumer
	Topic      string
	LogDir     string
	RedisPool  *redis.Pool
	Handler    func(*sarama.ConsumerMessage, *os.File)
	OffsetsKey string
	Group      string
}

func (c *Consume) Start() {
	if err := os.MkdirAll(c.LogDir, 0755); err != nil {
		panic(err)
	}
	if c.OffsetsKey == `` {
		c.OffsetsKey = `kafka-offsets-` + c.Topic + `-` + c.Group
	}
	partitions, err := c.Consumer.Partitions(c.Topic)
	if err != nil {
		panic(err)
	}
	for _, n := range partitions {
		go c.StartPartition(n)
	}
	select {}
}

func (c *Consume) StartPartition(n int32) {
	// TODO: error handling ignored
	logFile, _ := fs.OpenAppend(filepath.Join(c.LogDir, fmt.Sprintf(`%d.log`, n)))
	defer logFile.Close()

	offset := c.GetPartitionOffset(n, logFile)

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

	pc, err := c.Consumer.ConsumePartition(c.Topic, n, offset)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	for message := range pc.Messages() {
		fmt.Fprintf(logFile, "%s %d %d %d\n", time.Now().Format(utils.ISO8601),
			pc.HighWaterMarkOffset(), message.Offset, len(message.Value))
		c.Handle(message, logFile)
		c.SetPartitionOffset(n, message.Offset, logFile)
	}
}

func (c *Consume) Handle(message *sarama.ConsumerMessage, logFile *os.File) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(logFile, "%s PANIC: %s\n%s",
				time.Now().Format(utils.ISO8601), err, utils.Stack(4),
			)
		}
	}()
	c.Handler(message, logFile)
}

func (c *Consume) GetPartitionOffset(n int32, logFile *os.File) int64 {
	conn := c.RedisPool.Get()
	defer conn.Close()
	v, err := conn.Do(`HGET`, c.OffsetsKey, n)
	if err != nil && err != redis.ErrNil {
		fmt.Fprintf(logFile, "%s get offset error: %s\n", time.Now().Format(utils.ISO8601), err)
	}
	if v == nil {
		return sarama.OffsetOldest
	}
	offset, _ := redis.Int64(v, err)
	return offset + 1
}

func (c *Consume) SetPartitionOffset(n int32, offset int64, logFile *os.File) {
	conn := c.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do(`HSET`, c.OffsetsKey, n, offset)
	if err != nil {
		fmt.Fprintf(logFile, "%s set offset error: %s (offset: %d)\n",
			time.Now().Format(utils.ISO8601), err, offset,
		)
	}
}
