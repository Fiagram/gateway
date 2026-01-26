package logic

import (
	auth_logic "github.com/Fiagram/gateway/internal/logic/auth"
	http_logic "github.com/Fiagram/gateway/internal/logic/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"logic",
	fx.Provide(
		http_logic.NewOapiLogic,
		auth_logic.NewTokenLogic,
	),
)
