package main

import (
	"dbtestutil/dbModify"
	"dbtestutil/dbQuery"
	"dbtestutil/goSnowFlake"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	//"gopkg.in/cheggaaa/pb.v1"
	"gopkg.in/cheggaaa/pb.v1"
	//esClient "github.com/elastic/go-elasticsearch/client"
	//elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/olivere/elastic"
)

var Config = struct {
	APPName string `default:"cache bench"`

	Mysql struct {
		User     string `required:"true" default:"ms-test"`
		Password string `required:"true" env:"DBPassword"`
		Host     string `default:"ms-test"`
		Port     int    `default:"13306"`
		Schema   string `default:"music"`
		Table    string `default:"testCache"`
	}

	Redis struct {
		Host string `default:"localhost"`
		Port int    `default:"6379"`
		Db   int    `default:"6"`
	}

	Bench struct {
		MaxRecordNum  int `default:"100"`
		ThreadNum     int `default:"10"`
		MysqlInterval int `default:"1"`
		RedisInterval int `default:"5"`
	}
	Es struct {
		Url string `default:"localhost"`
	}
}{}
var N int = 10
var R int = 2

//var sumDb int32 = 0
//var sumCache int32 = 0

var client *redis.Client
var resultTime sync.Map
var esTime sync.Map
var dbMysql *gorm.DB
var schema, table string

//var mysqlConf *Config
func InsertDbAndFindRedis(id int64) {

	timeStart := time.Now().UnixNano()
	dbMysql = dbModify.ModifyDb(dbMysql, id)
	if dbMysql.RowsAffected == 1 {
		//atomic.AddInt32(&sumDb, 1)
		time.Sleep(time.Duration(50) * time.Millisecond)
		keyPrefix := schema + "-" + table + "-"
		key := keyPrefix + strconv.FormatInt(id, 10)

		dbQuery.FindRedis(client, key, &resultTime, timeStart)

		dbQuery.QueryEs(esClient, "qa-test_canal", "test", key, &esTime, timeStart)
	}
}
func Workers(task func(interface{}), climax func()) chan interface{} {
	input := make(chan interface{})
	ack := make(chan bool)
	for i := 0; i < R; i++ {
		go func() {
			for {
				v, ok := <-input
				if ok {
					task(v)
					ack <- true
				} else {
					return
				}
			}
		}()
	}
	go func() {
		for i := 0; i < N; i++ {
			<-ack
			bar.Increment()
		}
		climax()
	}()
	return input
}

var bar *pb.ProgressBar
var esClient *elastic.Client

func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println("total cost time:", terminal)
}
func main() {
	defer timeCost(time.Now())
	configor.Load(&Config, "config.yml")
	benchConf := Config.Bench
	N = benchConf.MaxRecordNum // 表示插入数据数量
	R = benchConf.ThreadNum    // mysql redis 线程
	if R > N {
		R = N
	}
	//fmt.Println(" N:", N, " R:", R)
	//fmt.Printf("config: %#v", Config)
	runtime.GOMAXPROCS(R * 2) //限制同时运行的goroutines数量

	var err1 error
	esClient, err1 = elastic.NewClient(elastic.SetURL(Config.Es.Url))
	if err1 != nil {
		panic(err1)
	}
	defer esClient.Stop()

	var err error
	mysqlConf := Config.Mysql
	mysqlUrl := mysqlConf.User + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + strconv.Itoa(mysqlConf.Port) + ")/" + mysqlConf.Schema
	mysqlUrl += "?charset=utf8&parseTime=True&loc=Local"
	fmt.Println("---mysqlUrl:" + mysqlUrl)
	schema = Config.Mysql.Schema
	table = Config.Mysql.Table
	dbMysql, err = gorm.Open("mysql", mysqlUrl)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer dbMysql.Close()

	rows, err := dbMysql.Raw("select id from " + table + " limit 1").Rows() // (*sql.Rows, error)
	if err != nil {
		fmt.Println(" table 不存在:", table)
		dbMysql.Exec(dbModify.T_new)
		fmt.Println(" 创建 table:", table)
	}
	defer rows.Close()

	redisConf := Config.Redis
	client = redis.NewClient(&redis.Options{
		Addr:     redisConf.Host + ":" + strconv.Itoa(redisConf.Port),
		Password: "",           // no password set
		DB:       redisConf.Db, // use default DB
		PoolSize: Config.Bench.ThreadNum,
	})
	fmt.Println("---redisConf:", redisConf)
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer client.Close()

	idGen, _ := goSnowFlake.NewIdWorker(123)
	//----------------------------
	bar = pb.StartNew(Config.Bench.MaxRecordNum)

	exit := make(chan bool)

	workers := Workers(func(a interface{}) {
		//Task(a.(int))
		InsertDbAndFindRedis(a.(int64))
	}, func() {
		exit <- true
	})

	for i := 0; i < N; i++ {
		//workers <- i
		id, _ := idGen.NextId()
		workers <- id
	}
	close(workers)

	<-exit

	costMax, costMin, costAve, totalNum := stats2(&resultTime)
	fmt.Println("----redis查询 有效记录数:", totalNum,"超时失败数:",int64(N)-totalNum, " 耗费毫秒数---最大=", costMax, "最小=", costMin, "平均=", costAve)
	costMax, costMin, costAve, totalNum = stats2(&esTime)
	fmt.Println("----es查询    有效记录数:", totalNum, "超时失败数:",int64(N)-totalNum, " 耗费毫秒数---最大=", costMax, "最小=", costMin, "平均=", costAve)
	//bar.Increment()
}
func stats2(dest *sync.Map) (costMax, costMin, costAve, totalNum int64) {
	//func stats(result map[string]int64) (costMax, costMin, costAve, totalNum int64) {
	costMax, costMin = 1, int64(^uint(0)>>1)
	var total int64

	var tmp int64
	dest.Range(func(k, v interface{}) bool {
		tmp = v.(int64)
		tmp = tmp
		//fmt.Println(k, v)
		if tmp > costMax {
			costMax = tmp
		}
		if tmp < costMin {
			costMin = tmp
		}
		total += tmp
		totalNum += 1
		return true
	})
	if totalNum == 0 {
		return
	}
	costAve = total / totalNum

	costMax = costMax / 1e6
	costMin = costMin / 1e6
	costAve = costAve / 1e6
	return
}
