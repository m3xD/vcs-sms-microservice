package service

import (
	"context"
	"encoding/json"
	"fmt"
	db "server-management/db/sqlc"
	"server-management/util"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-co-op/gocron/v2"
	_ "time/tzdata"
)

type ReportService struct {
	Elastic *Elastic
	Redis   *Redis
	Store   db.Store
	Mail    *Mail
	Kafka   *KafkaService
}

func NewReportService(elastic *Elastic, redis *Redis, store db.Store, mail *Mail, kafka *KafkaService) *ReportService {
	return &ReportService{
		Elastic: elastic,
		Redis:   redis,
		Store:   store,
		Mail:    mail,
		Kafka:   kafka,
	}
}

func (report *ReportService) SendReport(start int64, end int64, email string) error {
	log := util.NewLogger()

	upTime := report.Elastic.GetAVG(start, end)

	allServer, err := report.Redis.GetServer("servers")
	var servers []db.Server
	if err != nil {
		log.Error("error when get server from redis: " + err.Error())
		servers, err = report.Store.GetAllServers(context.Background())
		if err != nil {
			log.Error("error when get server from db: " + err.Error())
		}
		serversRedis, err := json.Marshal(servers)
		if err != nil {
			log.Error("error when marshal servers: " + err.Error())
		}
		err = report.Redis.SetServer("servers", string(serversRedis), 15*time.Minute)
	} else {
		err = json.Unmarshal([]byte(allServer), &servers)
		if err != nil {
			log.Error("error when unmarshal servers: " + err.Error())
		}
	}

	on := 0
	off := 0
	for _, server := range servers {
		if server.Status == 0 {
			off++
		} else {
			on++
		}
	}

	var body string
	body += fmt.Sprintf("Number of server on %d\n", on)
	body += fmt.Sprintf("Number of server off %d\n", off)
	for _, data := range upTime {
		if data.ID == 0 {
			continue
		}
		body += fmt.Sprintf("Server: %d, Uptime: %f\n", data.ID, data.AVGUpTime)
	}
	p := report.Kafka.p
	send, _ := json.Marshal(MailSent{
		Email: email,
		Body:  body,
	})
	topic := "vcs"
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          send,
	}, nil)
	return err
}

type MailSent struct {
	Email string
	Body  string
}

func (report *ReportService) SendEveryDayReport() {
	log := util.NewLogger()
	timezone, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Error("error occur when load timezone: " + err.Error())
	}
	s, err := gocron.NewScheduler(gocron.WithLocation(timezone))
	if err != nil {
		log.Error("error occur when schedule: " + err.Error())

	}
	j, err := s.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(23, 59, 00))), gocron.NewTask(report.SendReport, time.Now().AddDate(0, 0, -1).UnixMilli(), time.Now().UnixMilli(), "thaoproduction123@gmail.com"))

	if err != nil {
		log.Error("error occur when add new job to schedule: " + err.Error())
	}
	log.Info(fmt.Sprintf("start schedule with job id: %d", j.ID))
	s.Start()
}
