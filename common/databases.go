package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB;

//打开连接池
func InitDB() *gorm.DB{
	driverName := "mysql"
	host := "localhost"
	port := "3306";
	database := "testlogin"
	username := "root"
	password := "root"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db,err := gorm.Open(driverName,args)
	if err != nil {
		panic("failed to open database, err: "+err.Error())
	}

	return db;
}

func GetDB() *gorm.DB {
	return DB;
}
