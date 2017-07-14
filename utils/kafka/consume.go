package kafka

import (
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"github.com/lovego/xiaomei/utils"
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
		c.process(pc, n, message)
	}
}

func (c *Consume) process(pc sarama.PartitionConsumer, n int32, message *sarama.ConsumerMessage) {
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

func (c *Consume) handle(message *sarama.ConsumerMessage, logMap map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("offset: %d, PANIC: %s\n%s", message.Offset, err, utils.Stack(4))
		}
	}()
	c.Handler(message, logMap)
}
