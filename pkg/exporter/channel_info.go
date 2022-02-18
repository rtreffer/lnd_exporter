package exporter

import "github.com/prometheus/client_golang/prometheus"

var channelInfoLables = []string{
	"remotekey",
	"channel_id",
	"channel_point",
	"public",
	"opened",
	"closed",
	"close_type",
	"nodekey",
	"state",
	"commit_type",
	"close_address",
}

var channelInfo = prometheus.NewDesc(prometheus.BuildFQName(Namespace, Channels, "info"),
	"Information about the channel.",
	channelInfoLables,
	nil)

func channelInfoRecorder(static prometheus.Labels, labelName ...string) func(...string) prometheus.Metric {
	return func(labelValues ...string) prometheus.Metric {
		labelsByKey := make(map[string]string)
		for k, v := range static {
			labelsByKey[k] = v
		}
		for i, k := range labelName {
			labelsByKey[k] = labelValues[i]
		}
		values := make([]string, len(channelInfoLables))
		for i, n := range channelInfoLables {
			values[i] = labelsByKey[n]
		}
		return gauge(channelInfo, 1, values...)
	}
}
