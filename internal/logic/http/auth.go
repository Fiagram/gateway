package logic

import (
	"net/http"
	"strings"

	"github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"github.com/Fiagram/gateway/internal/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Refresh access token using refresh token
// (POST /auth/refresh)
func (o *oapiLogic) RefreshToken(c *gin.Context, params oapi.RefreshTokenParams) {

}

func (o *oapiLogic) SingIn(c *gin.Context) {
	// logger := utils.LoggerWithContext(c, o.logger).With(zap.Any("account", acc))

}

func (o *oapiLogic) SignOut(c *gin.Context, params oapi.SignOutParams) {
	// logger := utils.LoggerWithContext(c, o.logger).With(zap.Any("account", acc))
}

func (o *oapiLogic) SignUp(c *gin.Context) {
	logger := log.LoggerWithContext(c, o.logger)

	// Decode the incoming JSON object
	var req oapi.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.With(zap.Error(err)).Error("failed to decode incoming JSON")
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: "failed to bind JSON object",
		})
		return
	}

	// Check whether the username is taken
	username := req.Account.Username
	isTaken, err := o.usernamesTakenCache.Has(c, username)
	if isTaken {
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: "username is taken",
		})
		return
	}
	takenResp, err := o.accountGrpc.IsUsernameTaken(c, &account_service.IsUsernameTakenRequest{
		Username: username,
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to call grpc method")
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: "failed to call grpc method",
		})
		return
	}
	if takenResp.IsTaken {
		if err := o.usernamesTakenCache.Add(c, username); err != nil {
			logger.With(zap.Error(err)).
				With(zap.Any("account", req.Account)).
				Error("failed to add username to cache")
		}
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: "username is taken",
		})
		return
	}

	// Process the incoming request
	accResp, err := o.accountGrpc.CreateAccount(c, &account_service.CreateAccountRequest{
		AccountInfo: &account_service.AccountInfo{
			Username:    strings.TrimSpace(username),
			Fullname:    strings.TrimSpace(req.Account.Fullname),
			Email:       strings.TrimSpace(req.Account.Email),
			PhoneNumber: strings.TrimSpace(*req.Account.PhoneNumber.CountryCode + " " + *req.Account.PhoneNumber.Number),
			Role:        2,
		},
		Password: strings.TrimSpace(*req.Password),
	})
	if err != nil || accResp.AccountId == 0 {
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: "failed to process grpc method",
		})
		return
	}

	c.JSON(http.StatusCreated, oapi.SignupResponse{
		Username: username,
	})
}

func (o *oapiLogic) GetMe(c *gin.Context) {

}
