package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"healthcheck-worker/repo"
	"healthcheck-worker/util"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
)

type ESService struct {
	escli repo.ElasticRepo
	bi    esutil.BulkIndexer
}

func NewESService(escli repo.ElasticRepo) *ESService {
	log := util.NewLogger()
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: util.GetES(),
		Index:  "vcssms",
	})
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the indexer: %s", err))
		log.Fatal("Shutting down")
	}
	return &ESService{
		escli: escli,
		bi:    bi,
	}
}

func (service ESService) InsertInBatch(doc interface{}) {
	log := util.NewLogger()
	data, err := json.Marshal(doc)
	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling the document: %s", err))
		return
	}
	err = service.bi.Add(context.Background(), esutil.BulkIndexerItem{
		Action:     "create",
		DocumentID: uuid.New().String(),
		Body:       bytes.NewReader(data),
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			log.Info(fmt.Sprintf("Document added to the indexer: %s", res.Result))
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error adding the document to the indexer: %s", err))
	}
}
