package gormstore

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func GetDbConnection() *gorm.DB {

	var db *gorm.DB
	var err error

	retryCounter := 0

	for retryCounter < 3 {
		dsn := "host=db user=notsamsa password=notsamsapw dbname=notsamsa port=5432 sslmode=disable"
		fmt.Println("Attempting to connect to database. Retrycounter: ", retryCounter)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		} else {
			retryCounter++
			if retryCounter == 3 {
				panic(fmt.Errorf("failed to connect to database %s", err))
			}
			time.Sleep(8 * time.Second)
		}
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db
}
