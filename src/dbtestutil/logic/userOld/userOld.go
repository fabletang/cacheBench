package userOld

import (
	"dbtestutil/dbutil"
	"dbtestutil/goSnowFlake"
	"fmt"
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

const (
	Dest_t_drop = "DROP TABLE IF EXISTS t_customer_ref;"
	Dest_t_new  = `
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
)

func CopyOld(dbCopy dbutil.DbCopy, from int, len int, cnum chan int) {
	rows, err := dbCopy.DbSrc.Raw("select id, time_created, time_updated from user limit ?,?", from, len).Rows() // (*sql.Rows, error)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var id int64
	var time_created int64
	var time_updated int64

	var createTime time.Time
	var updateTime time.Time
	var newId int64
	for rows.Next() {
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
		//fmt.Printf("id = %+v\n", id)
		//fmt.Printf("countryname = %+v\n", countryname)
	}
	fmt.Printf(" userOld from:%+v  len: %+v \n ", from, len)
	//fmt.Printf(" userOld rows = %+v\n", num)
	cnum <- 1
}

func BatchCopyOld(dbCopy dbutil.DbCopy, from int, len int, cnum chan int) {
	rows, err := dbCopy.DbSrc.Raw("select id, time_created, time_updated from user limit ?,?", from, len).Rows() // (*sql.Rows, error)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var id int64
	var time_created int64
	var time_updated int64

	var createTime time.Time
	var updateTime time.Time
	var newId int64
	sql := "insert into t_customer_ref(id,oldId,old_time_created,old_time_updated,createtime,updatetime) values "
	for rows.Next() {
		rows.Scan(&id, &time_created, &time_updated)
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
		sql += fmt.Sprintf("(%v,%v,%v,%v,'%v','%v'),", newId, id, time_created, time_updated, createTime.Format("2006-01-02 15:04:05"), updateTime.Format("2006-01-02 15:04:05"))
	}
	sql = strings.TrimSuffix(sql, ",")
	dbCopy.DbDest.Exec(sql)
	fmt.Printf(" userOld from:%+v  len: %+v \n ", from, len)
	cnum <- 1
}
