package classLog

import (
	"dbmerge/dbutil"
	"dbmerge/goSnowFlake"
	"fmt"
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)
const(
Dest_t_drop string = "DROP TABLE IF EXISTS t_class_log_ref;"
Dest_t_new string = `
CREATE TABLE t_class_log_ref (
  id    bigint NOT NULL DEFAULT 0,
  oldId bigint NOT NULL DEFAULT 0,
  old_time_created varchar(30) NOT NULL DEFAULT '' COMMENT '添加时间',
  old_time_updated int(11) NOT NULL DEFAULT '0' COMMENT '时间',
  createTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  updateTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
  isDone boolean NOT NULL default false,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='排课 日志';
`
)
func Export(dbExport dbutil.DbExport, from int, len int, cnum chan int) {
	// Raw SQL
	//rows, err := dbCopy.DbSrc.Raw("select id, time_created, time_updated from user where role = ? order by id desc", 0).Rows() // (*sql.Rows, error)
	dbSrc, err := gorm.Open("mysql", dbExport.DbSrcUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbSrc.Close()
	dbDest, err := gorm.Open("mysql", dbExport.DbDestUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbDest.Close()
	//rows, err := dbSrc.Raw("select id, time_created, time_updated from user limit ?,?", from, len).Rows() // (*sql.Rows, error)
	rows, err := dbSrc.Raw("select id, time from class_log limit ?,?", from, len).Rows() // (*sql.Rows, error)
	if err != nil {
		log.Fatal(err)
	}
	//dbSrc.Close()
	//dbDest, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	//defer dbDest.Close()

	defer rows.Close()
	var id int64
	var time_created string
	var time_updated int64

	var createTime time.Time
	var updateTime time.Time
	var newId int64
	//var num int
	//获取本地location
	//toBeCharge := "2015-01-01 00:00:00"                             //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区

	fmt.Printf(" -- classLog start from:%+v  len: %+v \n ", from, len)
	for rows.Next() {
		//rows.Scan(&id, &time_created, &time_updated)
		rows.Scan(&id, &time_created)
		createTime, _ = time.ParseInLocation(timeLayout, time_created, loc) //使用模板在对应时区转化为time.time类型
		//createTime, _ = time.Parse(timeLayout, time_created) //使用模板在对应时区转化为time.time类型
		//fmt.Printf("id = %+v\n", id)
		if createTime.Second() == 0 {
			createTime = time.Now()
			newId, _ = dbExport.IdGen.NextId()
		} else {
			newId = goSnowFlake.JoinId(createTime.Unix(), dbExport.NodeId)
			//createTime = time.Unix(time_created, 0)
		}
		if time_updated == 0 {
			updateTime = time.Now()
		} else {
			updateTime = time.Unix(time_updated, 0)
		}
		dbDest.Exec("insert into t_class_log_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values (?,?,?,?,?,?)", newId, id, time_created, time_updated, createTime, updateTime)
		//dbDest.Exec("insert into t_customer_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values (?,?,?,?,?,?)", newId, id, time_created, time_updated, createTime, updateTime)
		//fmt.Printf("id = %+v\n", id)
		//fmt.Printf("countryname = %+v\n", countryname)
	}
	fmt.Printf(" classLog end from:%+v  len: %+v \n ", from, len)
	//fmt.Printf(" userOld rows = %+v\n", num)
	cnum <- 1
}
func CopyOld(dbCopy dbutil.DbCopy, from int, len int, cnum chan int) {
	// Raw SQL
	//rows, err := dbCopy.DbSrc.Raw("select id, time_created, time_updated from user where role = ? order by id desc", 0).Rows() // (*sql.Rows, error)
	//dbSrc, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/music?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer dbSrc.Close()
	//rows, err := dbSrc.Raw("select id, time_created, time_updated from user limit ?,?", from, len).Rows() // (*sql.Rows, error)
	rows, err := dbCopy.DbSrc.Raw("select id, time from class_log limit ?,?", from, len).Rows() // (*sql.Rows, error)
	if err != nil {
		log.Fatal(err)
	}
	//dbDest, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	//defer dbDest.Close()

	defer rows.Close()
	var id int64
	var time_created string
	var time_updated int64

	var createTime time.Time
	var updateTime time.Time
	var newId int64
	//var num int
	//获取本地location
	//toBeCharge := "2015-01-01 00:00:00"                             //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区

	fmt.Printf(" -- classLog start from:%+v  len: %+v \n ", from, len)
	for rows.Next() {
		//rows.Scan(&id, &time_created, &time_updated)
		rows.Scan(&id, &time_created)
		createTime, _ = time.ParseInLocation(timeLayout, time_created, loc) //使用模板在对应时区转化为time.time类型
		//createTime, _ = time.Parse(timeLayout, time_created) //使用模板在对应时区转化为time.time类型
		//fmt.Printf("id = %+v\n", id)
		if createTime.Second() == 0 {
			createTime = time.Now()
			newId, _ = dbCopy.IdGen.NextId()
		} else {
			newId = goSnowFlake.JoinId(createTime.Unix(), dbCopy.NodeId)
			//createTime = time.Unix(time_created, 0)
		}
		if time_updated == 0 {
			updateTime = time.Now()
		} else {
			updateTime = time.Unix(time_updated, 0)
		}
		dbCopy.DbDest.Exec("insert into t_class_log_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values (?,?,?,?,?,?)", newId, id, time_created, time_updated, createTime, updateTime)
		//dbDest.Exec("insert into t_customer_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values (?,?,?,?,?,?)", newId, id, time_created, time_updated, createTime, updateTime)
		//fmt.Printf("id = %+v\n", id)
		//fmt.Printf("countryname = %+v\n", countryname)
	}
	fmt.Printf(" classLog end from:%+v  len: %+v \n ", from, len)
	//fmt.Printf(" userOld rows = %+v\n", num)
	cnum <- 1
}
