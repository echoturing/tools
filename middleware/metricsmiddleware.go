package http

import (
	"strconv"
	"sync"
	"time"

	"github.com/echoturing/log"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// const label names
	constPrometheusLabelNameService = "service"

	// variable label names
	prometheusLabelNameMethod = "method"
	prometheusLabelNamePath   = "path"
	prometheusLabelNameStatus = "status"
)

// prometheusLabelNames returns variable label names
//
// the length and order of returned slice must same as prometheusLabelValues
func prometheusLabelNames() []string {
	return []string{prometheusLabelNameMethod, prometheusLabelNamePath, prometheusLabelNameStatus}
}

// prometheusLabelValues returns variable label values
//
// the length and order of returned slice must same as prometheusLabelNames
func prometheusLabelValues(method, path string, status int) []string {
	return []string{method, path, strconv.Itoa(status)}
}

var (
	prometheusHistogramVec         *prometheus.HistogramVec
	initPrometheusHistogramVecOnce sync.Once
)

func initPrometheusHistogramVec(service string) {
	initPrometheusHistogramVecOnce.Do(func() {
		prometheusHistogramVec = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace:   "",
				Subsystem:   "",
				Name:        "http_request_duration_seconds",
				Help:        "http request duration seconds",
				ConstLabels: prometheus.Labels{constPrometheusLabelNameService: service},
				Buckets:     prometheus.DefBuckets,
			},
			prometheusLabelNames(),
		)
		prometheus.MustRegister(prometheusHistogramVec)
	})
}

// PrometheusMetrics returns a prometheus metrics middleware
func PrometheusMetrics(service string) func(echo.HandlerFunc) echo.HandlerFunc {
	initPrometheusHistogramVec(service)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			ctx := request.Context()

			method := request.Method
			path := request.URL.Path
			if method == "HEAD" {
				return next(c)
			}
			switch path {
			case "/favicon.ico":
				return next(c)
			}
			start := time.Now()
			defer func() {
				status := c.Response().Status
				// The performance of GetMetricWithLabelValues will be a little better than GetMetricWith
				observer, err := prometheusHistogramVec.GetMetricWithLabelValues(prometheusLabelValues(method, path, status)...)
				if err != nil {
					log.ErrorWithContext(ctx, "prometheus-get-metric-with-label-values-failed", "method", method, "path", path, "status", status, "error", err.Error())
					return
				}
				observer.Observe(time.Since(start).Seconds())
			}()
			return next(c)
		}
	}
}
