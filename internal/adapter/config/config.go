package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		DB    *DB
		Cache *Cache
	}
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
	Cache struct {
		Host     string
		Port     string
		User     string
		Password string
	}
)

func New() (*Container, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	cache := &Cache{
		Host:     os.Getenv("CACHE_HOST"),
		Port:     os.Getenv("CACHE_PORT"),
		User:     os.Getenv("CACHE_USER"),
		Password: os.Getenv("CACHE_PASSWORD"),
	}

	return &Container{
		DB:    db,
		Cache: cache,
	}, nil
}
