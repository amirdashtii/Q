package repositories

import (
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	client *redis.Client
}

func createRedisConnection() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return redisClient

}


func RedisInit() *RedisDB {
	o := sync.Once{}
	var db *RedisDB
	o.Do(func() {
		db = &RedisDB{
			client: createRedisConnection(),
		}
	})
	return db
}
