package exporter

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

const Fwd = "fwd"

type fwdChanStats struct {
	in   uint64
	out  uint64
	fees uint64
}

type FwdCollector struct {
	updateLock sync.Mutex
	stats      map[uint64]fwdChanStats
	threshold  uint64

	fees *prometheus.Desc
	in   *prometheus.Desc
	out  *prometheus.Desc
}

func (c *FwdCollector) Update(client *lnc.Client) error {
	c.updateLock.Lock()
	defer c.updateLock.Unlock()
	if c.stats == nil {
		c.stats = make(map[uint64]fwdChanStats)
		c.threshold = uint64(time.Second)
	}
	for {
		start := time.Unix(int64(c.threshold/uint64(time.Second)), 0)
		resp, err := client.ForwardingHistory(start)
		if err != nil {
			// error :-/
			return err
		}
		if len(resp.ForwardingEvents) == 0 {
			// we are done
			return nil
		}
		newt := c.threshold
		for _, fwd := range resp.ForwardingEvents {
			if fwd.TimestampNs <= c.threshold {
				continue
			}
			if fwd.TimestampNs > newt {
				newt = fwd.TimestampNs
			}
			in := c.stats[fwd.ChanIdIn]
			in.in += fwd.AmtInMsat
			c.stats[fwd.ChanIdIn] = in
			out := c.stats[fwd.ChanIdOut]
			out.out += fwd.AmtInMsat
			out.fees += fwd.FeeMsat
			c.stats[fwd.ChanIdOut] = out
		}
		if newt == c.threshold {
			return nil
		}
		c.threshold = newt
	}
}

func (c *FwdCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.fees == nil {
		nodekey := prometheus.Labels{"nodekey": nodekey}
		c.fees = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Fwd, "fee_satoshis_total"),
			"Total number of satoshis earned with fees.",
			[]string{"channel_id"}, nodekey,
		)
		c.in = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Fwd, "received_satoshis_total"),
			"Total number of received satoshis.",
			[]string{"channel_id"}, nodekey,
		)
		c.out = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Fwd, "sent_satoshis_total"),
			"Total number of satoshis sent.",
			[]string{"channel_id"}, nodekey,
		)
	}

	ch <- c.fees
	ch <- c.in
	ch <- c.out
}

func (c *FwdCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric) {
	c.Update(client)
	c.updateLock.Lock()
	defer c.updateLock.Unlock()

	for id, stats := range c.stats {
		channelId := strconv.FormatUint(id, 10)
		ch <- gauge(c.fees, float64(stats.fees)/1000, channelId)
		ch <- gauge(c.in, float64(stats.in)/1000, channelId)
		ch <- gauge(c.out, float64(stats.out)/1000, channelId)
	}
}
