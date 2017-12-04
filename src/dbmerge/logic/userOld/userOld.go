package userOld

import (
	"dbmerge/dbutil"
	"dbmerge/goSnowFlake"
	"fmt"
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CopyOld(dbCopy dbutil.DbCopy, from int, len int, cnum chan int) {
	// Raw SQL
	//rows, err := dbCopy.DbSrc.Raw("select id, time_created, time_updated from user where role = ? order by id desc", 0).Rows() // (*sql.Rows, error)
	//dbSrc, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/music?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer dbSrc.Close()
	//rows, err := dbSrc.Raw("select id, time_created, time_updated from user limit ?,?", from, len).Rows() // (*sql.Rows, error)
	rows, err := dbCopy.DbSrc.Raw("select id, time_created, time_updated from user limit ?,?", from, len).Rows() // (*sql.Rows, error)
	if err != nil {
		log.Fatal(err)
	}
	//dbDest, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	//defer dbDest.Close()

	defer rows.Close()
	var id int64
	var time_created int64
	var time_updated int64

	var createTime time.Time
	var updateTime time.Time
	var newId int64
	//var num int
	for rows.Next() {
		//for 1==2 {
		//num += 1
		rows.Scan(&id, &time_created, &time_updated)
		//fmt.Printf("id = %+v\n", id)
		if time_created == 0 {
			createTime = time.Now()
			newId, _ = dbCopy.IdGen.NextId()
		} else {
			newId = goSnowFlake.JoinId(time_created, dbCopy.NodeId)
			createTime = time.Unix(time_created, 0)
		}
		if time_updated == 0 {
			updateTime = time.Now()
		} else {
			updateTime = time.Unix(time_updated, 0)
		}
		dbCopy.DbDest.Exec("insert into t_customer_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values (?,?,?,?,?,?)", newId, id, time_created, time_updated, createTime, updateTime)
		//dbDest.Exec("insert into t_customer_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values (?,?,?,?,?,?)", newId, id, time_created, time_updated, createTime, updateTime)
		//fmt.Printf("id = %+v\n", id)
		//fmt.Printf("countryname = %+v\n", countryname)
	}
	fmt.Printf(" userOld from:%+v  len: %+v \n ", from, len)
	//fmt.Printf(" userOld rows = %+v\n", num)
	cnum <- 1
}
