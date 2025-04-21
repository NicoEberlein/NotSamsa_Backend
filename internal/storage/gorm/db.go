package gormstore

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func GetDbConnection() *gorm.DB {

	dsn := "host=localhost user=notsamsa password=notsamsapw dbname=notsamsa port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db
}
