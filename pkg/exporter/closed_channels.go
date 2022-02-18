package exporter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

type ClosedChannelsCollector struct {
	info func(...string) prometheus.Metric
}

func (c *ClosedChannelsCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.info == nil {
		static := prometheus.Labels{"nodekey": nodekey, "state": "closed"}
		c.info = channelInfoRecorder(
			static,
			"remotekey",
			"channel_id",
			"channel_point",
			"opened",
			"closed",
			"close_type",
		)
	}

	ch <- channelInfo
}

func (c *ClosedChannelsCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric, up chan<- peerNodeFetch) {
	channels, err := client.GetClosedChannels()
	if err != nil {
		fmt.Println("error collecting closed channels:", err)
		return
	}
	for _, channel := range channels.Channels {
		up <- peerNodeFetch{
			pubkey:        channel.RemotePubkey,
			closedChannel: true,
		}
		id := strconv.FormatUint(channel.ChanId, 10)
		open := strings.ToLower(lnrpc.Initiator_name[int32(channel.OpenInitiator)])
		if strings.Index(open, "initiator_") == 0 {
			open = open[10:]
		}
		close := strings.ToLower(lnrpc.Initiator_name[int32(channel.CloseInitiator)])
		if strings.Index(close, "initiator_") == 0 {
			close = close[10:]
		}
		closeType := strings.ToLower(lnrpc.ChannelCloseSummary_ClosureType_name[int32(channel.CloseType)])
		switch channel.CloseType {
		case lnrpc.ChannelCloseSummary_ABANDONED:
			closeType = "abandoned"
		case lnrpc.ChannelCloseSummary_BREACH_CLOSE:
			closeType = "breach"
		case lnrpc.ChannelCloseSummary_COOPERATIVE_CLOSE:
			closeType = "cooperative"
		case lnrpc.ChannelCloseSummary_FUNDING_CANCELED:
			closeType = "cancelled"
		case lnrpc.ChannelCloseSummary_LOCAL_FORCE_CLOSE:
			closeType = "local_force"
		case lnrpc.ChannelCloseSummary_REMOTE_FORCE_CLOSE:
			closeType = "remote_force"
		default:
			fmt.Println("unknown channel close type", closeType)
		}
		ch <- c.info(
			channel.RemotePubkey,
			id,
			channel.ChannelPoint,
			open,
			close,
			closeType,
		)
	}
}
