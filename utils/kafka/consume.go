package kafka

import (
	"encoding/json"
	"log"
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
	Group      string
	Handler    func(*sarama.ConsumerMessage, map[string]interface{})
	RedisPool  *redis.Pool
	OffsetsKey string
	LogPath    string
	logFile    *os.File
}

func (c *Consume) Start() {
	c.setupLogFile()
	defer c.logFile.Close()

	if c.OffsetsKey == `` {
		c.OffsetsKey = `kafka-offsets-` + c.Topic + `-` + c.Group
	}
	partitions, err := c.Consumer.Partitions(c.Topic)
	if err != nil {
		panic(err)
	}
	for _, n := range partitions {
		go c.startPartition(n)
	}
	select {}
}

func (c *Consume) startPartition(n int32) {
	offset := c.getPartitionOffset(n)

	pc, err := c.Consumer.ConsumePartition(c.Topic, n, offset)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	for message := range pc.Messages() {
		logMap := map[string]interface{}{
			`partition`: n,
			`at`:        time.Now().Format(utils.ISO8601),
			`now`:       message.Offset,
			`max`:       pc.HighWaterMarkOffset(),
			`bytes`:     len(message.Value),
		}
		c.handle(message, logMap)
		c.setPartitionOffset(n, message.Offset)
		c.writeLog(logMap)
	}
}

func (c *Consume) handle(message *sarama.ConsumerMessage, logMap map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("offset: %d, PANIC: %s\n%s", message.Offset, err, utils.Stack(4))
		}
	}()
	c.Handler(message, logMap)
}

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

func (c *Consume) setupLogFile() {
	if dir := filepath.Dir(c.LogPath); dir != `.` {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
	if logFile, err := fs.OpenAppend(c.LogPath); err != nil {
		log.Fatal(err)
	} else {
		c.logFile = logFile
	}
}

func (c *Consume) writeLog(m map[string]interface{}) {
	buf, err := json.Marshal(m)
	if err != nil {
		log.Printf("marshal log err: %v", err)
		return
	}
	buf = append(buf, '\n')
	_, err = c.logFile.Write(buf)
	if err != nil {
		log.Printf("write log err: %v", err)
		return
	}
}
