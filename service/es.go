package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"logging-helper/config"
	"os"
	"strings"
	"time"

	"logging-helper/utils/sse"

	"github.com/olivere/elastic/v7"
)

type Elastic struct {
}

var Client *elastic.Client

type LogMessage struct {
	Log string `json:"log"`
}

func InitElasticClinet() (err error) {
	Client, err = elastic.NewClient(
		elastic.SetURL(config.EsHost),
		elastic.SetSniff(config.Sniff),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if err != nil {
		return err
	}
	info, code, err := Client.Ping(config.EsHost).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := Client.ElasticsearchVersion(config.EsHost)
	if err != nil {
		return err
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	return nil
}

func (es *Elastic) Query(taskName string) (res []string) {

	pq := elastic.NewMatchPhraseQuery("kubernetes.labels.ev-logger-sign", taskName)
	bQ := elastic.NewBoolQuery()
	bQ.Must(pq)
	numLine := 10000
	searchRes, err := Client.Search().Index(config.EsIndex).Query(bQ).Size(numLine).SortBy(elastic.NewFieldSort("@timestamp").Asc()).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	res = EachTime(searchRes)
	return res
}

func (es *Elastic) QueryBySse(taskName string, stream *sse.Event) {

	rQ := elastic.NewRangeQuery("@timestamp")

	start := time.Now().AddDate(-1, 0, 0)
	end := time.Now()
	rQ.Gte(start.Format("2006-01-02T15:04:05Z07:00"))

	rQ.Lt(end.Format("2006-01-02T15:04:05Z07:00"))
	pq := elastic.NewMatchPhraseQuery(config.TaskName, taskName)
	//pq := elastic.NewMatchPhraseQuery("kubernetes.labels.k8s-app", taskName)

	bQ := elastic.NewBoolQuery()
	bQ.Must(rQ, pq)
queryLoop:
	for {
		numLine := 10000

		searchRes, err := Client.Search(config.EsIndex).Query(bQ).Size(numLine).SortBy(elastic.NewFieldSort("@timestamp").Asc()).Do(context.Background())
		if err != nil {
			panic(err)
		}

		var tempstr string
		for _, tempstr := range EachTime(searchRes) {
			if strings.Contains(tempstr, "[91m") && strings.Contains(tempstr, "[0m") {
				continue
			}
			if len(tempstr) <= 3 {
				tempstr = ""
			}
			if tempstr == "\n" {
				continue
			}
			stream.Message <- tempstr
		}

		if stream.Stop == 1 {
			break queryLoop
		}

		if tempstr != "" {
			start = end
		}
		time.Sleep(3 * time.Second)
		end = time.Now()
		rQ.Gte(start.Format("2006-01-02T15:04:05Z07:00"))
		rQ.Lt(end.Format("2006-01-02T15:04:05Z07:00"))
		pq = elastic.NewMatchPhraseQuery(config.TaskName, taskName)

		//pq := elastic.NewMatchPhraseQuery("kubernetes.labels.k8s-app", taskName)
		bQ = elastic.NewBoolQuery()
		bQ.Must(rQ, pq)
	}
	return
}

func EachTime(r *elastic.SearchResult) (slice []string) {
	if r.Hits == nil || r.Hits.Hits == nil || len(r.Hits.Hits) == 0 {
		return nil
	}
	for _, hit := range r.Hits.Hits {
		var newLog LogMessage
		json.Unmarshal(hit.Source, &newLog)
		slice = append(slice, newLog.Log)
	}
	return slice
}
