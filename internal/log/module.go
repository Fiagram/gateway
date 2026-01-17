package log

import (
	"context"

	"github.com/Fiagram/gateway/internal/configs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"log",
	fx.Provide(
		NewLogger,
	),
)

func NewLogger(lc fx.Lifecycle, cfg configs.Log) (*zap.Logger, error) {
	logger, cleanup, err := InitializeLogger(cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cleanup()
			return nil
		},
	})

	return logger, nil
}
