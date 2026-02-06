package logic

import (
	"net/http"

	account_grpc "github.com/Fiagram/gateway/internal/dataaccess/account_service"
	"github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	oapi "github.com/Fiagram/gateway/internal/generated/openapi"
	"github.com/Fiagram/gateway/internal/log"
	"github.com/Fiagram/gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UsersLogic interface {
	GetMe(c *gin.Context)
}

var _ UsersLogic = (oapi.ServerInterface)(nil)

type usersLogic struct {
	accountGrpc account_grpc.Client
	logger      *zap.Logger
}

func NewUsersLogic(
	accountGrpc account_grpc.Client,
	logger *zap.Logger,
) UsersLogic {
	return &usersLogic{
		accountGrpc: accountGrpc,
		logger:      logger,
	}
}

func (u *usersLogic) GetMe(c *gin.Context) {
	logger := log.LoggerWithContext(c, u.logger)

	accountId, exists := c.Get("accountId")
	if !exists || accountId.(uint64) == 0 {
		errMsg := "accountId not existed in context"
		logger.Error(errMsg)
		c.JSON(http.StatusUnauthorized, oapi.Unauthorized{
			Code:    "Unauthorized",
			Message: errMsg,
		})
		return
	}

	account, err := u.accountGrpc.GetAccount(c, &account_service.GetAccountRequest{
		AccountId: accountId.(uint64),
	})
	if err != nil {
		errMsg := "failed to get account from account service"
		logger.Error(errMsg, zap.Error(err))
		c.JSON(http.StatusInternalServerError, oapi.InternalServerError{
			Code:    "InternalServerError",
			Message: errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, oapi.UsersMeResponse{
		Account: oapi.Account{
			Username: account.Account.Username,
			Fullname: account.Account.Fullname,
			Email:    account.Account.Email,
			PhoneNumber: &oapi.PhoneNumber{
				CountryCode: utils.Ptr("none"), // TODO: fill proper country code
				Number:      utils.Ptr(account.Account.PhoneNumber),
			},
			Role: "member", // TODO: fill proper role with converters int <-> string
		},
	})

}
