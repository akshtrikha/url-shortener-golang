package database

import (
	"context"
	"os"

	https://youtu.be/iDQAZEJK8lI?si=bcpJFjNW-ziQrLcy
)

// Ctx is the background context.
// It is used to handle the database connection operation.
var Ctx = context.Background()

// CreateClient function creates and returns the redis client.
// This redis client can be used to perform operations on the redis db.
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADD"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})

	return rdb
}
