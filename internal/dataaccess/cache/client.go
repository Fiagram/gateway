package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Fiagram/gateway/internal/configs"
	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Del(ctx context.Context, key ...string) error
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataInSet(ctx context.Context, key string, data any) (bool, error)
}

func NewClient(
	config configs.Cache,
	logger *zap.Logger,
) (Client, error) {
	switch config.Type {
	case configs.CacheTypeRam:
		return NewRamClient(logger), nil
	case configs.CacheTypeRedis:
		return NewRedisClient(config, logger), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", config.Type)
	}

}
