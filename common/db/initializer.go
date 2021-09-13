package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/poniteru/go-coin-watcher/app/config"
)

var db *sqlx.DB

func init() {
	initDB()
}

func initDB() {
	dsn := config.MysqlDsn
	// 也可以使用MustConnect连接不成功就panic
	//db, err = sqlx.Connect("mysql", dsn)
	//if err != nil {
	//	fmt.Printf("connect DB failed, err:%v\n", err)
	//	return
	//}
	db = sqlx.MustConnect("mysql", dsn)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
}

func Instance() *sqlx.DB {
	return db
}
