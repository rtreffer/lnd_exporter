package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rtreffer/lnd_exporter/pkg/lnc"
)

const Wallet = "wallet"

type WalletCollector struct {
	balance *prometheus.Desc
}

func (c *WalletCollector) Describe(ch chan<- *prometheus.Desc, nodekey string) {
	if c.balance == nil {
		static := prometheus.Labels{"nodekey": nodekey}
		c.balance = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, Wallet, "balance_satoshis"),
			"Bitcoin wallet balance.",
			[]string{
				"confirmed",
				"account",
			},
			static,
		)
	}

	ch <- c.balance
}

func (c *WalletCollector) Collect(client *lnc.Client, ch chan<- prometheus.Metric) {
	balance, err := client.GetWalletBalance()
	if err != nil {
		fmt.Println("error collecting wallet balance:", err)
		return
	}
	for account, balance := range balance.AccountBalance {
		ch <- gauge(c.balance, balance.ConfirmedBalance, "true", account)
		ch <- gauge(c.balance, balance.UnconfirmedBalance, "false", account)
	}
}
