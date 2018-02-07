package dbQuery

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"sync"
	"time"
)

func QueryEs(esClient *elastic.Client, index string, mapType string, key string, result *sync.Map, timeStart int64) (sucess bool) {
	esService := esClient.Exists().Index(index).Type(mapType).Id(key)
	exists, err2 := esService.Do(context.TODO())
	if err2 != nil {
		panic(err2)
	}
	if !exists {
		//t.Fatal("expected document to exist")
		for i := 1; i < 10; i++ {
			time.Sleep(time.Duration(50*i) * time.Millisecond)
			exists, _ = esService.Do(context.TODO())
			if exists {
				break
			}

		}
	}
	if !exists {
		fmt.Println("expected document to exist:" + key)
		return false

	} else {
		result.Store(key, (time.Now().UnixNano())-timeStart)
		return true
	}
}
