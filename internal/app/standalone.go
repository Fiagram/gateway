package app

import (
	"github.com/Fiagram/gateway/internal/configs"
	"github.com/Fiagram/gateway/internal/handler"
	"github.com/Fiagram/gateway/internal/log"
	"github.com/Fiagram/gateway/internal/logic"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	configs.Module,
	log.Module,
	logic.Module,
	handler.Module,
)
