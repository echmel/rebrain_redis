package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"os"
	"time"
)

type RedisService struct {
	redis *redis.Client
}

func New(redis *redis.Client) *RedisService {
	return &RedisService{
		redis: redis,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		os.Exit(1)
	}

	args := os.Args[1:]
	name := args[0]

	service := New(redisClient)

	pattern := name + ":" + "?0?0?0?[3456789]?0"
	pubSub := service.redis.PSubscribe(ctx, pattern)
	msg, err := pubSub.ReceiveMessage(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg.Payload)
}
