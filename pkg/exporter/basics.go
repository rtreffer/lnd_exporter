package exporter

import (
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
)

const Basics = ""

type BasicsCollector struct {
	info        *prometheus.Desc
	versionInfo *prometheus.Desc
	chainSynced *prometheus.Desc
	graphSynced *prometheus.Desc
	blockHeight *prometheus.Desc
	blockTime   *prometheus.Desc
	peers       *prometheus.Desc
	channels    *prometheus.Desc
}

func (b *BasicsCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if b.info == nil {
		static := prometheus.Labels{"nodekey": nodekey}
		b.info = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "node_info"),
			"Information about the node.",
			[]string{"alias", "color", "chain", "network"},
			static,
		)
		b.versionInfo = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "info"),
			"Information about the lnd software version.",
			[]string{"version", "commit"},
			static,
		)
		b.chainSynced = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "chain_synced"),
			"1 if this node is synced up with the bitcoin chain.",
			nil,
			static,
		)
		b.graphSynced = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "graph_synced"),
			"1 if this node is synced up with the lightning graph.",
			nil,
			static,
		)
		b.blockHeight = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "block_height_total"),
			"Current blockchain height as known to lnd.",
			nil,
			static,
		)
		b.blockTime = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "latest_block_time_seconds"),
			"Unix time of the latest blockchain block.",
			nil,
			static,
		)
		b.peers = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "peers_total"),
			"Number of peers.",
			nil,
			static,
		)
		b.channels = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Basics, "channels_total"),
			"Number of channels in a given state.",
			[]string{"state"},
			static,
		)
	}

	ch <- b.info
	ch <- b.versionInfo
	ch <- b.chainSynced
	ch <- b.graphSynced
	ch <- b.blockHeight
	ch <- b.blockTime
	ch <- b.peers
	ch <- b.channels
}

func (b *BasicsCollector) Collect(info *lnrpc.GetInfoResponse, ch chan<- prometheus.Metric) {
	// this will usually be one row, bitcoin/mainnet
	// split it if multiple networks become common
	for _, chain := range info.Chains {
		ch <- gauge(b.info, 1, info.Alias, info.Color, chain.Chain, chain.Network)
	}
	ch <- gauge(b.versionInfo, 1, info.Version, info.CommitHash)
	ch <- gauge(b.chainSynced, info.SyncedToChain)
	ch <- gauge(b.graphSynced, info.SyncedToGraph)
	ch <- gauge(b.blockHeight, info.BlockHeight)
	ch <- gauge(b.blockTime, info.BestHeaderTimestamp)
	ch <- gauge(b.peers, info.NumPeers)
	ch <- gauge(b.channels, info.NumActiveChannels, "active")
	ch <- gauge(b.channels, info.NumInactiveChannels, "inactive")
	ch <- gauge(b.channels, info.NumPendingChannels, "pending")
}
