package utils

import (
	"context"
	"fmt"
	"time"

	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitCache() {
	host := Config.RedisHost
	port := Config.RedisPort
	pswd := Config.RedisPassword

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		DB:           0,
		Password:     pswd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  2 * time.Second,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		logger.Errorf(err.Error())
		panic(err)
	}
	logger.Infof("Connected to Redis")
}
