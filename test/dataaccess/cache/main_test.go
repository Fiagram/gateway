package cache_test

import (
	"log"
	"os"
	"testing"

	"github.com/Fiagram/gateway/internal/configs"
	"go.uber.org/zap"
)

var config configs.Cache
var logger *zap.Logger

func TestMain(m *testing.M) {
	// Use the default config to test database connection
	cfg, err := configs.NewConfig("")
	if err != nil {
		log.Fatal("failed to init config default")
	}

	config = cfg.Cache
	config.Type = configs.CacheTypeRedis
	logger = zap.NewNop()

	os.Exit(m.Run())
}
