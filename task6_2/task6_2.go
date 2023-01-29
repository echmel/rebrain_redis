package main

import (
	"os"
	"strconv"
)
import "github.com/go-redis/redis"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := redisClient.Ping().Err(); err != nil {
		os.Exit(1)
	}

	argsWithoutProg := os.Args[1:]
	listKey := argsWithoutProg[0]
	strLimit := argsWithoutProg[1]
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		panic(err)
	}

	for i := 0; i < limit; i++ {
		strCmd := redisClient.LPop(listKey)
		if err := strCmd.Err(); err != nil {
			return
		}
		bytes, err := strCmd.Bytes()
		if err != nil {
			panic(err)
		}

		if _, e := os.Stdout.Write(bytes); e != nil {
			panic(e)
		}
	}

}
