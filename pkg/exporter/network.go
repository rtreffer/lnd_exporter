package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

const Network = "graph"

type NetworkCollector struct {
	nodes             *prometheus.Desc
	channels          *prometheus.Desc
	degree            *prometheus.Desc
	capacity          *prometheus.Desc
	medianChannelSize *prometheus.Desc
	minChannelSize    *prometheus.Desc
	maxChannelSize    *prometheus.Desc
}

func (n *NetworkCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if n.nodes == nil {
		nodekey := prometheus.Labels{"nodekey": nodekey}
		n.nodes = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "nodes"),
			"Total number of nodes in the graph.",
			nil, nodekey,
		)
		n.channels = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "channels"),
			"Total number of channels in the graph.",
			nil, nodekey,
		)
		n.degree = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "max_degree"),
			"Maximum number of outgoing channels.",
			nil, nodekey,
		)
		n.capacity = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "capacity_satoshis"),
			"Total graph capacity.",
			nil, nodekey,
		)
		n.medianChannelSize = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "median_channel_size_satoshis"),
			"Median channel size in the lightning graph.",
			nil, nodekey,
		)
		n.maxChannelSize = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "max_channel_size_satoshis"),
			"Maximum channel size in the lightning graph.",
			nil, nodekey,
		)
		n.minChannelSize = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Network, "min_channel_size_satoshis"),
			"Minimum channel size in the lightning graph.",
			nil, nodekey,
		)
	}
	ch <- n.nodes
	ch <- n.channels
	ch <- n.degree
	ch <- n.capacity
	ch <- n.medianChannelSize
	ch <- n.minChannelSize
	ch <- n.maxChannelSize
}

func (n *NetworkCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric) {
	info, err := client.GetNetworkInfo()
	if err != nil {
		return
	}

	ch <- gauge(n.nodes, info.NumNodes)
	ch <- gauge(n.channels, info.NumChannels)
	ch <- gauge(n.degree, info.MaxOutDegree)
	ch <- gauge(n.capacity, info.TotalNetworkCapacity)
	ch <- gauge(n.minChannelSize, info.MinChannelSize)
	ch <- gauge(n.maxChannelSize, info.MaxChannelSize)
	ch <- gauge(n.medianChannelSize, info.MedianChannelSizeSat)
}
