package repo

import (
	"github.com/elastic/go-elasticsearch/v8"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticRepo interface {
	Query(query string) (*esapi.Response, error)
}

type ESClient struct {
	*elasticsearch.Client
}

func (escli ESClient) Query(query string) (*esapi.Response, error) {
	return escli.Search(
		escli.Search.WithIndex("vcssms"),
		escli.Search.WithBody(strings.NewReader(query)),
		escli.Search.WithPretty(),
	)
}
