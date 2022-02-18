package lnc

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) GetInfo() (*lnrpc.GetInfoResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.GetInfoRequest{}
	return client.GetInfo(context.Background(), req)
}

func (c *Client) GetNodeInfo(node string) (*lnrpc.NodeInfo, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.NodeInfoRequest{
		PubKey:          node,
		IncludeChannels: true,
	}
	return client.GetNodeInfo(context.Background(), req)
}
