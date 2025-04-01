package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"scrooge/config"
	"scrooge/utils"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
}

func Set(key string, value interface{}) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func GetAll(keyPattern string) (map[string]string, error) {
	var cursor uint64
	var keys []string
	var err error

	for {
		var newKeys []string
		newKeys, cursor, err = rdb.Scan(ctx, cursor, keyPattern, 10).Result()
		if err != nil {
			utils.Error("Failed fetching keys from Redis: %v", err)
			return nil, err
		}
		keys = append(keys, newKeys...)

		if cursor == 0 {
			break
		}
	}

	values := make(map[string]string)
	for _, key := range keys {
		val, err := Get(key)
		if err == nil {
			values[key] = val
		}
	}

	return values, nil
}
