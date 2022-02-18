package lnc

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

func (c *Client) GetFeereport() (*lnrpc.FeeReportResponse, error) {
	client := lnrpc.NewLightningClient(grpc.ClientConnInterface((*grpc.ClientConn)(c)))
	req := &lnrpc.FeeReportRequest{}
	return client.FeeReport(context.Background(), req)
}
