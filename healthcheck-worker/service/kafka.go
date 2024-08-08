package service

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaService struct {
	C *kafka.Consumer
	P *kafka.Producer
}

type KafkaMessage struct {
	Id        string `json:"id"`
	SendId    string `json:"sendID"`
	ReceiveId string `json:"receiveID"`
	Duration  int    `json:"duration"`
	Timestamp int64  `json:"timestamp"`
}

func NewKafkaService(c *kafka.Consumer, p *kafka.Producer) *KafkaService {
	return &KafkaService{C: c, P: p}
}

func (k *KafkaService) ConsumeMessage() (*KafkaMessage, error) {
	msg, err := k.C.ReadMessage(-1)
	if err == nil {
		message := &KafkaMessage{}
		err := json.Unmarshal(msg.Value, message)
		if err != nil {
			return nil, err
		}
		return message, nil
	} else {
		return nil, err
	}
}

func (k *KafkaService) SendMessage(topic string, value []byte) error {
	return k.P.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
}
