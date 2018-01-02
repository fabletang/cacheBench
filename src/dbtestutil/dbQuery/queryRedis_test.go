package dbQuery

import (
	//"math/rand"
	//"strconv"
	"fmt"
	"github.com/go-redis/redis"

	"testing"
	"time"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "192.168.40.221:6379",
		Password: "", // no password set
		DB:       6,  // use default DB
	})
	client.Set("music-testCanal-123","123456",time.Second*60)
}
func TestQueryRedis(t *testing.T) {
	//func QueryRedis(client *redis.Client,key string,src,result map[string]int64) (sucess bool) {
	//var err error

	row, err := client.Get("music-testCanal-123").Result()
	if err != nil {
		//fmt.Println("----err:",err)
		//panic(err)
		t.Errorf("expecting connNum > %s, got %s", row, err)
		//return false
	}
	fmt.Println("row", row)
	//result[key] = (time.Now().UnixNano() / 1e6) - v
	//result.Store(key,int64((time.Now().UnixNano() / 1e6) - v.(int64)))
	//result.Store(key,1)
	//sucess = true

	//return
	//fmt.Println("key", val)
}
