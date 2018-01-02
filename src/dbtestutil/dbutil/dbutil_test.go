package dbutil

import (
	"fmt"
	"log"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "ms-test:ms-test@tcp(www.anycloud.top:13306)/ms-test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetCount(t *testing.T) {
	result := GetCount(db, "user")
	fmt.Printf("result = %+v\n", result)
	if result == 0 {
		t.Errorf("expecting >0, got %s", result)
	}
}
