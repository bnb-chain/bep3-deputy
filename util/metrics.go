package util

import "github.com/prometheus/client_golang/prometheus"

var PrometheusMetrics *Metrics

const subsystem = "deputy"

type Metrics struct {
	NumSwaps           *prometheus.GaugeVec
	FetchedBlockHeight *prometheus.GaugeVec
	Balance            *prometheus.GaugeVec
}

func MustRegisterMetrics() {
	PrometheusMetrics = &Metrics{}

	// initialize metrics
	PrometheusMetrics.NumSwaps = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "num_swaps",
			Subsystem: subsystem,
			Help:      "Swap numbers",
		},
		[]string{"type", "status"},
	)

	PrometheusMetrics.FetchedBlockHeight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "fetched_block_height",
			Subsystem: subsystem,
			Help:      "Fetched block height of blockchains",
		},
		[]string{"chain"},
	)

	PrometheusMetrics.Balance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "balance",
			Subsystem: subsystem,
			Help:      "Balance of deputy account",
		},
		[]string{"chain"},
	)

	// register metrics
	prometheus.MustRegister(PrometheusMetrics.NumSwaps)
	prometheus.MustRegister(PrometheusMetrics.FetchedBlockHeight)
	prometheus.MustRegister(PrometheusMetrics.Balance)
}
