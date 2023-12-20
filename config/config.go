package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host         string
	UserName     string
	UserPassword string
	DbName       string
	Port         uint16
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	userName := os.Getenv("USER_NAME")
	userPassword := os.Getenv("USER_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	uport, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		panic(err)
	}

	return Config{
		Host:         host,
		UserName:     userName,
		UserPassword: userPassword,
		DbName:       dbName,
		Port:         uint16(uport),
	}
}
