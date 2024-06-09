package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	err      error
)

func OpenConnection() {
	dsn := "host=localhost user=postgres password=admin dbname=temukaDB port=5432 sslmode=disable"
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetDBInstance() *gorm.DB {
	return Database
}
