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

func (o *oapiLogic) RefreshToken(c *gin.Context) {
	logger := log.LoggerWithContext(c, o.logger)

	// Extract refresh token from header
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || strings.TrimSpace(refreshToken) == "" {
		errMsg := "refresh token is required"
		logger.Error(errMsg)
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: errMsg,
		})
		return
	}

	// Get account ID from refresh token cache
	accountId, err := o.refreshTokenCache.Get(c, refreshToken)
	if err != nil || accountId == 0 {
		errMsg := "invalid or expired refresh token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusUnauthorized, oapi.Unauthorized{
			Code:    "Unauthorized",
			Message: errMsg,
		})
		return
	}

	// Create a new access token
	accessToken, accessTokenExpiresAt, err := o.tokenLogic.GenerateAccessToken(c, auth_logic.TokenPayload{
		AccountId: accountId,
	})
	if err != nil {
		errMsg := "failed to generate access token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Generate new refresh token (rotation)
	newRefreshToken, newRefreshTokenExpiresAt, err := o.tokenLogic.GenerateRefreshToken(c)
	if err != nil {
		errMsg := "failed to generate refresh token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Save the new refresh token to cache
	err = o.refreshTokenCache.Set(c, newRefreshToken, accountId, o.authConfig.Token.RefreshTokenTTL)
	if err != nil {
		errMsg := "failed to save refresh token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Revoke old refresh token
	if _, err := o.refreshTokenCache.Del(c, refreshToken); err != nil {
		logger.With(zap.Error(err)).Error("failed to revoke old refresh token")
		// Continue anyway, not a critical error
	}

	// Return the refresh token to cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/v1/auth/token",
		Domain:   o.authConfig.Domain,
		Expires:  newRefreshTokenExpiresAt,
		MaxAge:   int(time.Until(newRefreshTokenExpiresAt).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Return the new access token in response
	c.JSON(http.StatusOK, oapi.RefreshResponse{
		AccessToken: oapi.AccessTokenResponse{
			Token: accessToken,
			Exp:   utils.Ptr(accessTokenExpiresAt.Unix()),
		},
	})
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
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/token",
		Domain:   o.authConfig.Domain,
		Expires:  refreshTokenExpiresAt,
		MaxAge:   int(time.Until(refreshTokenExpiresAt).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Return the access token to the response
	c.JSON(http.StatusOK, oapi.SigninResponse{
		AccessToken: oapi.AccessTokenResponse{
			Token: accessToken,
			Exp:   utils.Ptr(accessTokenExpiresAt.Unix()),
		},
	})
}

func (o *oapiLogic) SignOut(c *gin.Context) {
	logger := log.LoggerWithContext(c, o.logger)

	// Extract refresh token from header
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || strings.TrimSpace(refreshToken) == "" {
		errMsg := "refresh token is required"
		logger.Error(errMsg)
		c.JSON(http.StatusBadRequest, oapi.BadRequest{
			Code:    "BadRequest",
			Message: errMsg,
		})
		return
	}

	// Revoke the refresh token from cache
	isDone, err := o.refreshTokenCache.Del(c, refreshToken)
	if err != nil || isDone == false {
		errMsg := "failed to revoke refresh token"
		logger.With(zap.Error(err)).Error(errMsg)
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	// Return the refresh token to cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/token",
		Domain:   o.authConfig.Domain,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.Status(http.StatusNoContent)
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

	// Create a new access token
	accessToken, accessTokenExpiresAt, err := o.tokenLogic.GenerateAccessToken(c, auth_logic.TokenPayload{
		AccountId: accResp.AccountId,
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
		refreshToken, accResp.AccountId,
		o.authConfig.Token.RefreshTokenTTL)
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
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/token",
		Domain:   o.authConfig.Domain,
		Expires:  refreshTokenExpiresAt,
		MaxAge:   int(time.Until(refreshTokenExpiresAt).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Return the access token to the response
	c.JSON(http.StatusOK, oapi.SigninResponse{
		AccessToken: oapi.AccessTokenResponse{
			Token: accessToken,
			Exp:   utils.Ptr(accessTokenExpiresAt.Unix()),
		},
	})
}
