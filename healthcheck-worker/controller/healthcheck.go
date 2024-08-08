package controller

import (
	"encoding/json"
	"fmt"
	"healthcheck-worker/model"
	"healthcheck-worker/service"
	"healthcheck-worker/util"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var statusMapping *util.ConcurrentMap

type HealthCheckRequest struct {
	IP       string `json:"ip"`
	Time     int64  `json:"time"`
	Duration int64  `json:"duration"`
}

type HealthCheckController struct {
	elastic *service.ESService
	db      *service.ServerService
	redis   *service.Redis
}

type DocElastic struct {
	Server   *model.Server
	Time     int64
	Duration int64
}

func NewHealthCheckController(elastic *service.ESService, db *service.ServerService, redis *service.Redis) *HealthCheckController {
	return &HealthCheckController{
		elastic: elastic,
		db:      db,
		redis:   redis,
	}
}

func (h *HealthCheckController) HealthCheck() {
	log := util.NewLogger()
	// get information from message broker
	k := service.NewKafkaService(util.GetConsumer(), util.GetProducer())
	megs := make(chan *kafka.Message, 100)

	// init worker
	workers := 10
	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case msg := <-megs:
					// check trong redis -> nếu có thì lấy ra, ko thì lấy trong db va update trong redis
					//
					m := &HealthCheckRequest{}
					if err := json.Unmarshal(msg.Value, m); err != nil {
						log.Error("error when unmarshal mail request: " + err.Error())
						continue
					}
					statusMapping.Set(m.IP, 1)
					// find server in redis
					server := &model.Server{}
					redis := h.redis
					serverString, err := redis.GetServer("server:" + m.IP)
					var tmp interface{}
					if err != nil { // not found in redis
						log.Error("error when get server from redis: " + err.Error())
						// get server from db
						tmp = h.db.GetServerByIP(m.IP)
						if tmp == nil { // if server is not in db
							log.Error("server not found")
							server = &model.Server{
								Name:        m.IP,
								Status:      1,
								IP:          m.IP,
								CreatedTime: time.Now(),
								LastUpdated: time.Now(),
							}
							err = h.db.CreateServer(server)
							log.Info(fmt.Sprintf("Server created: %v", server))
							if err != nil {
								log.Error("error when create server: " + err.Error())
							}
						}
					} else {
						err = json.Unmarshal([]byte(serverString), server)
						server.LastUpdated = time.Now()
						server.Status = 1
						if err != nil {
							log.Error("error when unmarshal server from redis: " + err.Error())
						}
					}
					// cache to redis
					serverMarshal, err := json.Marshal(server)
					if err != nil {
						log.Error("error when marshal server: " + err.Error())
					}
					err = redis.SetServer("server:"+m.IP, serverMarshal, time.Minute*10)
					if err != nil {
						log.Error("error when set server to redis: " + err.Error())
					}
					log.Info("server cached to redis")
					// insert doc in elastic
					doc := &DocElastic{
						Server:   server,
						Time:     m.Time,
						Duration: m.Duration,
					}
					h.elastic.InsertInBatch(doc)
				}
			}
		}()
	}

	go func() {
		for {
			// read message
			msg, err := k.C.ReadMessage(-1)
			if err != nil {
				log.Error("error when consume message: " + err.Error())
			} else {
				megs <- msg
			}
		}
	}()
}

func (h *HealthCheckController) UpdateCurrentStatus() {
	log := util.NewLogger()
	// get all server from db
	db := h.db
	statusMapping = util.NewConcurrentMap()
	// update status of all server
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info("update status of all server")
			err := db.UpdateServersOn(statusMapping.Items())
			if err != nil {
				log.Error("error when update status of all server: " + err.Error())
			}
		}
	}
}
