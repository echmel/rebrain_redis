package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"os"
	"strconv"
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

/*
Напишите программу, которая принимает в качестве первого аргумента значение $count,
а второго и последующих — имена ключей. Программа должна дождаться ровно $count значений элементов из списков,
находящихся в указанных ключах и вывести каждое из них на своей строке в формате $key $item.
Общее время ожидания не должно превышать 5 секунд, иначе программа должна завершиться с кодом, отличным от 0.

Ключ cmd-wait должен содержать команду для запуска этой программы.
*/
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
	/*
		count, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}*/
	count := args[0]
	intCount, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lists := args[1:]

	ch := make(chan struct{})
	service := New(redisClient)
	timeCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	go func() {
		funcName(timeCtx, ch, intCount, lists, service)
	}()

	select {
	case <-timeCtx.Done():
		_, ok := timeCtx.Deadline()
		if !ok {
			os.Exit(1)
		}
		if err := timeCtx.Err(); err != nil {
			os.Exit(1)
		}
	case <-ch:
	}
}

func funcName(timeCtx context.Context, ch chan struct{}, intCount int, lists []string, service *RedisService) {
	for intCount > 0 {
		// BLMPOP timeout numkeys key [key ...] LEFT|RIGHT [COUNT count]
		keyNumbers := strconv.Itoa(len(lists))
		cmd := []string{"BLMPOP", "5", keyNumbers}
		cmd = append(cmd, lists...)
		cmd = append(cmd, "RIGHT", "COUNT")
		lCount := strconv.Itoa(intCount)
		cmd = append(cmd, lCount)

		b := make([]interface{}, len(cmd), len(cmd))
		for i := range cmd {
			b[i] = cmd[i]
		}
		res, err := service.redis.Do(timeCtx, b...).Result()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sliceArray, ok := res.([]interface{})
		if !ok {
			fmt.Println("failed to cast interface response")
			os.Exit(1)
		}
		key := sliceArray[0].(string)
		strSlice := interfaceSLiceToStrSlice(sliceArray[1].([]interface{}))
		intCount = intCount - len(strSlice)
		for _, v := range strSlice {
			fmt.Println(fmt.Sprintf("%s %s", key, v))
		}
	}
	ch <- struct{}{}
}

func interfaceSLiceToStrSlice(vals []interface{}) []string {
	b := make([]string, len(vals), len(vals))
	for i := range vals {
		b[i] = vals[i].(string)
	}
	return b
}
