package cache_test

import (
	"log"
	"os"
	"testing"

	"github.com/Fiagram/gateway/internal/configs"
	"github.com/Fiagram/gateway/internal/dataaccess/cache"
	"go.uber.org/zap"
)

var client cache.Client

func TestMain(m *testing.M) {
	// Use the default config to test database connection
	config, err := configs.NewConfig("")
	if err != nil {
		log.Fatal("failed to init config default")
	}

	config.Cache.Type = configs.CacheTypeRedis
	logger := zap.NewNop()

	client, err = cache.NewClient(config.Cache, logger)
	if err != nil {
		log.Fatal("failed to init redis cache")
	}

	os.Exit(m.Run())
}
