package account_grpc

import (
	"context"
	"fmt"

	"github.com/Fiagram/gateway/internal/configs"
	"github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	pb "github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is the interface for account service gRPC client operations
type Client interface {
	account_service.AccountServiceClient
	Close() error
}

type client struct {
	stub pb.AccountServiceClient
	conn *grpc.ClientConn
}

func NewClient(
	config configs.AccountService,
	logger *zap.Logger,
) (Client, error) {
	logger.With(zap.Any("account_service_config", config))

	conn, err := grpc.NewClient(
		config.Address+":"+config.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to init grpc connection")
		return nil, fmt.Errorf("failed to init grpc connection")
	}
	stub := pb.NewAccountServiceClient(conn)

	return &client{
		stub: stub,
		conn: conn,
	}, nil
}

func (c *client) CreateAccount(
	ctx context.Context,
	in *pb.CreateAccountRequest,
	opts ...grpc.CallOption,
) (*pb.CreateAccountResponse, error) {
	return c.stub.CreateAccount(ctx, in, opts...)
}

func (c *client) CheckAccountValid(
	ctx context.Context,
	in *pb.CheckAccountValidRequest,
	opts ...grpc.CallOption,
) (*pb.CheckAccountValidResponse, error) {
	return c.stub.CheckAccountValid(ctx, in, opts...)
}

func (c *client) DeleteAccount(
	ctx context.Context,
	in *pb.DeleteAccountRequest,
	opts ...grpc.CallOption,
) (*pb.DeleteAccountResponse, error) {
	return c.stub.DeleteAccount(ctx, in, opts...)
}

func (c *client) DeleteAccountByUsername(
	ctx context.Context,
	in *pb.DeleteAccountByUsernameRequest,
	opts ...grpc.CallOption,
) (*pb.DeleteAccountByUsernameResponse, error) {
	return c.stub.DeleteAccountByUsername(ctx, in, opts...)
}

func (c *client) GetAccount(
	ctx context.Context,
	in *pb.GetAccountRequest,
	opts ...grpc.CallOption,
) (*pb.GetAccountResponse, error) {
	return c.stub.GetAccount(ctx, in, opts...)
}

func (c *client) GetAccountAll(
	ctx context.Context,
	in *pb.GetAccountAllRequest,
	opts ...grpc.CallOption,
) (*pb.GetAccountAllResponse, error) {
	return c.stub.GetAccountAll(ctx, in, opts...)
}

func (c *client) GetAccountList(
	ctx context.Context,
	in *pb.GetAccountListRequest,
	opts ...grpc.CallOption,
) (*pb.GetAccountListResponse, error) {
	return c.stub.GetAccountList(ctx, in, opts...)
}

func (c *client) IsUsernameTaken(
	ctx context.Context,
	in *pb.IsUsernameTakenRequest,
	opts ...grpc.CallOption,
) (*pb.IsUsernameTakenResponse, error) {
	return c.stub.IsUsernameTaken(ctx, in, opts...)
}

func (c *client) UpdateAccountInfo(
	ctx context.Context,
	in *pb.UpdateAccountInfoRequest,
	opts ...grpc.CallOption,
) (*pb.UpdateAccountInfoResponse, error) {
	return c.stub.UpdateAccountInfo(ctx, in, opts...)
}

func (c *client) UpdateAccountPassword(
	ctx context.Context,
	in *pb.UpdateAccountPasswordRequest,
	opts ...grpc.CallOption,
) (*pb.UpdateAccountPasswordResponse, error) {
	return c.stub.UpdateAccountPassword(ctx, in, opts...)
}

func (c *client) Close() error {
	return c.conn.Close()
}
