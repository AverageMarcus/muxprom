package muxprom

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MuxProm struct {
	MetricName  string
	MetricsPath string
}

func (m *MuxProm) ensureDefaults() {
	if m.MetricName == "" {
		m.MetricName = "requests"
	}
	if m.MetricsPath == "" {
		m.MetricsPath = "/metrics"
	}
}

func (m *MuxProm) RegisterPrometheus(router *mux.Router) error {
	m.ensureDefaults()

	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: m.MetricName,
		},
		[]string{"code"},
	)

	router.Handle(m.MetricsPath, promhttp.Handler()).Methods("GET")
	router.Use(m.recordMetrics(histogram))

	if router.NotFoundHandler == nil {
		router.NotFoundHandler = router.NewRoute().HandlerFunc(http.NotFound).GetHandler()
	}

	return prometheus.Register(histogram)
}

func (m *MuxProm) recordMetrics(histogram *prometheus.HistogramVec) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Don't record metrics for the metrics endpoint
			if r.URL.Path == m.MetricsPath {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			mrw := newMetricsResponseWriter(w)

			defer func() {
				duration := time.Since(start)
				histogram.WithLabelValues(strconv.Itoa(mrw.statusCode)).Observe(duration.Seconds())
			}()

			next.ServeHTTP(mrw, r)
		})
	}
}
