package exporter

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

type Exporter struct {
	RpcPath string

	CollectBasics          bool
	CollectNetwork         bool
	CollectChannels        bool
	CollectClosedChannels  bool
	CollectPendingChannels bool
	CollectPeers           bool
	CollectWallet          bool
	CollectPeerNodes       bool
	CollectForwardHistory  bool

	NodeKey string
	up      *prometheus.Desc

	basics          BasicsCollector
	network         NetworkCollector
	channels        ChannelsCollector
	closedChannels  ClosedChannelsCollector
	pendingChannels PendingChannelsCollector
	peers           PeersCollector
	wallet          WalletCollector
	peerNodes       PeerNodeCollector
	fwd             FwdCollector
}

// peerNodeFetch captures why we might want to have this info
type peerNodeFetch struct {
	pubkey         string
	connected      bool
	channel        bool
	closedChannel  bool
	pendingChannel bool
}

func (p peerNodeFetch) merge(o peerNodeFetch) peerNodeFetch {
	return peerNodeFetch{
		pubkey:         p.pubkey,
		connected:      p.connected || o.connected,
		channel:        p.channel || o.channel,
		closedChannel:  p.closedChannel || o.closedChannel,
		pendingChannel: p.pendingChannel || o.pendingChannel,
	}
}

const Namespace = "lnd"

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	client, err := lnc.NewClient("", "", "")
	if err != nil {
		fmt.Println("error conencting to lnd:", err)
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return
	}
	defer client.Close()
	info, err := client.GetInfo()
	if err != nil {
		fmt.Println("error conencting to lnd:", err)
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return
	}

	// our client is working, run the different collectors in parallel, then fetch
	// extra node information if needed

	peerNodeFetchNeeded := make(chan peerNodeFetch, 5)
	peerNodeFetches := make(map[string]peerNodeFetch)

	var wg sync.WaitGroup

	if e.CollectBasics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.basics.Collect(info, ch)
		}()
	}

	if e.CollectNetwork {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.network.Collect(client, ch)
		}()
	}

	if e.CollectChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.channels.Collect(client, ch, peerNodeFetchNeeded)
		}()
	}

	if e.CollectClosedChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.closedChannels.Collect(client, ch, peerNodeFetchNeeded)
		}()
	}

	if e.CollectPendingChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.pendingChannels.Collect(client, ch, peerNodeFetchNeeded)
		}()
	}

	if e.CollectPeers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.peers.Collect(client, ch, peerNodeFetchNeeded)
		}()
	}

	if e.CollectWallet {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.wallet.Collect(client, ch)
		}()
	}

	if e.CollectForwardHistory {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.fwd.Collect(client, ch)
		}()
	}

	go func() {
		wg.Wait()
		close(peerNodeFetchNeeded)
	}()

	// check which peers we need to fetch
	for f := range peerNodeFetchNeeded {
		if e.CollectPeerNodes {
			peerNodeFetches[f.pubkey] = peerNodeFetches[f.pubkey].merge(f)
		}
	}
	if e.CollectPeerNodes {
		e.peerNodes.Collect(e.NodeKey, peerNodeFetches, client, ch)
	}
}

func (e *Exporter) Initialize() error {
	client, err := lnc.NewClient("", "", "")
	if err != nil {
		return err
	}
	if e.CollectForwardHistory {
		go e.fwd.Update(client)
	}
	return nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.up == nil {
		e.up = prometheus.NewDesc(prometheus.BuildFQName(Namespace, Basics, "up"),
			"Indicate if the lnd daemon can be reached.", nil, prometheus.Labels{"nodekey": e.NodeKey},
		)
	}

	ch <- e.up

	var wg sync.WaitGroup

	if e.CollectBasics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.basics.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectNetwork {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.network.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.channels.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectClosedChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.closedChannels.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectPendingChannels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.pendingChannels.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectPeers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.peers.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectWallet {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.wallet.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectPeerNodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.peerNodes.Describe(ch, e.NodeKey)
		}()
	}

	if e.CollectForwardHistory {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.fwd.Describe(ch, e.NodeKey)
		}()
	}

	wg.Wait()
}
