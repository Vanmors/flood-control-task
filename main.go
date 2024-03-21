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

	// Получаем значения из конфигурации
	host := viper.GetString("redis.host")
	password := viper.GetString("redis.password")
	count := viper.GetInt("redis.countOfDataBase")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       count,
	})

	fc := flood_control.NewFloodControl(5*time.Second, 3, redisClient) // Проверка на флуд каждые 5 секунд, максимум 3 запроса
	userID := int64(123)                                               // Идентификатор пользователя
	for i := 0; i < 10; i++ {
		ok, _ := fc.Check(context.Background(), userID)
		if ok {
			fmt.Printf("Request for user: %d passed", userID)
		} else {
			fmt.Printf("Request for user: %d rejected", userID)
		}
		time.Sleep(time.Second) // Имитация запросов с интервалом в 1 секунду
	}
}
