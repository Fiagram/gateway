package account_grpc_test

import (
	"context"
	"testing"

	pb "github.com/Fiagram/gateway/internal/generated/grpc/account_service"
	"github.com/stretchr/testify/require"
)

func TestCreateAndDeleteAccount(t *testing.T) {
	ctx := context.Background()

	accInfo := RandomAccountInfo()
	accPass := RandomString(30)
	req := &pb.CreateAccountRequest{
		AccountInfo: &accInfo,
		Password:    accPass,
	}
	createAccResp, err := client.CreateAccount(ctx, req)
	require.NoError(t, err)
	require.NotEqual(t, 0, createAccResp.AccountId)

	deleteAccResp, err := client.DeleteAccount(ctx, &pb.DeleteAccountRequest{
		Username: accInfo.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, "", deleteAccResp.Username)
}

func TestCheckAccountValid(t *testing.T) {
	ctx := context.Background()

	accInfo := RandomAccountInfo()
	accPass := RandomString(30)
	req1 := &pb.CreateAccountRequest{
		AccountInfo: &accInfo,
		Password:    accPass,
	}
	createAccResp, err := client.CreateAccount(ctx, req1)
	require.NoError(t, err)
	require.NotEqual(t, 0, createAccResp.AccountId)

	req2 := &pb.CheckAccountValidRequest{
		Username: accInfo.Username,
		Password: accPass,
	}
	resp2, err := client.CheckAccountValid(ctx, req2)
	require.NoError(t, err)
	require.NotEqual(t, 0, resp2.AccountId)

	deleteAccResp, err := client.DeleteAccount(ctx, &pb.DeleteAccountRequest{
		Username: accInfo.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, "", deleteAccResp.Username)
}

// func TestGetAccount(t *testing.T) {
// 	ctx := context.Background()

// 	req := &pb.GetAccountRequest{
// 		AccountId: 1,
// 	}

// 	resp, err := client.GetAccount(ctx, req)
// 	if err != nil {
// 		t.Fatalf("GetAccount failed: %v", err)
// 	}

// 	if resp.AccountId == 0 {
// 		t.Error("expected account_id to be non-zero")
// 	}

// 	if resp.Account == nil {
// 		t.Error("expected account info to be present")
// 	}

// 	t.Logf("Successfully retrieved account: %+v", resp.Account)
// }

// func TestGetAccountAll(t *testing.T) {
// 	ctx := context.Background()

// 	req := &pb.GetAccountAllRequest{}

// 	resp, err := client.GetAccountAll(ctx, req)
// 	if err != nil {
// 		t.Fatalf("GetAccountAll failed: %v", err)
// 	}

// 	if len(resp.AccountIdList) == 0 {
// 		t.Log("No accounts found in the system")
// 	}

// 	if len(resp.AccountIdList) != len(resp.AccountInfoList) {
// 		t.Error("account_id_list and account_info_list length mismatch")
// 	}

// 	t.Logf("Successfully retrieved %d accounts", len(resp.AccountIdList))
// }

// func TestGetAccountList(t *testing.T) {
// 	ctx := context.Background()

// 	req := &pb.GetAccountListRequest{
// 		AccountIdList: []uint64{1, 2},
// 	}

// 	resp, err := client.GetAccountList(ctx, req)
// 	if err != nil {
// 		t.Fatalf("GetAccountList failed: %v", err)
// 	}

// 	if len(resp.AccountIdList) != len(resp.AccountInfoList) {
// 		t.Error("account_id_list and account_info_list length mismatch")
// 	}

// 	t.Logf("Successfully retrieved %d accounts from list", len(resp.AccountIdList))
// }

// func TestUpdateAccount(t *testing.T) {
// 	ctx := context.Background()

// 	req := &pb.UpdateAccountRequest{
// 		AccountId: 1,
// 		UpdatedAccountInfo: &pb.AccountInfo{
// 			Username:    "testuser_updated",
// 			Fullname:    "Updated Test User",
// 			Email:       "updated@example.com",
// 			PhoneNumber: "9876543210",
// 			Role:        pb.AccountInfo_MEMBER,
// 		},
// 	}

// 	resp, err := client.UpdateAccount(ctx, req)
// 	if err != nil {
// 		t.Fatalf("UpdateAccount failed: %v", err)
// 	}

// 	if resp.AccountId != 1 {
// 		t.Errorf("expected account_id to be 1, got %d", resp.AccountId)
// 	}

// 	t.Logf("Successfully updated account with ID: %d", resp.AccountId)
// }
