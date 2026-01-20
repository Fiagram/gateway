package account_grpc

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"account_grpc",
	fx.Provide(
		NewClient,
	),
)
