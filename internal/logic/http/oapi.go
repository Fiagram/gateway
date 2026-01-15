package logic

import (
	"net/http"

	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"github.com/gin-gonic/gin"
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

// Refresh access token using refresh token
// (POST /auth/refresh)
func (o oapiLogic) RefreshToken(c *gin.Context, params oapi.RefreshTokenParams) {

}

func (o oapiLogic) SingIn(c *gin.Context) {
	// logger := utils.LoggerWithContext(c, o.logger).With(zap.Any("account", acc))

}

func (o oapiLogic) SignOut(c *gin.Context, params oapi.SignOutParams) {
	// logger := utils.LoggerWithContext(c, o.logger).With(zap.Any("account", acc))
}

func (o oapiLogic) SignUp(c *gin.Context) {
	// logger := utils.LoggerWithContext(c, o.logger)
	response := oapi.SignupResponse{
		AccessToken: oapi.AccessTokenResponse{
			AccessToken: "DXMDKGLA42K21SKHV",
			ExpiresIn:   9000,
			TokenType:   oapi.Bearer,
		},
		Username: "thaivd",
	}
	c.JSON(http.StatusOK, response)

}

func (o oapiLogic) GetMe(c *gin.Context) {

}
