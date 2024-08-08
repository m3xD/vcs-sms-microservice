package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"log"
	"server-management/util"
	"strconv"
)

type Elastic struct {
	elasticClient *elasticsearch.Client
}

func NewElastic() *Elastic {
	return &Elastic{elasticClient: util.GetES()}
}

type UptimeResponse struct {
	ID        int
	AVGUpTime float64
}

func (elastic *Elastic) GetAVG(start int64, end int64) []UptimeResponse {
	es := elastic.elasticClient
	var buf bytes.Buffer
	// trong 1 days co nhieu record tu 1 server
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"Time": map[string]interface{}{
					"gte": start,
					"lte": end,
				},
			},
		},
		"aggs": map[string]interface{}{
			"servers": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "Server.ID.keyword",
					"size":  100000,
				},
				"aggs": map[string]interface{}{
					"sum_duration": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "Duration",
						},
					},
					"min_time": map[string]interface{}{
						"min": map[string]interface{}{
							"field": "Time",
						},
					},
					"uptime_avg": map[string]interface{}{
						"bucket_script": map[string]interface{}{
							"buckets_path": map[string]interface{}{
								"uptime":   "sum_duration",
								"min_time": "min_time",
							},
							"script": "params.min_time == null ? 100 : (params.uptime / ((new Date().getTime() / 1000) - (params.min_time / 1000))) * 100",
						},
					},
				},
			},
		},
	}
	err := json.NewEncoder(&buf).Encode(query)

	if err != nil {
		log.Println(err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("vcssms"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty())
	if err != nil {
		log.Println(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	var upTimeServers []UptimeResponse
	if res.IsError() {
		log.Println(res.String())
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Println(err.Error())
		} else {
			log.Printf(
				"[%s] %d hits; took: %dms",
				res.Status(),
				int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
				int(r["took"].(float64)),
			)

			for _, result := range r["aggregations"].(map[string]interface{})["servers"].(map[string]interface{})["buckets"].([]interface{}) {
				s, _ := strconv.Atoi(result.(map[string]interface{})["key"].(string))
				upTimeServers = append(upTimeServers, UptimeResponse{
					ID:        s,
					AVGUpTime: result.(map[string]interface{})["uptime_avg"].(map[string]interface{})["value"].(float64),
				})
			}
			fmt.Println()
		}
	}
	return upTimeServers
}
