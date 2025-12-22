package middlewares

import (
	"net/http"
	"time"
	"youth-convention-2025/internal/logs"

	"go.uber.org/zap"
)

type Middleware struct {
	Next              http.Handler
	Secure            bool
	HTTPOnly          bool
	RequestDurMetrics bool
}

type MiddlewareOpts func(*Middleware)

func NewMiddleware(next http.Handler, opts ...MiddlewareOpts) http.Handler {
	mw := Middleware{
		Next:              next,
		Secure:            true,
		HTTPOnly:          false,
		RequestDurMetrics: false,
	}
	for _, opt := range opts {
		opt(&mw)
	}

	logs.Log().Info(
		"Middlewares",
		zap.Bool("secure", mw.Secure),
		zap.Bool("HTTP only", mw.HTTPOnly),
		zap.Bool("request duration metrics", mw.RequestDurMetrics),
	)

	return mw
}

func (mw Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	if mw.RequestDurMetrics {
		now = time.Now()
	}
	defer func() {
		if mw.RequestDurMetrics {
			logs.Log().Debug(
				"Request duration",
				zap.String("path", r.URL.Path),
				zap.Duration("dur", time.Since(now)),
			)
		}
	}()

	mw.Next.ServeHTTP(w, r)
}
