package database

import (
	"context"
	"os"
	"github.com/akshtrikha/url-shortener-golang/logger"

	"github.com/go-redis/redis/v8"
)

// Ctx is the background context.
// It is used to handle the database connection operation.
var Ctx = context.Background()

// CreateClient function creates and returns the redis client.
// This redis client can be used to perform operations on the redis db.
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})

	logger.Log.Println("Created the redis client", rdb)

	return rdb
}
