package lnc

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) GetNetworkInfo() (*lnrpc.NetworkInfo, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.NetworkInfoRequest{}
	return client.GetNetworkInfo(context.Background(), req)
}
