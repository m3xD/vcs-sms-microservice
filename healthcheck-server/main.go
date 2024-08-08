package main

import (
	"github.com/gin-gonic/gin"
	"healthcheck-server/controller"
	"healthcheck-server/util"
)

func main() {

	log := util.NewLogger()

	log.Info("Starting server at port 8082")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/healthcheck", controller.HealthCheck)
	err := router.Run(":8082")
	if err != nil {
		log.Error("error when start server at port 8081: " + err.Error())
		return
	}
}
