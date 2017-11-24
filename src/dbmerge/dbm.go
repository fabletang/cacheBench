package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

func main() {
	db, err := gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	//db.Exec("DROP TABLE users;")
	var country Country
	//db.Raw("SELECT name, age FROM country WHERE name = ?", 3).Scan(&result)
	db.Raw("SELECT id, countryname,countrycode,status FROM country WHERE id = ?", 3).Scan(&country)
	fmt.Printf("country = %+v\n", country)
	//fmt.Printf("country = %#v\n", country.CountryCode)

	// Raw SQL
	rows, err := db.Raw("select id, countryname from country where id < ? order by id desc", 9).Rows() // (*sql.Rows, error)
	defer rows.Close()
	var id int64
	var countryname string
	for rows.Next() {
		rows.Scan(&id, &countryname)
		fmt.Printf("id = %+v\n", id)
		fmt.Printf("countryname = %+v\n", countryname)
	}
}
