package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CreateConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	psqlinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, DBName, password)

	return gorm.Open("postgres", psqlinfo)
}
