package main

import (
	"bufio"
	"os"
)
import "github.com/go-redis/redis"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := redisClient.Ping().Err(); err != nil {
		os.Exit(1)
	}
	_, err := redisClient.Info("server").Result()
	if err != nil {
		os.Exit(1)
	}

	argsWithoutProg := os.Args[1:]
	listKey := argsWithoutProg[0]

	scanner := bufio.NewScanner(os.Stdin)
	pipe := redisClient.Pipeline()
	for scanner.Scan() {
		cmd := pipe.RPush(listKey, scanner.Bytes())
		if cmd.Err() != nil {
			panic(cmd.Err())
		}
		statusCmd := pipe.LTrim(listKey, -20, -1)
		if statusCmd.Err() != nil {
			panic(statusCmd.Err())
		}
		cmdArr, err := pipe.Exec()
		if err != nil {
			panic(err)
		}
		for i := range cmdArr {
			if err := cmdArr[i].Err(); err != nil {
				panic(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
