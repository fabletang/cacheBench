package dbModify

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"strings"
	//"context"
	"github.com/jinzhu/gorm"
	"math/rand"
	//"strconv"
)

const (
	T_drop = "DROP TABLE IF EXISTS test_canal;"
	T_new  = `
		CREATE TABLE test_canal (
  		id    bigint NOT NULL DEFAULT 0,
		name  varchar(1000) NOT NULL DEFAULT 'name-fdasdlkklsdfkldlafdlsakfdlasfdaslfdsalfdalskfldasklfldslafldaslkfl',
		remark  varchar(1000) NOT NULL DEFAULT 'remark-aaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbcccccccccccccccccccc',
		demo  varchar(1000) NOT NULL DEFAULT 'demo-GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG',
  		status bigint NOT NULL DEFAULT 0,
  		birth int(11) NOT NULL DEFAULT '0' COMMENT '生日',
		score  int NOT NULL DEFAULT 0,
  		registerDate int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
  		createTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  		updateTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
  		isDone boolean NOT NULL default false,
  		PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='cache性能测试';
`
)

func ModifyDb(db *gorm.DB, id int64) {
	db.Exec("insert into test_canal (id,score) values (?,?)", id, rand.Intn(10000))

	//if (rand.Intn(100)>30){
	//db.Exec("update testCanal set name= ? where id= ?","abc"+strconv.Itoa(rand.Intn(10000)),id)
	//}
}
