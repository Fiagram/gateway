package configs_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/Fiagram/gateway/internal/configs"
	"github.com/stretchr/testify/require"
)

var config configs.Config

func TestMain(m *testing.M) {
	cfg, err := configs.NewConfig("")
	if err != nil {
		log.Fatal("failed to init config default")
	}
	config = cfg
	os.Exit(m.Run())
}

func TestAuth(t *testing.T) {
	require.Equal(t, 24*time.Hour, config.Auth.Token.RefreshTokenTTL)
	require.Equal(t, 15*time.Minute, config.Auth.Token.AccessTokenTTL)
}
