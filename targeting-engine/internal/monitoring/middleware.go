package monitoring

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				labels := prometheus.Labels{
					"method": r.Method,
					"path":   r.URL.Path,
					"status": http.StatusText(ww.Status()),
				}

				metrics.HTTPRequests.With(labels).Inc()
				metrics.HTTPRequestLatency.
					With(prometheus.Labels{"method": r.Method, "path": r.URL.Path}).
					Observe(time.Since(start).Seconds())
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
