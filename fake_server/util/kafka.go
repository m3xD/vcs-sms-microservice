package util

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	err   error
	once  sync.Once
	once2 sync.Once
	c     *kafka.Consumer
	p     *kafka.Producer
)

func GetProducer() *kafka.Producer {
	if p == nil {
		once2.Do(func() {
			p, err = kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": "localhost:9092",
			})

			if err != nil {
				panic(err)
			}
		})
	}
	return p
}
