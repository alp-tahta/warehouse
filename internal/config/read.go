package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func BuiltConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(os.Getenv("DB-PORT"))
	if err != nil {
		return nil, err
	}
	return &Config{
		Port:       os.Getenv("PORT"),
		DBHost:     os.Getenv("DB-HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB-USER"),
		DBPassword: os.Getenv("DB-PASSWORD"),
		DBName:     os.Getenv("DB-NAME"),
	}, nil
}
