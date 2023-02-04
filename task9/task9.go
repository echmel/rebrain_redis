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

// Redis transactions use optimistic locking.
const maxRetries = 1000

func New(redis *redis.Client) *RedisService {
	return &RedisService{
		redis: redis,
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
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
	// 0..999999999

	/*
		Напишите программу, которая принимает единственный аргумент $name и в бесконечном цикле непрерывно генерирует
		случайные целые числа в диапазоне 0..999999999 (включительно). Из каждого такого числа $n берутся цифры среди
		содержащихся в его значении и для каждой из них считается, сколько раз она встречается в исходном числе $n,
		на основе этого формируется строка $s из десяти цифр, в которой каждая позиция соответствует своей цифре
		и представляет количество раз, которое она присутствует в числе $n. Например, для числа 982668 результатом
		будет строка 0010002021 (2 — 1 раз, 6 — 2 раза, 8 — 2 раза и 9 — 1 раз, остальных нет); для 2689 — 0010001011,
		для 4242424242 — 0050500000 и так далее. Каждое значение $n публикуется в канал с именем $name:$s.
	*/
	for {
		n := randInt(0, 999999999)
		s := makeKeySuffix(n)
		channel := name + ":" + s
		intCmd := service.redis.Publish(ctx, channel, n)
		if err := intCmd.Err(); err != nil {
			fmt.Println(err)
		}
	}

}

const sliceSize = 10

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
func makeKeySuffix(interval int) string {
	//strInt := strconv.Itoa(interval)
	tmpInterval := interval
	counter := make([]int, sliceSize)
	for {
		q, r := divmod(tmpInterval, 10)
		counter[r] = counter[r] + 1
		if q == 0 {
			break
		}
		tmpInterval = q
	}

	ret := make([]byte, sliceSize)
	for i := range counter {
		ret[i] = byte(counter[i] + 48)
	}
	return string(ret)
}
