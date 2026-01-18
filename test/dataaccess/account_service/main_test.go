package account_dao_test

import (
	"log"
	"os"
	"testing"

	"github.com/Fiagram/gateway/internal/configs"
	account_dao "github.com/Fiagram/gateway/internal/dataaccess/account_service"
	"go.uber.org/zap"
)

var client account_dao.Client

func TestMain(m *testing.M) {
	// Use the default config to test database connection
	cfg, err := configs.NewConfig("")
	if err != nil {
		log.Fatal("failed to init config default")
	}

	client, err = account_dao.NewClient(
		cfg.Grpc.AccountService,
		zap.NewNop(),
	)
	if err != nil {
		log.Fatal("failed to init new client for account_service")
	}

	code := m.Run()

	// Clean up
	if err := client.Close(); err != nil {
		log.Printf("failed to close account_dao client: %v", err)
	}

	os.Exit(code)
}
