package prometheus

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_count",
		Help: "Requests count",
	}, []string{"method", "path", "status"})

	ErrorRequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "error_request_count",
		Help: "Error request count",
	}, []string{"method", "path", "status"})

	RequestCurrent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "request_current",
		Help: "Number of current requests",
	}, []string{"method", "path"})

	RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_duration",
		Help: "Requests duration in second",
	}, []string{"method", "path"})
)

func RegisterPrometheus(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestCurrent)
	prometheus.MustRegister(RequestDuration)
}
