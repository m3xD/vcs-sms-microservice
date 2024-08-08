package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"healthcheck-server/service"
	"healthcheck-server/util"
	"time"
)

type HealthCheckRequest struct {
	IP       string `json:"ip"`
	Time     int64  `json:"time"`
	Duration int64  `json:"duration"`
}

type HealthCheckResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HealthCheck(c *gin.Context) {
	var req HealthCheckRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, HealthCheckResponse{
			Code:    400,
			Message: "Bad Request :" + err.Error(),
		})
		return
	}

	kafka := service.NewKafkaService(util.GetConsumer(), util.GetProducer())
	topic := "healthcheck"
	req.Time = time.Now().UnixMilli()
	bytes, _ := json.Marshal(req)
	err := kafka.SendMessage(topic, bytes)
	if err != nil {
		c.JSON(500, HealthCheckResponse{
			Code:    500,
			Message: "Internal Server Error :" + err.Error(),
		})
		return
	}

	c.JSON(200, HealthCheckResponse{
		Code:    200,
		Message: "OK",
	})
}
