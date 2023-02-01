package main

import (
	"fmt"
	"os"
	"strings"
)
import "github.com/go-redis/redis"

const prefix = "data:"

func scanKeys(client *redis.Client, match string) error {
	strSliceCmd := client.Keys(match)
	err := strSliceCmd.Err()
	if err != nil {
		return err
	}

	keys, err := strSliceCmd.Result()
	if err != nil {
		return err
	}

	for i := range keys {
		in := client.StrLen(keys[i]) //  $size $file.
		size := in.Val()
		fileName := strings.TrimPrefix(keys[i], prefix)
		fmt.Println(fmt.Sprintf("%d %s", size, fileName))
	}
	return nil
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := redisClient.Ping().Err(); err != nil {
		os.Exit(1)
	}
	// получает и печатате все ключи
	err := scanKeys(redisClient, prefix+"*")
	if err != nil {
		fmt.Print(err)
	}

}
