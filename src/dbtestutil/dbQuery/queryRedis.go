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
	} else {
		fmt.Println("not found in src map,key=", key)
		return false
	}
	return
}

func FindRedis(client *redis.Client, key string, result *sync.Map, timeStart int64) (sucess bool) {
	_,err := client.Get(key).Result()
	//c, err := client.Exists(key).Result()
	//rs:=client.Exists(key)
	//fmt.Println(" --getfromredis:", row)

	//fmt.Println(" key:",key," err:",err)
	//if len(c) < 2 {
	if err != nil {
		//if len(c) < 2 {
		for i := 0; i < 9; i++ {
			time.Sleep(time.Duration(50) * time.Millisecond)
			//c, err = client.Exists(key).Result()
			_, err = client.Get(key).Result()
			//fmt.Println(" key:",key," err:",err)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return false
	}
	//if len(c) < 2 {
	//	return false
	//}
	//fmt.Println(" key",key," exist:",c)
	result.Store(key, (time.Now().UnixNano())-timeStart)
	return true
}
