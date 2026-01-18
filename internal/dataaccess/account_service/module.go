package account_dao

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"account_dao",
	fx.Provide(
		NewClient,
	),
)
