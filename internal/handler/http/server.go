package handler

import (
	"context"

	"github.com/Fiagram/gateway/internal/configs"
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"github.com/Fiagram/gateway/internal/log"
	logic "github.com/Fiagram/gateway/internal/logic/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpServer interface {
	Start(ctx context.Context) error
}

type httpServer struct {
	httpConfig configs.Http
	oapiLogic  logic.OapiLogic
	logger     *zap.Logger
}

func NewHttpServer(
	httpConfig configs.Http,
	oapiLogic logic.OapiLogic,
	logger *zap.Logger,
) HttpServer {
	return &httpServer{
		httpConfig: httpConfig,
		oapiLogic:  oapiLogic,
		logger:     logger,
	}
}

func (s httpServer) Start(ctx context.Context) error {
	logger := log.LoggerWithContext(ctx, s.logger)

	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	oapi.RegisterHandlers(apiV1, s.oapiLogic)

	address := s.httpConfig.Address
	port := s.httpConfig.Port
	logger.With(zap.String("address", address)).
		With(zap.String("port", port)).
		Info("starting http server")

	return r.Run(address + ":" + port)
}
