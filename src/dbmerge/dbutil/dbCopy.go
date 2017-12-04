package dbutil

import (
	"dbmerge/goSnowFlake"

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
	DbDestUrl  string
	IdGen  *goSnowFlake.IdWorker
	NodeId int64
}
const (
DbSrcUrl= "ms-test:ms-test@tcp(www.anycloud.top:13306)/music?charset=utf8&parseTime=True&loc=Local"
DbDestUrl= "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local"
)
