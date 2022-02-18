package exporter

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

const PeerNode = "peer_node"

type PeerNodeCollector struct {
	info    *prometheus.Desc
	addr    *prometheus.Desc
	feeBase *prometheus.Desc
	feeRate *prometheus.Desc
}

func (c *PeerNodeCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.info == nil {
		static := prometheus.Labels{"nodekey": nodekey}
		c.info = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, PeerNode, "info"),
			"Remote peer info.",
			[]string{
				"remotekey",
				"alias",
				"color",
			},
			static,
		)
		c.addr = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, PeerNode, "address"),
			"Remote peer address.",
			[]string{
				"remotekey",
				"address",
			},
			static,
		)
		c.feeBase = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "incoming_base_fee_satoshis"),
			"Incoming base fee.",
			[]string{
				"remotekey",
				"chan_id",
			},
			static,
		)
		c.feeRate = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Channels, "incoming_fee_rate"),
			"Incoming fee rate.",
			[]string{
				"remotekey",
				"chan_id",
			},
			static,
		)
	}

	ch <- c.info
	ch <- c.addr
	ch <- c.feeBase
	ch <- c.feeRate
}

func (c *PeerNodeCollector) Collect(self string, nodes map[string]peerNodeFetch, client *lnc.Client, ch chan<- prometheus.Metric) {
	for key, _ := range nodes {
		resp, err := client.GetNodeInfo(key)
		if err != nil {
			fmt.Println("can't fetch info for", key, "-", err)
			continue
		}
		ch <- gauge(c.info, 1, resp.Node.PubKey, resp.Node.Alias, resp.Node.Color)
		for _, addr := range resp.Node.Addresses {
			ch <- gauge(c.addr, 1, resp.Node.PubKey, addr.Addr)
		}
		for _, channel := range resp.Channels {
			if channel.Node1Pub == self {
				ch <- gauge(c.feeBase,
					channel.Node2Policy.FeeBaseMsat/1000,
					resp.Node.PubKey,
					strconv.FormatUint(channel.ChannelId, 10))
				ch <- gauge(c.feeRate,
					channel.Node2Policy.FeeRateMilliMsat/1000000,
					resp.Node.PubKey,
					strconv.FormatUint(channel.ChannelId, 10))
				continue
			}
			if channel.Node2Pub == self {
				ch <- gauge(c.feeBase,
					float64(channel.Node1Policy.FeeBaseMsat)/1000,
					resp.Node.PubKey,
					strconv.FormatUint(channel.ChannelId, 10))
				ch <- gauge(c.feeRate,
					float64(channel.Node1Policy.FeeRateMilliMsat)/1000000,
					resp.Node.PubKey,
					strconv.FormatUint(channel.ChannelId, 10))
				continue
			}
		}
	}
}
