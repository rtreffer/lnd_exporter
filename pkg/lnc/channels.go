package lnc

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) GetChannels() (*lnrpc.ListChannelsResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.ListChannelsRequest{}
	return client.ListChannels(context.Background(), req)
}

func (c *Client) GetClosedChannels() (*lnrpc.ClosedChannelsResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.ClosedChannelsRequest{}
	return client.ClosedChannels(context.Background(), req)
}

func (c *Client) GetPendingChannels() (*lnrpc.PendingChannelsResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.PendingChannelsRequest{}
	return client.PendingChannels(context.Background(), req)
}
