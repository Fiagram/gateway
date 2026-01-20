package logic

import (
	account_grpc "github.com/Fiagram/gateway/internal/dataaccess/account_service"
	"github.com/Fiagram/gateway/internal/dataaccess/cache"
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"go.uber.org/zap"
)

type OapiLogic interface {
	oapi.ServerInterface
}

type oapiLogic struct {
	logger              *zap.Logger
	usernamesTakenCache cache.UsernamesTaken
	accountGrpc         account_grpc.Client
}

func NewOapiLogic(
	logger *zap.Logger,
	usernamesTakenCache cache.UsernamesTaken,
	accountGrpc account_grpc.Client,
) OapiLogic {
	return &oapiLogic{
		logger:              logger,
		usernamesTakenCache: usernamesTakenCache,
		accountGrpc:         accountGrpc,
	}
}
