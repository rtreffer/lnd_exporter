package exporter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

type ChannelsCollector struct {
	active    *prometheus.Desc
	info      func(...string) prometheus.Metric
	cap       *prometheus.Desc
	local     *prometheus.Desc
	remote    *prometheus.Desc
	unsettled *prometheus.Desc
	sent      *prometheus.Desc
	received  *prometheus.Desc
	lifetime  *prometheus.Desc
	uptime    *prometheus.Desc
	updates   *prometheus.Desc
	feeBase   *prometheus.Desc
	feeRate   *prometheus.Desc
}

const Channels = "channel"

func (c *ChannelsCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.active == nil {
		static := prometheus.Labels{"nodekey": nodekey}
		c.active = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "active"),
			"Is the channel active / live.",
			[]string{"channel_id"},
			static,
		)
		c.info = channelInfoRecorder(
			prometheus.Labels{"nodekey": nodekey, "state": "open"},
			"remotekey",
			"channel_id",
			"channel_point",
			"public",
			"opened",
			"commit_type",
			"close_address",
		)
		c.cap = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "capacity_satoshis"),
			"The total channel capacity.",
			[]string{"channel_id"},
			static,
		)
		c.local = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "local_balance_satoshis"),
			"The local channel balance.",
			[]string{"channel_id"},
			static,
		)
		c.remote = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "remote_balance_satoshis"),
			"The remote channel balance.",
			[]string{"channel_id"},
			static,
		)
		c.unsettled = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "unsettled_satoshis"),
			"The unsettled balance.",
			[]string{"channel_id", "direction"},
			static,
		)
		c.sent = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "sent_satoshis_total"),
			"The total amount sent over the channel.",
			[]string{"channel_id"},
			static,
		)
		c.received = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "received_satoshis_total"),
			"The total amount received over the channel.",
			[]string{"channel_id"},
			static,
		)
		c.lifetime = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "lifetime_seconds"),
			"The monitored lifetime of the channel.",
			[]string{"channel_id"},
			static,
		)
		c.uptime = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "uptime_seconds"),
			"The uptime of the channel.",
			[]string{"channel_id"},
			static,
		)
		c.updates = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "updates_total"),
			"The number of channel updates.",
			[]string{"channel_id"},
			static,
		)

		c.feeBase = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "base_fee_satoshis"),
			"The base fee in satoshis, with 0.001 satoshi resolution.",
			[]string{"channel_id"},
			static,
		)
		c.feeRate = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "base_rate"),
			"The fee rate as percent, ranging from [0..1].",
			[]string{"channel_id"},
			static,
		)
	}

	ch <- c.active
	ch <- channelInfo
	ch <- c.cap
	ch <- c.local
	ch <- c.remote
	ch <- c.unsettled
	ch <- c.sent
	ch <- c.received
	ch <- c.lifetime
	ch <- c.uptime
	ch <- c.updates
}

func (c *ChannelsCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric, up chan<- peerNodeFetch) {
	channels, err := client.GetChannels()
	if err != nil {
		fmt.Println("error collecting channels:", err)
		return
	}
	for _, channel := range channels.Channels {
		up <- peerNodeFetch{
			pubkey:  channel.RemotePubkey,
			channel: true,
		}
		public := "true"
		if channel.Private {
			public = "false"
		}
		initiator := "remote"
		if channel.Initiator {
			initiator = "local"
		}
		id := strconv.FormatUint(channel.ChanId, 10)
		ch <- c.info(
			channel.RemotePubkey,
			id,
			channel.ChannelPoint,
			public,
			initiator,
			strings.ToLower(channel.CommitmentType.String()),
			channel.CloseAddress)
		ch <- gauge(c.active, channel.Active, id)
		ch <- gauge(c.cap, channel.Capacity, id)
		ch <- gauge(c.local, channel.LocalBalance, id)
		ch <- gauge(c.remote, channel.RemoteBalance, id)
		ch <- gauge(c.sent, channel.TotalSatoshisSent, id)
		ch <- gauge(c.received, channel.TotalSatoshisReceived, id)
		ch <- gauge(c.lifetime, channel.Lifetime, id)
		ch <- gauge(c.uptime, channel.Uptime, id)
		ch <- gauge(c.updates, channel.NumUpdates, id)

		uin := int64(0)
		uout := int64(0)
		for _, htlc := range channel.PendingHtlcs {
			if htlc.Incoming {
				uin += htlc.Amount
			} else {
				uout += htlc.Amount
			}
		}
		ch <- gauge(c.unsettled, float64(uin), id, "incoming")
		ch <- gauge(c.unsettled, float64(uout), id, "outgoing")
	}

	fees, err := client.GetFeereport()
	if err != nil {
		fmt.Println("error collecting channel fees:", err)
		return
	}
	for _, channelFee := range fees.ChannelFees {
		id := strconv.FormatUint(channelFee.ChanId, 10)
		ch <- gauge(c.feeBase, float64(channelFee.BaseFeeMsat)*0.001, id)
		ch <- gauge(c.feeRate, channelFee.FeeRate, id)
	}
}
