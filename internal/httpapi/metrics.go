package httpapi

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics is this service's Prometheus registry: a request counter and latency
// histogram plus the Go runtime and process collectors.
type Metrics struct {
	registry *prometheus.Registry
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
	jobs     *prometheus.GaugeVec
	queue    *prometheus.GaugeVec
}

// NewMetrics builds an isolated registry so the exposition contains only this
// service's signals and the standard runtime ones.
func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "terrakube_http_requests_total",
		Help: "HTTP requests handled, by method, route and status.",
	}, []string{"method", "route", "status"})
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "terrakube_http_request_duration_seconds",
		Help:    "HTTP request latency in seconds, by method and route.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})
	jobs := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "terrakube_jobs",
		Help: "Live jobs by status.",
	}, []string{"status"})
	queue := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "terrakube_queue_depth",
		Help: "Messages waiting on a Redis stream.",
	}, []string{"stream"})
	registry.MustRegister(requests, duration, jobs, queue)

	return &Metrics{registry: registry, requests: requests, duration: duration, jobs: jobs, queue: queue}
}

// Middleware records one observation per request. The route label is chi's
// pattern, not the URL, so path parameters cannot inflate the series count.
func (m *Metrics) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			route := chi.RouteContext(r.Context()).RoutePattern()
			if route == "" {
				route = "unmatched"
			}
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}
			m.requests.WithLabelValues(r.Method, route, strconv.Itoa(status)).Inc()
			m.duration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
		})
	}
}

// TrackJobStatuses replaces the job gauge series, so a status that disappears
// from the database stops being reported instead of freezing at its old value.
func (m *Metrics) TrackJobStatuses(counts map[string]int64) {
	m.jobs.Reset()
	for status, count := range counts {
		m.jobs.WithLabelValues(status).Set(float64(count))
	}
}

// TrackQueueDepth sets one stream's gauge.
func (m *Metrics) TrackQueueDepth(stream string, depth int64) {
	m.queue.WithLabelValues(stream).Set(float64(depth))
}

// StatsSource is what the sampler reads. The command wires it to the database
// and the queue so this package stays independent of both.
type StatsSource interface {
	JobStatusCounts(ctx context.Context) (map[string]int64, error)
	QueueDepth(ctx context.Context) (int64, error)
}

// Sample refreshes the job and queue gauges until ctx is cancelled. A scrape
// must not run those queries itself: the job count grows with history and the
// queue call is a round trip.
func (m *Metrics) Sample(ctx context.Context, logger *slog.Logger, src StatsSource, interval time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if counts, err := src.JobStatusCounts(ctx); err != nil {
			if ctx.Err() == nil {
				logger.Warn("sample job statuses", "error", err)
			}
		} else {
			m.TrackJobStatuses(counts)
		}
		if depth, err := src.QueueDepth(ctx); err != nil {
			if ctx.Err() == nil {
				logger.Warn("sample queue depth", "error", err)
			}
		} else {
			m.TrackQueueDepth(QueueLabel, depth)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// QueueLabel is the stream label the depth gauge carries; the jobs queue is
// the only stream with a depth worth alerting on.
const QueueLabel = "jobs"

// Handler serves the exposition format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// requireMetricsToken guards the exposition. An empty token disables it.
func requireMetricsToken(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if token == "" {
			return http.NotFoundHandler()
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("Authorization")
			if len(got) > 7 && got[:7] == "Bearer " {
				got = got[7:]
			}
			if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
