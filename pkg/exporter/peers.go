package exporter

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

const Peers = "peer"

type PeersCollector struct {
	info     *prometheus.Desc
	addr     *prometheus.Desc
	bytesIn  *prometheus.Desc
	bytesOut *prometheus.Desc
	satsIn   *prometheus.Desc
	satsOut  *prometheus.Desc
	flapCnt  *prometheus.Desc
	flapTime *prometheus.Desc
}

func (c *PeersCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.info == nil {
		static := prometheus.Labels{"nodekey": nodekey}
		c.info = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "info"),
			"Information about this peer.",
			[]string{
				"remotekey",
				"open",
				"graph_sync",
			},
			static,
		)
		c.addr = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "address"),
			"Current peer address.",
			[]string{
				"remotekey",
				"address",
			},
			static,
		)
		c.bytesIn = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "received_bytes_total"),
			"Total number of received bytes.",
			[]string{"remotekey"},
			static,
		)
		c.bytesOut = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "sent_bytes_total"),
			"Total number of sent bytes.",
			[]string{"remotekey"},
			static,
		)
		c.satsIn = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "received_satoshis_total"),
			"Total number of received satoshis.",
			[]string{"remotekey"},
			static,
		)
		c.satsOut = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "sent_satoshis_total"),
			"Total number of satoshis bytes.",
			[]string{"remotekey"},
			static,
		)
		c.flapCnt = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "reconnects_total"),
			"Total number of reconnects.",
			[]string{"remotekey"},
			static,
		)
		c.flapTime = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Peers, "last_reconnect_time_seconds"),
			"Unix time of the last reconnect.",
			[]string{"remotekey"},
			static,
		)
	}

	ch <- c.info
	ch <- c.addr
	ch <- c.bytesIn
	ch <- c.bytesOut
	ch <- c.satsIn
	ch <- c.satsOut
	ch <- c.flapCnt
	ch <- c.flapTime
}

func (c *PeersCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric, up chan<- peerNodeFetch) {
	peers, err := client.GetPeers()
	if err != nil {
		fmt.Println("error collecting pending channels:", err)
		return
	}
	for _, peer := range peers.Peers {
		up <- peerNodeFetch{
			pubkey:    peer.PubKey,
			connected: true,
		}
		open := "local"
		if peer.Inbound {
			open = "remote"
		}
		graphSync := strings.ToLower(peer.SyncType.String())
		if strings.HasSuffix(graphSync, "_sync") {
			graphSync = graphSync[0 : len(graphSync)-5]
		}

		ch <- gauge(c.info, 1, peer.PubKey, open, graphSync)
		ch <- gauge(c.addr, 1, peer.PubKey, peer.Address)
		ch <- gauge(c.bytesIn, peer.BytesRecv, peer.PubKey)
		ch <- gauge(c.bytesOut, peer.BytesSent, peer.PubKey)
		ch <- gauge(c.satsIn, peer.SatRecv, peer.PubKey)
		ch <- gauge(c.satsOut, peer.SatSent, peer.PubKey)
		ch <- gauge(c.flapCnt, peer.FlapCount, peer.PubKey)
		ch <- gauge(c.flapTime, peer.LastFlapNs/1000000, peer.PubKey)
	}
}
