package logic

import (
	"net/http"
	"strings"
	"time"

	"github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"github.com/Fiagram/gateway/internal/log"
	auth_logic "github.com/Fiagram/gateway/internal/logic/auth"
	"github.com/Fiagram/gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Refresh access token using refresh token
// (POST /auth/refresh)
func (o *oapiLogic) RefreshToken(c *gin.Context, params oapi.RefreshTokenParams) {
	// logger := log.LoggerWithContext(c, o.logger)

	// token, err := c.Cookie("refresh_token")
	// if err != nil {
	// 	errMsg := "failed to get the refresh token"
	// 	logger.With(zap.Error(err)).Error(errMsg)
	// 	c.JSON(http.StatusUnauthorized, oapi.Unauthorized{
	// 		Code:    "Unauthorized",
	// 		Message: errMsg,
	// 	})
	// 	return
	// }

	// // check token

	// c.JSON(http.StatusOK, oapi.RefreshResponse{
	// 	AccessToken: oapi.AccessTokenResponse{
	// 		Token:     "...",
	// 		ExpiresAt: 9999,
	// 	},
	// 	Username: "xxx",
	// })
}

func (o *oapiLogic) SignIn(c *gin.Context) {
	logger := log.LoggerWithContext(c, o.logger)

	// Decode the incoming JSON object
	var req oapi.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := "failed to bind JSON object"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: errMsg,
		})
		return
	}

	// Verify data input is not empty
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(*req.Password)
	isRememberMe := *req.IsRememberMe
	logger.With(zap.String("username", username))
	if username == "" || password == "" {
		errMsg := "invalid username or password"
		logger.Error(errMsg)
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: errMsg,
		})
		return
	}

	// Checking account valid
	validResp, err := o.accountGrpc.CheckAccountValid(c,
		&account_service.CheckAccountValidRequest{
			Username: username,
			Password: password,
		})
	if err != nil {
		errMsg := "failed to check account valid"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	} else if validResp.AccountId == 0 {
		c.JSON(http.StatusUnauthorized, oapi.Unauthorized{
			Code:    "Unauthorized",
			Message: "invalid username or password",
		})
		return
	}

	// Create a new access token
	accessToken, accessTokenExpiresAt, err := o.tokenLogic.GenerateAccessToken(c, auth_logic.TokenPayload{
		AccountId: validResp.AccountId,
	})
	if err != nil {
		errMsg := "failed to gen access token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Create refresh token
	refreshToken, refreshTokenExpiresAt, err := o.tokenLogic.GenerateRefreshToken(c)
	if err != nil {
		errMsg := "failed to gen refresh token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Save the refresh token to the redis
	err = o.refreshTokenCache.Set(c,
		refreshToken, validResp.AccountId,
		utils.If(isRememberMe,
			o.authConfig.Token.RefreshTokenLongTTL,
			o.authConfig.Token.RefreshTokenTTL),
	)
	if err != nil {
		errMsg := "failed to save refresh token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Return the refresh token to cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(time.Until(refreshTokenExpiresAt).Seconds()),
		"/auth/refresh",
		"localhost:8080",
		true,
		true,
	)

	// Return the access token to the response
	c.JSON(http.StatusOK, oapi.SigninResponse{
		AccessToken: oapi.AccessTokenResponse{
			Token:         accessToken,
			ExpiresInSecs: int(time.Until(accessTokenExpiresAt).Seconds()),
		},
		Username: username,
	})
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
