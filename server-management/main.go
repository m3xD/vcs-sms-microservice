package main

import (
	"context"
	"log"
	"os"
	"server-management/controller"
	db "server-management/db/sqlc"
	"server-management/service"
	"server-management/util"

	"github.com/jackc/pgx/v5"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		//log.Fatal("cannot load config file")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_SOURCE"))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal(err)
	}
	reportService := service.NewReportService(
		service.NewElastic(),
		service.NewRedis(),
		store,
		service.NewMailService(),
		service.NewKafkaService(util.GetConsumer(), util.GetProducer()),
	)

	reportService.SendEveryDayReport()

	err = server.Start()

	if err != nil {
		log.Fatal(err)
	}
}
