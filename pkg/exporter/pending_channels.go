package exporter

import (
	"fmt"
	"strings"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

const PendingChannels = "pending_channel"

type PendingChannelsCollector struct {
	info   func(...string) prometheus.Metric
	cap    *prometheus.Desc
	local  *prometheus.Desc
	remote *prometheus.Desc
}

func (c *PendingChannelsCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.info == nil {
		static := prometheus.Labels{"nodekey": nodekey, "state": "pending"}
		c.info = channelInfoRecorder(
			static,
			"remotekey",
			"channel_point",
			"opened",
			"commit_type",
		)
		c.cap = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, PendingChannels, "capacity_satoshis"),
			"The total channel capacity.",
			[]string{"remotekey", "channel_point"},
			static,
		)
		c.local = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, PendingChannels, "local_balance_satoshis"),
			"The local channel balance.",
			[]string{"remotekey", "channel_point"},
			static,
		)
		c.remote = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, PendingChannels, "remote_balance_satoshis"),
			"The remote channel balance.",
			[]string{"remotekey", "channel_point"},
			static,
		)

	}

	ch <- channelInfo
	ch <- c.cap
	ch <- c.local
	ch <- c.remote
}

func (c *PendingChannelsCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric, up chan<- peerNodeFetch) {
	channels, err := client.GetPendingChannels()
	if err != nil {
		fmt.Println("error collecting pending channels:", err)
		return
	}
	for _, channel := range channels.PendingOpenChannels {
		up <- peerNodeFetch{
			pubkey:         channel.Channel.RemoteNodePub,
			pendingChannel: true,
		}
		open := strings.ToLower(lnrpc.Initiator_name[int32(channel.Channel.Initiator)])
		if strings.Index(open, "initiator_") == 0 {
			open = open[10:]
		}
		ch <- c.info(
			channel.Channel.RemoteNodePub,
			channel.Channel.ChannelPoint,
			open,
			strings.ToLower(channel.Channel.CommitmentType.String()),
		)
		ch <- gauge(c.cap, channel.Channel.Capacity, channel.Channel.RemoteNodePub, channel.Channel.ChannelPoint)
		ch <- gauge(c.local, channel.Channel.LocalBalance, channel.Channel.RemoteNodePub, channel.Channel.ChannelPoint)
		ch <- gauge(c.remote, channel.Channel.RemoteBalance, channel.Channel.RemoteNodePub, channel.Channel.ChannelPoint)
	}
}
