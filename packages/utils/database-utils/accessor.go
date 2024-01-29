package database_utils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
	configurationmanager "update-microservice/packages/utils/configuration-manager"
)

var db *gorm.DB

func OpenDatabaseConnection() (*gorm.DB, error) {
	dsn := configurationmanager.Configuration.DatabaseConnectionString
	if db != nil {
		return db, nil
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Second * 5)
	sqlDB.SetConnMaxIdleTime(time.Second * 5)
	return db, nil
}
