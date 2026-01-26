package logic

import (
	"github.com/Fiagram/gateway/internal/configs"
	account_grpc "github.com/Fiagram/gateway/internal/dataaccess/account_service"
	"github.com/Fiagram/gateway/internal/dataaccess/cache"
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	auth_logic "github.com/Fiagram/gateway/internal/logic/auth"
	"go.uber.org/zap"
)

type OapiLogic interface {
	oapi.ServerInterface
}

type oapiLogic struct {
	authConfig          configs.Auth
	usernamesTakenCache cache.UsernamesTaken
	refreshTokenCache   cache.RefreshToken
	accountGrpc         account_grpc.Client
	tokenLogic          auth_logic.Token
	logger              *zap.Logger
}

func NewOapiLogic(
	authConfig configs.Auth,
	usernamesTakenCache cache.UsernamesTaken,
	refreshTokenCache cache.RefreshToken,
	accountGrpc account_grpc.Client,
	tokenLogic auth_logic.Token,
	logger *zap.Logger,
) OapiLogic {
	return &oapiLogic{
		authConfig:          authConfig,
		usernamesTakenCache: usernamesTakenCache,
		refreshTokenCache:   refreshTokenCache,
		accountGrpc:         accountGrpc,
		tokenLogic:          tokenLogic,
		logger:              logger,
	}
}
