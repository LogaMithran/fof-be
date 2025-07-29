package connectors

import (
	"friends-of-friends-be/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var (
	Db     *gorm.DB
	err    error
	entity = []interface{}{entities.User{}}
)

func InitializeDbConnection() {
	dsn := os.Getenv("MYSQL_DSN")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error in establishing database connection %v", err)
	}

	if migrateErr := Db.AutoMigrate(entity...); migrateErr != nil {
		panic(migrateErr)
	}

	log.Print("Database connection established successfully")
}
