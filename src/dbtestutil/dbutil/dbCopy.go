package dbutil

import (
	"dbtestutil/goSnowFlake"

	"github.com/jinzhu/gorm"
)

type DbCopy struct {
	DbSrc  *gorm.DB
	DbDest *gorm.DB
	IdGen  *goSnowFlake.IdWorker
	NodeId int64
}

type DbExport struct {
	DbSrcUrl  string
	DbDestUrl string
	IdGen     *goSnowFlake.IdWorker
	NodeId    int64
}

const (
	//DbSrcUrl= "ms-test:ms-test@tcp(www.anycloud.top:13306)/music?charset=utf8&parseTime=True&loc=Local"
	//DbDestUrl= "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local"
	//db, err := gorm.Open("mysql", "ms-test:ms-test@tcp(223.167.128.39:13306)/music?charset=utf8&parseTime=True&loc=Local")
	DbSrcUrl  = "viptest:viptest_2017@tcp(192.168.40.219:3306)/music?charset=utf8&parseTime=True&loc=Local"
	DbDestUrl = "ms-test:ms-test@tcp(192.168.40.203:13306)/test-music?charset=utf8&parseTime=True&loc=Local"
)
