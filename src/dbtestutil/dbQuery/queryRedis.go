package dbQuery

import (
	//"math/rand"
	//"strconv"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
	//"strconv"
)

func QueryRedis(client *redis.Client, key string, src, result *sync.Map) (sucess bool) {
	if v, ok := src.Load(key); ok {
		//if v, ok := src[key]; ok {
		_, err := client.Get(key).Result()
		//fmt.Println(" --getfromredis:", row)
		if err != nil {
			//fmt.Println("----err:",err)
			//panic(err)
			return false
		}
		//result[key] = (time.Now().UnixNano() / 1e6) - v
		result.Store(key, (time.Now().UnixNano())-v.(int64))
		//result.Store(key,1)
		sucess = true
	}else{
		fmt.Println("not found in src map,key=", key)
		return false
	}
	return
}
