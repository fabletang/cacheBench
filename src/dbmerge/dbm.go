package main

import (
	"dbmerge/dbutil"
	"dbmerge/goSnowFlake"
	//"dbmerge/logic/userOld"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"runtime"
	"dbmerge/logic/classLog"
	"dbmerge/logic/userOld"
)

type Country struct {
	//gorm.Model
	Id          int
	CountryCode string //`gorm:"type:varchar(100);unique"`
	CountryName string //`gorm:"size:255"`
	Status      int
	//CountryName2 string
	//Name        string `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	//Num         int    `gorm:"AUTO_INCREMENT"`
	//Birthday time.Time
	//CreditCard        CreditCard      // One-To-One relationship (has one - use CreditCard's UserID as foreign key)
	//Emails            []Email         // One-To-Many relationship (has many - use Email's UserID as foreign key)
}

//type DbCopy struct {
//	DbSrc  *gorm.DB
//	DbDest *gorm.DB
//	IdGen  *goSnowFlake.IdWorker
//	NodeId int64
//}

var dest_t_drop string = "DROP TABLE IF EXISTS t_customer_ref;"
var dest_t_new string = `
CREATE TABLE t_customer_ref (
  id    bigint NOT NULL DEFAULT 0,
  oldId bigint NOT NULL DEFAULT 0,
  old_time_created int(11) NOT NULL DEFAULT '0' COMMENT '添加时间',
  old_time_updated int(11) NOT NULL DEFAULT '0' COMMENT '时间',
  createTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  updateTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
  isDone boolean NOT NULL default false,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='国家信息';
`
var doStart , doEnd time.Time
func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println(terminal)
}
//var wg sync.WaitGroup  //定义一个同步等待的组
var cnum chan int



func main() {
	maxProcs := runtime.NumCPU() // 获取cpu个数
	runtime.GOMAXPROCS(maxProcs) //限制同时运行的goroutines数量
	fmt.Printf("maxProcs = %+v\n", maxProcs)

	defer timeCost(time.Now())
	//db, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/music?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open("mysql", "ms-test:ms-test@tcp(223.167.128.39:13306)/music?charset=utf8&parseTime=True&loc=Local")

    //db.DB().SetMaxOpenConns(20)
	//db.DB().SetMaxOpenConns(dbutil.ConnMax/2)
    //db.DB().SetMaxIdleConns(dbutil.ConnMax/8)
	//db, err := gorm.Open("mysql", "viptest:viptest_2017@tcp(192.168.40.219:3306)/music?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	//db.Exec("set GLOBAL max_connections=1000;")
	//db_dest, err := gorm.Open("mysql", "ms-test:ms-test@tcp(192.168.40.203:13306)/test-music?charset=utf8&parseTime=True&loc=Local")
	//db_dest, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	db_dest, err := gorm.Open("mysql", "ms-test:ms-test@tcp(223.167.128.39:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	//db_dest.DB().SetMaxOpenConns(dbutil.ConnMax/2)
	//db_dest.DB().SetMaxIdleConns(dbutil.ConnMax/8)
	defer db_dest.Close()
	if err != nil {
		log.Fatal(err)
	}

	db_dest.Exec(dest_t_drop)
	db_dest.Exec(dest_t_new)
	//var country Country
	//db.Raw("SELECT name, age FROM country WHERE name = ?", 3).Scan(&result)
	//db.Raw("SELECT id,countrycode, countryname,status FROM country WHERE id = ?", 3).Scan(&country)
	//fmt.Printf("country = %+v\n", country)
	//fmt.Printf("country = %#v\n", country.CountryCode)
	var workId int64 = 123
	iw, _ := goSnowFlake.NewIdWorker(workId)

	var dbCopy dbutil.DbCopy
	dbCopy.DbSrc = db
	dbCopy.DbDest = db_dest
	dbCopy.IdGen = iw
	dbCopy.NodeId = workId

	var dbExport dbutil.DbExport
	dbExport.DbSrcUrl = dbutil.DbSrcUrl
	dbExport.DbDestUrl = dbutil.DbDestUrl
	dbExport.IdGen = iw
	dbExport.NodeId = workId
	//userOld.CopyOld(dbCopy, 10, 100)
	//records := dbuitl.GetCount(db, "user")
	records := dbutil.GetCount(db, "user")
	//records := 308
	//pages := records / perPage
	//left := records % perPage
	//
	//if left != 0 {
	//	pages += 1
	//}
	pages, perPage := dbutil.GetConnNum(records)
	cnum = make(chan int, pages) //make一个chan,缓存为num
	fmt.Printf("pages = %+v\n", pages)
for i := 0; i < pages; i++ {
		go userOld.CopyOld(dbCopy, i*perPage, perPage, cnum)
	}
	for i := 0; i < pages; i++ {
		<-cnum
	}
	//wg.Wait() //阻塞等待所有组内成员都执行完毕退栈
	//time.Sleep(time.Second*100)
	fmt.Println(" userOld end-----")

	db_dest.Exec(classLog.Dest_t_drop)
	db_dest.Exec(classLog.Dest_t_new)
    doStart=time.Now()
	records = dbutil.GetCount(db, "class_log")
	//records=125
	pages, perPage = dbutil.GetConnNum(records)
	cnum = make(chan int, pages) //make一个chan,缓存为num
	//pages=1
	fmt.Printf("class_log pages = %+v\n", pages)
	for i := 0; i < pages; i++ {
		//go classLog.CopyOld(dbCopy, i*perPage, perPage, cnum)
		go classLog.Export(dbExport, i*perPage, perPage, cnum)
	}
	for i := 0; i < pages; i++ {
		<-cnum
	}
	fmt.Println(" classLog end-----")
	//doEnd=time.Now()
	secondDru:=time.Since(doStart)
	fmt.Printf("classLog records:%v spend time: %v ,ops: %v/second \n",records,secondDru,records/(int)(secondDru.Seconds()))
}
