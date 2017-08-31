package kafka

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/alarm"
	"github.com/lovego/xiaomei/utils/fs"
)

type Consume struct {
	Client     sarama.Client
	Consumer   sarama.Consumer
	Topic      string
	Group      string
	Handler    func(*sarama.ConsumerMessage, map[string]interface{}) error
	RedisPool  *redis.Pool
	OffsetsKey string

	LogPath   string
	logWriter io.Writer
	Alarm     alarm.Engine
}

func (c *Consume) Start() {
	c.logWriter = fs.NewLogFile(c.LogPath)

	if c.OffsetsKey == `` {
		c.OffsetsKey = `kafka-offsets-` + c.Topic + `-` + c.Group
	}
	partitions, err := c.Consumer.Partitions(c.Topic)
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
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
	defer func() {
		if err := recover(); err != nil {
			c.Alarm.Alarmf("PANIC: %v\n%s", err, utils.Stack(1))
		}
	}()
	if err := c.Handler(message, logMap); err != nil {
		c.Alarm.Alarmf("error: %v\n%s", err, utils.Stack(4))
	}
}

func (c *Consume) writeLog(m map[string]interface{}) {
	buf, err := json.Marshal(m)
	if err != nil {
		log.Printf("marshal log err: %v", err)
		return
	}
	buf = append(buf, '\n')
	_, err = c.logWriter.Write(buf)
	if err != nil {
		log.Printf("write log err: %v", err)
		return
	}
}
