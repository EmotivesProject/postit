package db

import (
	"fmt"
	"os"
	"postit/model"

	myLogger "github.com/TomBowyerResearchProject/common/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	database *gorm.DB
)

//ConnectDB function: Make database connection
func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		myLogger.Fatal(err)
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	//Define DB connection string and connect
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		myLogger.Error(err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Like{},
		&model.Comment{},
	)
	if err != nil {
		myLogger.Error(err)
	}

	myLogger.Info("Successfully connected to Database! ALL SYSTEMS ARE GO")
	database = db
}

func GetDB() *gorm.DB {
	return database
}
