package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	redisURL := os.Getenv("REDIS_URL")
	log.Printf("Connect to Redis at %s\n", redisURL)
	client = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func putToRedis(key string, value string) error {
	log.Printf("Add to Redis key=%s, value=%s", key, value)
	result, err := client.Set(key, value, 0).Result()
	ttl, _ := client.TTL(key).Result()
	log.Printf("Result %s TTL: %f\n", result, ttl.Seconds())
	return err
}

func getFromRedis(key string) (string, error) {
	log.Printf("Query from Redis key=%s\n", key)
	status := client.Get(key)
	result, err := status.Result()
	log.Printf("Get result %s value: %s\n", result, status.Val())
	return result, err
}

func initKeys(n int) error {
	// insert n keys with random values
	for i := 0; i < n; i++ {
		orgID := fmt.Sprintf("C-%06d", i)
		n := rand.Intn(10)
		effDate := time.Now().AddDate(0, 0, n-8)
		expDate := effDate.AddDate(0, 0, 30)
		value := fmt.Sprintf("%s,%s,%s", orgID, effDate.Format("2006-01-02"), expDate.Format("2006-01-02"))
		log.Printf("Put to Redis %s=%s\n", orgID, value)
		if err := putToRedis(orgID, value); err != nil {
			log.Printf("Error put key %s: %+v\n", orgID, err)
			return err
		}
	}
	return nil
}
