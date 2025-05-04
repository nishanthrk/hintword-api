package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
)

var (
	// RedisDB is the redis connection handle
	RedisDB *redis.Client
)

func ConnectRedis() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		//Addr:     fmt.Sprintf("%s:%s", cfg.GetConfig().RedisHost, cfg.GetConfig().RedisPost), // Redis server address
		//Password: cfg.GetConfig().RedisPassword,                                              // No password set
		//DB:       utility.StringToInt(cfg.GetConfig().RedisDB),                               // Use default DB
	})

	// Ping the Redis server to check the connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println(strings.Repeat("!", 40))
		log.Printf("😏 Could Not Establish Redis Connection: %v\n", err)
		log.Println(strings.Repeat("!", 40))
		return
	}
	log.Println(strings.Repeat("-", 40))
	log.Println("😀 Connected to Redis:", pong)
	log.Println(strings.Repeat("-", 40))

	// Close the Redis client when done
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Println(strings.Repeat("!", 40))
			log.Printf("😏 Error closing Redis Connection: %v\n", err)
			log.Println(strings.Repeat("!", 40))
		}
	}()

	RedisDB = rdb
}
