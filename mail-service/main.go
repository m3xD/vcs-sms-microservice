package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"mail-service/service"
	"mail-service/util"
)

func main() {
	log := util.NewLogger()

	c := util.GetConsumer()
	mailRequest := make(chan *kafka.Message, 100)
	mailService := service.NewMailService()
	log.Info("Start consumer")
	defer c.Close()
	go func() {
		for {
			select {
			case msg := <-mailRequest:
				m := &MailSent{}
				err := json.Unmarshal(msg.Value, m)
				if err != nil {
					log.Error("error when unmarshal mail request: " + err.Error())
					continue
				}
				err = mailService.SendEmail([]string{m.Email}, m.Body)
				if err != nil {
					log.Error("error when send email: " + err.Error())
					continue
				}
				log.Info("Sent email to admin: " + m.Email)
			}
		}
	}()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Info("received message: " + string(msg.Value))
			mailRequest <- msg
		} else {
			log.Error("consumer is broken: " + err.Error())
		}
	}
}

type MailSent struct {
	Email string
	Body  string
}
