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

const (
	DbUrl = "192.168.40.221:6379"
	DB    = 6
)

func QueryRedis2(client *redis.Client, id string, src, result sync.Map) (sucess bool) {

	var incr func(string) error

	// Transactionally increments key using GET and SET commands.
	incr = func(key string) error {
		err := client.Watch(func(tx *redis.Tx) error {
			//value, err := tx.Get(key).Int64()
			//value, err := tx.Get(key).Bytes()
			var value string
			err := tx.Set(key, value, time.Second*10).Err()
			if err != nil && err != redis.Nil {
				return err
			}
			fmt.Println("-----row:", value)

			//_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			//	pipe.Set(key, strconv.FormatInt(n+1, 10), 0)
			//	return nil
			//})

			if v, ok := src.Load(id); ok {
				//fmt.Println("key", int64(v))
				//if v, ok := src[key]; ok {

				//result[key] = (time.Now().UnixNano() / 1e6) - v
				result.Store(id, int64((time.Now().UnixNano()/1e6)-v.(int64)))
				//result.Store(id,int64((time.Now().UnixNano() / 1e6) - value))
				//result.Store(key,1)
				//sucess = true
				return nil
			}
			return err
		}, key)
		if err == redis.TxFailedErr {
			return incr(key)
		}
		return err
	}
	if incr(id) != nil {
		return true
	}
	return false

	//fmt.Println("key", val)
}
func QueryRedis(client *redis.Client, key string, src, result *sync.Map) (sucess bool) {
	//func QueryRedis(client *redis.Client,key string,src,result map[string]int64) (sucess bool) {
	//var err error
	if v, ok := src.Load(key); ok {
		//fmt.Println("sr key", v.(int64))
		//if v, ok := src[key]; ok {
		_, err := client.Get(key).Result()
		//fmt.Println(" --getfromredis:", row)
		if err != nil {
			//fmt.Println("----err:",err)
			//panic(err)
			return false
		}
		//result[key] = (time.Now().UnixNano() / 1e6) - v
		result.Store(key, int64((time.Now().UnixNano()/1e6)-v.(int64)))
		//result.Store(key,1)
		sucess = true
	}else{
		fmt.Println("not found in src map,key=", key)
		return false
	}

	return
	//fmt.Println("key", val)
}
