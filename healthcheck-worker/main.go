package main

import (
	"healthcheck-worker/controller"
	"healthcheck-worker/repo"
	"healthcheck-worker/service"
	"healthcheck-worker/util"
)

func main() {
	healthcheck := controller.NewHealthCheckController(service.NewESService(&repo.ESClient{Client: util.GetES()}),
		service.NewServerService(service.NewDatabase()),
		service.NewRedis())

	healthcheck.HealthCheck()

	go func() {
		healthcheck.UpdateCurrentStatus()
	}()

	go func() {

	}()

	select {}
}
