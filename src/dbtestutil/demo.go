package main

import (
	//"dbtestutil/goSnowFlake"
	"fmt"
	"time"

	"context"
	"dbtestutil/goSnowFlake"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"math/rand"
	"runtime"
	//"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"github.com/jinzhu/gorm"
	//"dbtestutil/dbutil"
	"log"
	"os"
	"dbtestutil/dbModify"
	//"dbtestutil/logic/userOld"
	 "github.com/go-redis/redis"
	//"dbtestutil/dbQuery"
)

//var doStart time.Time

const (
	//MaxRecordNum=10000*5
	MaxRecordNum  = 10
	TestDB        = "music"
	TestTableName = "testCanal"
	IdPrefix      = TestDB + "-" + TestTableName + "-"
	ModifyDbNum   = 10
	QueryCacheNum = 1000
)

func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println("total cost time:", terminal)
}

var wg sync.WaitGroup //定义一个同步等待的组
//var endGoroutine chan int

var sumDb int32 = 0
var sumCache int32 = 0
var idGen *goSnowFlake.IdWorker

//var cost [MaxRecordNum]int64
//var m map[string]int64
var m sync.Map
//var result map[string]int64
var result sync.Map

var db  *gorm.DB
//iw, _ := goSnowFlake.NewIdWorker(workId)
func main() {
	//costCache = []int64{}
	//costDb=[]int64{}
	//m = make(map[string]int64)
	//result = make(map[string]int64)
	defer timeCost(time.Now())
	//startTime:=time.Now()
	maxProcs := runtime.NumCPU()     // 获取cpu个数
	runtime.GOMAXPROCS(maxProcs * 8) //限制同时运行的goroutines数量
	fmt.Printf("maxProcs = %+v\n", maxProcs*8)
   var err error
	db, err = gorm.Open("mysql", dbModify.DbUrl)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.40.221:6379",
		Password: "", // no password set
		DB:       6,  // use default DB
	})
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer client.Close()

	db.Exec(dbModify.T_drop)
	db.Exec(dbModify.T_new)

	var workId int64 = 123
	idGen, _ = goSnowFlake.NewIdWorker(workId)
	//iw.NextId();
	//timeDelta,_:=idGen.NextId()
	//queue := make(chan int, 10)  // 这里的10表示管道的容量，根据应用的需求进行设置

	dbIdqueue := make(chan int64, ModifyDbNum)       // 这里的10表示管道的容量，根据应用的需求进行设置
	//cacheIdqueue := make(chan string, QueryCacheNum) // 这里的10表示管道的容量，根据应用的需求进行设置
	cacheIdqueue := make(chan string, QueryCacheNum) // 这里的10表示管道的容量，根据应用的需求进行设置
	wg.Add(1)
	background := context.Background()
	//ctx, _ := context.WithTimeout(background, 2*time.Second)
	ctxDb, cancelDb := context.WithCancel(background)

	//background2 := context.Background()
	//ctxCache, cancelCache := context.WithCancel(background)
	//pubsub := client.PSubscribe("test")

	pubsub := client.Subscribe("test")

	defer pubsub.Close()

	msg, err := pubsub.ReceiveMessage()
	fmt.Println("---sub:",msg.Payload)

	//fmt.Println("-msgi:",pubsub.r)
	//msgi, err := pubsub.ReceiveTimeout(time.Second*3)
	//subscr := msgi.(*redis.Subscription)
	////subscr := msgi.(*redis.Message)
	//
	////subscr.Count
	////fmt.Println("-msgi:",subscr.Payload)
	////fmt.Println("-msgi count:",subscr.Payload)
	//fmt.Println("-msgi count:",subscr.Channel)


	for i := 0; i < MaxRecordNum; i++ {
		//go QueryCache(client,cacheIdqueue, ctxCache,m,result)
	}

	for i := 0; i < ModifyDbNum; i++ {
		go ModifyDb(dbIdqueue, ctxDb, cacheIdqueue)
	}

	go P(dbIdqueue)


	//time.Sleep(time.Second*1)
	wg.Wait()

	//for {
	//	if sumDb == MaxRecordNum {
	//		cancelDb()
	//		break
	//	}
	//}
	for {
		if sumCache == MaxRecordNum {
			//cancelCache()
			//time.Sleep(time.Second*1)
			cancelDb()
			break
		}
	}
	time.Sleep(time.Second * 1)
	fmt.Println("sumDb=", sumDb, " sumCache=", sumCache)
	costMax, costMin, costAve, totalNum := stats(result)
	fmt.Println("共统计记录数:", totalNum, "----ms---costMax=", costMax, "costMin=", costMin, "costAve=", costAve)


	//time.Sleep(2 * time.Second)
	//cancel();

	//fmt.Println("total cost time:",time.Since(startTime) )
	//terminal := time.Since(start)
}
func stats(result sync.Map) (costMax, costMin, costAve, totalNum int64) {
	//func stats(result map[string]int64) (costMax, costMin, costAve, totalNum int64) {

	costMax, costMin = 0, int64(^uint(0)>>1)
	var total int64

	//for k, v := range result.(int64) {
	//	//for k, v := range result.(interface{}).(int64) {
	//	//for k, v := range result.(interface{}).(map[string]int64) {
	//	fmt.Println(k, v)
	//	//fmt.Println("name:", v.Name)
	//}

	var tmp int64
	m.Range(func(k, v interface{}) bool {
	tmp=v.(int64)
		//fmt.Println(k, v)
			if tmp > costMax {
				costMax = tmp
			}
			if tmp < costMin {
				costMin = tmp
			}
			total += tmp
			totalNum += 1
			return true;
	})
	// 遍历map
	//for _, v := range result {
	//	if v > costMax {
	//		costMax = v
	//	}
	//	if v < costMin {
	//		costMin = v
	//	}
	//	total += v
	//	totalNum += 1
	//}
	costAve = total / totalNum
	return
}

