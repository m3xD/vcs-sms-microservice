package util

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
	"sync"
	"time"
)

var (
	instance *elasticsearch.Client
	err      error
	mutex    sync.Mutex
)

func initES() *elasticsearch.Client {
	var es *elasticsearch.Client
	cfg := elasticsearch.Config{
		APIKey:  os.Getenv("API_KEY"),
		CloudID: os.Getenv("CLOUD_ID"),
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
		fmt.Println(instance.Info())
	}
	return instance
}
