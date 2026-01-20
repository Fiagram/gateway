package app

import (
	"github.com/Fiagram/gateway/internal/configs"
	account_grpc "github.com/Fiagram/gateway/internal/dataaccess/account_service"
	"github.com/Fiagram/gateway/internal/dataaccess/cache"
	"github.com/Fiagram/gateway/internal/handler"
	"github.com/Fiagram/gateway/internal/log"
	"github.com/Fiagram/gateway/internal/logic"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	configs.Module,
	log.Module,

	cache.Module,
	account_grpc.Module,

	logic.Module,
	handler.Module,
)
