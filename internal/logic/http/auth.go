package logic

import (
	"net/http"

	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"github.com/gin-gonic/gin"
)

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
		},
		Username: "thaivd",
	}
	c.JSON(http.StatusOK, response)

}

func (o oapiLogic) GetMe(c *gin.Context) {

}
