package logic

import (
	http_logic "github.com/Fiagram/gateway/internal/logic/http"
	token_logic "github.com/Fiagram/gateway/internal/logic/token"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"logic",
	fx.Provide(
		token_logic.NewTokenLogic,

		http_logic.NewAuthLogic,
		http_logic.NewUsersLogic,
	),
)
