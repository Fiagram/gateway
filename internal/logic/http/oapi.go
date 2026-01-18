package logic

import (
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"go.uber.org/zap"
)

type OapiLogic interface {
	oapi.ServerInterface
}

type oapiLogic struct {
	logger *zap.Logger
}

func NewOapiLogic(
	logger *zap.Logger,
) OapiLogic {
	return &oapiLogic{
		logger: logger,
	}
}
