package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"os"
)

type RedisService struct {
	redis *redis.Client
}

// Redis transactions use optimistic locking.
const maxRetries = 1000

func New(redis *redis.Client) *RedisService {
	return &RedisService{
		redis: redis,
	}
}

// IncValue ...
func (r *RedisService) IncValue(ctx context.Context, key string, increment int) error {
	// Transactional function.
	txf := func(tx *redis.Tx) error {
		// Get the current value or zero.
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// Actual operation (local in optimistic lock).
		n = n + increment

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			return nil
		})
		return err
	}

	// Retry if the key has been changed.
	for i := 0; i < maxRetries; i++ {
		err := r.redis.Watch(ctx, txf, key)
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}

	return errors.New("increment reached maximum number of retries")
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		os.Exit(1)
	}

	args := os.Args[1:]
	keyName := args[0]
	incStr := args[1]

	increment, err := strconv.Atoi(incStr)
	if err != nil {
		panic(err)
	}

	service := New(redisClient)

	err = service.IncValue(ctx, keyName, increment)
	if err != nil {
		fmt.Print(err)
	}

}