func P(dbIdQueue chan<- int64) {
	defer wg.Done()
	var id int64
	for i := 0; i < MaxRecordNum; i++ {
		//dbIdQueue <- (int64)(i)
		id, _ = idGen.NextId()
		dbIdQueue <- id
		//cacheIdQueue <- IdPrefix+string(id)
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		//cacheIdQueue <- IdPrefix + strconv.FormatInt(id, 10)
	}
}
func ModifyDb(dbIdQueue <-chan int64, ctx context.Context, cacheIdQueue chan string) {
	//defer wg.Done()
	var key string
	for {
		select {
		case i := <-dbIdQueue:
			//randNum := rand.Intn(100)
			//time.Sleep(time.Millisecond * time.Duration(randNum))
			//fmt.Println("--receive:", i, "sleep ms:", randNum)
			//costDb[sumDb] = time.Now().UnixNano() / 1e6
			dbModify.ModifyDb(db,i)
			key = IdPrefix + strconv.FormatInt(i, 10)
			//m[key] = time.Now().UnixNano() / 1e6
			m.Store(key,(time.Now().UnixNano() / 1e6))
			cacheIdQueue <- key
			atomic.AddInt32(&sumDb, 1)
			//time.Sleep(time.Microsecond*2)
			fmt.Println("--sumDb:", sumDb)
			//case wg.Done:
			//	fmt.Println("--receive end:")
		case <-ctx.Done():
			fmt.Println("------------ goroutine ModifyDb done")
			return
		}
	}
}

func QueryCache(client *redis.Client,cacheIdQueue <-chan string, ctx context.Context,src,result sync.Map) {
	//func QueryCache(client *redis.Client,cacheIdQueue <-chan string, ctx context.Context,src,result map[string]int64) {
	//defer wg.Done()
	for {
		select {
		case key := <-cacheIdQueue:
			//if v, ok := src[key]; ok {
			//	result[key] = (time.Now().UnixNano() / 1e6) - v
			//	//fmt.Println(v)
			//} else {
				fmt.Println("Key Not Found:", key)
			//}
			//dbQuery.QueryRedis(client, key, src, result)

			//for {
			//  success:=	dbQuery.QueryRedis(client, key, src, result)
			//  if success {
			//  	break;
			//  }
			//}
			atomic.AddInt32(&sumCache, 1)
			fmt.Println("--sumCache:", sumCache)
		case <-ctx.Done():
			fmt.Println("------------ goroutine QueryCache done")
			return
		}
	}

	//wg.Done()
}
