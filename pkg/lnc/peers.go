package lnc

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) GetPeers() (*lnrpc.ListPeersResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.ListPeersRequest{}
	return client.ListPeers(context.Background(), req)
}
