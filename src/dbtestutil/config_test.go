package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {

	configor.Load(&Config, "config.yml")
	fmt.Printf("config: %#v\n", Config)
	fmt.Printf("Config.Mysql.Port: %d\n", Config.Mysql.Port)
	fmt.Printf("Config.Mysql.Pwd: %s\n", Config.Mysql.Password)

	mysqlConf := Config.Mysql
	mysqlUrl := mysqlConf.User + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + strconv.Itoa(mysqlConf.Port) + ")/" + mysqlConf.Schema
	mysqlUrl += "?charset=utf8&parseTime=True&loc=Local"
	fmt.Println("---mysqlUrl:" + mysqlUrl)

}
