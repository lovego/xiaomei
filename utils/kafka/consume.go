package kafka

import (
	"encoding/json"
	"io"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
	"github.com/lovego/xiaomei/utils/logger"
)

type Consume struct {
	Client     sarama.Client
	consumer   sarama.Consumer
	Topic      string
	Group      string
	Handler    func(*sarama.ConsumerMessage, map[string]interface{}) error
	RedisPool  *redis.Pool
	OffsetsKey string

	LogPath   string
	Logger    *logger.Logger
	logWriter io.Writer
}

func (c *Consume) Start() {
	if consumer, err := sarama.NewConsumerFromClient(c.Client); err == nil {
		c.consumer = consumer
	} else {
		c.Logger.Panic(err)
	}

	c.logWriter = fs.NewLogFile(c.LogPath)

	if c.OffsetsKey == `` {
		c.OffsetsKey = `kafka-offsets-` + c.Topic + `-` + c.Group
	}
	partitions, err := c.Client.Partitions(c.Topic)
	if err != nil {
		c.Logger.Panic(err)
	}
	for _, n := range partitions {
		go c.startPartition(n)
	}
	select {}
}

func (c *Consume) startPartition(n int32) {
	offset := c.getPartitionOffset(n)

	pc, err := c.consumer.ConsumePartition(c.Topic, n, offset)
	if err != nil {
		c.Logger.Panic(err)
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
	c.callHandler(message, logMap)
	c.setPartitionOffset(n, message.Offset)
	c.writeLog(logMap)
}

func (c *Consume) callHandler(message *sarama.ConsumerMessage, logMap map[string]interface{}) {
	defer c.Logger.Recover()
	if err := c.Handler(message, logMap); err != nil {
		c.Logger.Printf("consume handler error: %v", err)
	}
}

func (c *Consume) writeLog(m map[string]interface{}) {
	buf, err := json.Marshal(m)
	if err != nil {
		c.Logger.Printf("marshal log err: %v", err)
		return
	}
	buf = append(buf, '\n')
	_, err = c.logWriter.Write(buf)
	if err != nil {
		c.Logger.Printf("write log err: %v", err)
		return
	}
}
