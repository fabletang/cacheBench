package dbutil

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GetCount(db *gorm.DB, tableName string) (nums int) {
	//db.Raw("select count(1) as nums from " + tableName).Scan(&nums)
	//nums = count.nums
	//defer db.Close()
	rows, err := db.Raw("select count(1) as nums  from " + tableName).Rows() // (*sql.Rows, error)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&nums)
	}
	//fmt.Printf("nums = %+v\n", nums)
	return
}
