package lnc

import (
	"context"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) ForwardingHistory(since time.Time) (*lnrpc.ForwardingHistoryResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.ForwardingHistoryRequest{
		StartTime:    uint64(since.UTC().Unix()),
		EndTime:      uint64(time.Now().UTC().Unix()),
		NumMaxEvents: 50000,
	}
	return client.ForwardingHistory(context.Background(), req)

}
