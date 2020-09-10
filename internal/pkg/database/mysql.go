package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"

	// 引用数据库驱动初始化
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	maxOpenConns    = 150
	maxIdleConns    = 100
	connMaxLifetime = 100
)

var engine *gorm.DB
var dbName string = ""

// GetDB get gorm.DB
func GetDB() *gorm.DB {
	return engine
}

// GetDbName get database name
func GetDbName() string {
	return dbName
}

func SetDbName(name string) {
	dbName = name
}
// Close closes current db connection
func Close() error {
	return engine.Close()
}

// SetUp set up db connection
func SetUp(user, pwd, host, port string) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		user,
		pwd,
		host,
		port,
		"information_schema",
		"utf8mb4",
		true,
		"Local")

	engine, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Println("connect to mysql fail, ", dsn, err)
		panic(err)
	}
	engine.LogMode(false)

	engine.DB().SetConnMaxLifetime(connMaxLifetime * time.Second)
	engine.DB().SetMaxOpenConns(maxOpenConns)
	engine.DB().SetMaxIdleConns(maxIdleConns)
}
