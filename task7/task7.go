package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
)
import "github.com/go-redis/redis"

func isReadable(path string) bool {
	return unix.Access(path, unix.R_OK) == nil
}

// удаление ключей по префиксу
func deleteKeys(client *redis.Client, match string) error {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = client.Scan(cursor, match, 0).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			intCmd := client.Del(keys...)
			if intCmd.Err() != nil {
				return intCmd.Err()
			}
		}

		if cursor == 0 { // no more keys
			return nil
		}
	}
}

const basePattern = "data:%s"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := redisClient.Ping().Err(); err != nil {
		os.Exit(1)
	}

	args := os.Args[1:]
	base := args[0]

	/*
		Перед началом загрузки ключ data:${base} и все ключи, начинающиеся с data:${base}/, должны быть удалены (обратите внимание на слеш в конце).
	*/
	// удаляем сет ключей
	err := deleteKeys(redisClient, fmt.Sprintf("data:%s/*", base))
	if err != nil {
		panic(err)
	}
	// удаляем один ключ
	err = deleteKeys(redisClient, fmt.Sprintf(basePattern, base))
	if err != nil {
		panic(err)
	}
	//redisClient.Del(fmt.Sprintf("data:%s", base))

	pipe := redisClient.Pipeline()
	fileErr := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !isReadable(path) {
			return nil
		}
		mode := info.Mode()
		switch {
		case mode.IsDir():
			return nil
		case mode.IsRegular():
			fullName := path
			fileBody, err := os.ReadFile(fullName) // read file body
			if err != nil {
				return err
			}
			key := fmt.Sprintf(basePattern, fullName)
			pipe.Set(key, fileBody, 0)
			fmt.Println(path, info.Size())
		}
		return nil
	})
	if fileErr != nil {
		log.Println(fileErr)
	}
	cmd, err := pipe.Exec()
	if err != nil {
		log.Println(err)
	}
	for i := range cmd {
		if er := cmd[i].Err(); er != nil {
			log.Println(er)
		}
	}
}
