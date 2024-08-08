package util

import (
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	e     error
	once  sync.Once
	once2 sync.Once
	c     *kafka.Consumer
	p     *kafka.Producer
)

func GetConsumer() *kafka.Consumer {
	if c == nil {
		once.Do(func() {
			c, e = kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": os.Getenv("KAFKA_ADDRESS"),
				"group.id":          "worker18",
				"auto.offset.reset": "earliest",
			})

			if e != nil {
				panic(e)
			}
			if err := c.SubscribeTopics([]string{"in"}, nil); err != nil {
				panic(err)
			}
		})
	}
	return c
}

func GetProducer() *kafka.Producer {
	if p == nil {
		once2.Do(func() {
			p, e = kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": os.Getenv("KAFKA_ADDRESS"),
			})

			if e != nil {
				panic(e)
			}
		})
	}
	return p
}
