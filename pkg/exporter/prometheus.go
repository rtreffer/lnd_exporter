package exporter

import "github.com/prometheus/client_golang/prometheus"

func gauge(desc *prometheus.Desc, value interface{}, lables ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, any2float64(value), lables...)
}

func counter(desc *prometheus.Desc, value interface{}, lables ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(desc, prometheus.CounterValue, any2float64(value), lables...)
}
