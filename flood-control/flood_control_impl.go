package flood_control

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// FloodControlImpl представляет реализацию интерфейса FloodControl.
type FloodControlImpl struct {
	period      time.Duration // Промежуток времени для флуд-контроля
	limit       int           // Максимальное количество запросов в промежутке времени
	redisClient *redis.Client
}

// NewFloodControl создает новый экземпляр FloodControlImpl.
func NewFloodControl(period time.Duration, limit int, redisClient *redis.Client) *FloodControlImpl {
	return &FloodControlImpl{
		period:      period,
		limit:       limit,
		redisClient: redisClient,
	}
}

// Check проверяет, пройдена ли проверка на флуд-контроль для данного пользователя.
func (fc *FloodControlImpl) Check(ctx context.Context, userID int64) (bool, error) {

	// Генерация ключа для пользователя
	key := "flood_control:" + strconv.FormatInt(userID, 10)

	currentTime := time.Now().Unix()
	// Удаляем элементы которые не входят в промежуток времени, проверяемый нами
	_, err := fc.redisClient.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(currentTime-int64(fc.period.Seconds()), 10)).Result()
	if err != nil {
		return false, err
	}

	startTime := time.Now().Add(-fc.period)

	// Получаем все временные метки запросов за последние N секунд
	requestTimes, err := fc.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(startTime.Unix(), 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return false, err
	}

	fmt.Println(requestTimes)

	// Проверяем количество запросов за последние N секунд
	if len(requestTimes) >= fc.limit {
		err = fc.redisClient.ZAdd(ctx, key, redis.Z{
			Score:  float64(currentTime),
			Member: strconv.FormatInt(currentTime, 10),
		}).Err()
		if err != nil {
			return false, err
		}

		return false, nil
	}

	// Добавляем текущий запрос в список запросов
	err = fc.redisClient.ZAdd(ctx, key, redis.Z{
		Score:  float64(currentTime),
		Member: strconv.FormatInt(currentTime, 10),
	}).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}
