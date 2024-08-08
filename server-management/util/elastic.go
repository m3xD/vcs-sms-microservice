package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

var (
	instance *elasticsearch.Client
	err      error
	mutex    sync.Mutex
)

func initES() *elasticsearch.Client {
	var es *elasticsearch.Client
	cfg := elasticsearch.Config{
		CloudID: os.Getenv("CLOUD_ID"),
		APIKey:  os.Getenv("API_KEY"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	_, err = es.Info(es.Info.WithContext(ctx))
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	return es
}
func GetES() *elasticsearch.Client {
	mutex.Lock()
	defer mutex.Unlock()
	if instance == nil {
		instance = initES()
		fmt.Println(time.Now())
	}
	return instance
}
