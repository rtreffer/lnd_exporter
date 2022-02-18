package lnc

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) GetWalletBalance() (*lnrpc.WalletBalanceResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.WalletBalanceRequest{}
	return client.WalletBalance(context.Background(), req)
}
