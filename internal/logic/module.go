package logic

import (
	http_logic "github.com/Fiagram/gateway/internal/logic/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"logic",
	fx.Provide(
		http_logic.NewOapiLogic,
	),
)
