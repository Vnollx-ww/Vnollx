package redis

import (
	"Vnollx/pkg/viper"
	"Vnollx/pkg/zaplog"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var client *redis.Client
var (
	config               = viper.Init("sql")
	host                 = config.Viper.GetString("redis.host")
	port                 = config.Viper.GetString("redis.port")
	password             = config.Viper.GetString("redis.password")
	db                   = config.Viper.GetInt("redis.db")
	logger   *zap.Logger = zaplog.GetLogger()
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port), // Redis 服务器地址
		Password: password,                         // 没有密码则留空
		DB:       db,                               // 使用默认数据库
	})
}
func SetKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}
func GetKeyValue(ctx context.Context, key string) (string, error) {
	value, err := client.Get(ctx, key).Result()
	return value, err
}
