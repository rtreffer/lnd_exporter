package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/rtreffer/lnd_exporter/pkg/exporter"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

func main() {
	listenAddress := flag.String("web.listen-address", ":9939", "Address on which to expose metrics and web interface.")
	metricsPath := flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	webConfig := flag.String("web.config.file", "", "[EXPERIMENTAL] Path to configuration file that can enable TLS or authentication.")
	collectorBasics := flag.Bool("collector.basics", true, "export lnd_info metrics")
	collectorNetwork := flag.Bool("collector.network", true, "export lnd_graph_* metrics")
	collectorChannels := flag.Bool("collector.channels", true, "export lnd_channel_* metrics")
	collectorClosedChannels := flag.Bool("collector.channels.closed", true, "export lnd_channel_* metrics for closed channels")
	collectorPendingChannels := flag.Bool("collector.channels.pending", true, "export lnd_channel_* metrics for pending channels")
	collectorPeers := flag.Bool("collector.peers", true, "export lnd_peer_* metrics")
	collectorWallet := flag.Bool("collector.wallet", true, "export lnd_wallet_* metrics")
	collectorPeerNodes := flag.Bool("collector.peer.nodes", true, "export relevant lnd_peer_node_* metrics")
	collectorFwd := flag.Bool("collector.earnings", false, "export the sum over the full forwarding history. THIS IS VERY EXPENSIVE ON STARTUP")
	flag.Parse()

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	c, err := lnc.NewClient("", "", "")
	if err != nil {
		level.Error(logger).Log("msg", "Error connecting to lnd", "err", err)
		os.Exit(1)
	}
	info, err := c.GetInfo()
	if err != nil {
		level.Error(logger).Log("msg", "Error getting node info", "err", err)
		os.Exit(1)
	}
	c.Close()

	prometheus.MustRegister(version.NewCollector("lnd_exporter"))
	exporter := &exporter.Exporter{
		NodeKey:                info.IdentityPubkey,
		CollectBasics:          *collectorBasics,
		CollectNetwork:         *collectorNetwork,
		CollectChannels:        *collectorChannels,
		CollectClosedChannels:  *collectorClosedChannels,
		CollectPendingChannels: *collectorPendingChannels,
		CollectPeers:           *collectorPeers,
		CollectPeerNodes:       *collectorPeerNodes,
		CollectWallet:          *collectorWallet,
		CollectForwardHistory:  *collectorFwd,
	}
	if err := exporter.Initialize(); err != nil {
		level.Error(logger).Log("msg", "Error initializing exporter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>LND Exporter</title></head>
             <body>
             <h1>LND Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error running HTTP server", "err", err)
		os.Exit(1)
	}
}
