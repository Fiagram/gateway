package account_grpc

import (
	"context"
	"fmt"

	"github.com/Fiagram/gateway/internal/configs"
	pb "github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is the interface for account service gRPC client operations
type Client interface {
	CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error)
	CheckAccountValid(ctx context.Context, req *pb.CheckAccountValidRequest) (*pb.CheckAccountValidResponse, error)
	IsUsernameTaken(ctx context.Context, in *pb.IsUsernameTakenRequest) (*pb.IsUsernameTakenResponse, error)
	GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error)
	GetAccountAll(ctx context.Context, req *pb.GetAccountAllRequest) (*pb.GetAccountAllResponse, error)
	GetAccountList(ctx context.Context, req *pb.GetAccountListRequest) (*pb.GetAccountListResponse, error)
	UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error)
	DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error)
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

func (c *client) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return c.stub.CreateAccount(ctx, req)
}

func (c *client) CheckAccountValid(ctx context.Context, req *pb.CheckAccountValidRequest) (*pb.CheckAccountValidResponse, error) {
	return c.stub.CheckAccountValid(ctx, req)
}

func (c *client) IsUsernameTaken(ctx context.Context, req *pb.IsUsernameTakenRequest) (*pb.IsUsernameTakenResponse, error) {
	return c.stub.IsUsernameTaken(ctx, req)
}

func (c *client) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	return c.stub.GetAccount(ctx, req)
}

func (c *client) GetAccountAll(ctx context.Context, req *pb.GetAccountAllRequest) (*pb.GetAccountAllResponse, error) {
	return c.stub.GetAccountAll(ctx, req)
}

func (c *client) GetAccountList(ctx context.Context, req *pb.GetAccountListRequest) (*pb.GetAccountListResponse, error) {
	return c.stub.GetAccountList(ctx, req)
}

func (c *client) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	return c.stub.UpdateAccount(ctx, req)
}

func (c *client) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	return c.stub.DeleteAccount(ctx, req)
}

func (c *client) Close() error {
	return c.conn.Close()
}
