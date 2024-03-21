package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	flood_control "task/flood-control"
	"time"
)

func main() {
	// Пример использования

	viper.SetConfigFile("./config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	host := viper.GetString("redis.host")
	password := viper.GetString("redis.password")
	count := viper.GetInt("redis.countOfDataBase")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       count,
	})

	n := viper.GetInt("flood_control.time_interval")
	k := viper.GetInt("flood_control.max_requests")

	duration := time.Duration(n) * time.Second
	fc := flood_control.NewFloodControl(duration, k, redisClient)
	userID := int64(123)
	for i := 0; i < 10; i++ {
		ok, _ := fc.Check(context.Background(), userID)
		if ok {
			fmt.Printf("Request for user: %d passed\n", userID)
		} else {
			fmt.Printf("Request for user: %d rejected\n", userID)
		}
		time.Sleep(time.Second) // Имитация запросов с интервалом в 1 секунду
	}
}
