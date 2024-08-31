package setup

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	kit "github.com/eden-quan/go-biz-kit"
	"github.com/eden-quan/go-biz-kit/config/def"
)

type redisImpl struct {
	db redis.UniversalClient
}

func newRedis(db redis.UniversalClient) kit.Redis {
	return &redisImpl{
		db: db,
	}
}

func (r *redisImpl) Get() redis.UniversalClient {
	return r.db
}

// NewRedis 创建 redis 客户端
func NewRedis(conf *def.Configuration, logger log.Logger) (kit.Redis, error) {
	redisConfig := &conf.Redis
	if !redisConfig.GetEnable() {
		return nil, nil
	}

	db := NewRedisClient(redisConfig, logger)

	return newRedis(db), nil
}

func NewRedisClient(config *def.Redis, logger log.Logger) redis.UniversalClient {
	lh := log.NewHelper(log.With(logger, "module", "redis"))

	opt := &redis.UniversalOptions{
		Addrs:        config.GetAddresses(),
		DB:           int(config.GetDb()),
		Username:     config.GetUsername(),
		Password:     config.GetPassword(),
		MaxRetries:   int(config.GetMaxRetries()),
		PoolSize:     int(config.GetMaxPoolSize()),
		MinIdleConns: int(config.GetMinPoolIdleSize()),
	}
	if config.GetDialTimeout() != nil {
		opt.DialTimeout = config.GetDialTimeout().AsDuration()
	}
	if config.GetReadTimeout() != nil {
		opt.ReadTimeout = config.GetReadTimeout().AsDuration()
	}
	if config.GetWriteTimeout() != nil {
		opt.WriteTimeout = config.GetWriteTimeout().AsDuration()
	}

	client := redis.NewUniversalClient(opt)
	err := client.Ping(context.Background()).Err()
	if err != nil {
		lh.Fatalw("msg", "redis ping failed", "err", err)
	}

	lh.Info("redis successfully connected and ping")
	return client
}
