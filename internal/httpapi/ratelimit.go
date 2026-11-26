package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
)

// RateLimit bounds how often one caller may hit the API: rate requests per
// second with a burst of burst requests. A rate of zero disables the limiter,
// which is the default because the UI fans out many queries on load and the
// CLI runs state requests in parallel.
//
// Callers are keyed by bearer token when there is one, so clients behind one
// NAT do not throttle each other, and by address otherwise. Buckets are
// dropped after ten idle minutes.
func RateLimit(rate, burst int) func(http.Handler) http.Handler {
	if rate <= 0 {
		return func(next http.Handler) http.Handler { return next }
	}
	if burst <= 0 {
		burst = rate
	}
	l := newLimiter(float64(rate), float64(burst), time.Now)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions || r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
				next.ServeHTTP(w, r)
				return
			}
			ok, retry := l.allow(key(r), l.now())
			if !ok {
				seconds := int(math.Ceil(retry.Seconds()))
				if seconds < 1 {
					seconds = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(seconds))
				jsonapi.WriteError(w, http.StatusTooManyRequests, "Too Many Requests",
					"The request rate exceeds the limit; retry after the given delay.")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func key(r *http.Request) string {
	if token := r.Header.Get("Authorization"); token != "" {
		sum := sha256.Sum256([]byte(token))
		return "token:" + hex.EncodeToString(sum[:8])
	}
	if r.RemoteAddr == "" {
		return "unknown"
	}
	return "addr:" + r.RemoteAddr
}

type bucket struct {
	tokens float64
	last   time.Time
}

type limiter struct {
	mu        sync.Mutex
	rate      float64
	burst     float64
	now       func() time.Time
	buckets   map[string]*bucket
	lastSweep time.Time
}

func newLimiter(rate, burst float64, now func() time.Time) *limiter {
	return &limiter{rate: rate, burst: burst, now: now, buckets: map[string]*bucket{}}
}

// allow takes a token from the caller's bucket. When the bucket is empty it
// reports how long the caller should wait for the next one.
func (l *limiter) allow(key string, now time.Time) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.sweep(now)
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens = math.Min(l.burst, b.tokens+elapsed*l.rate)
		b.last = now
	}
	if b.tokens >= 1 {
		b.tokens--
		return true, 0
	}
	return false, time.Duration((1 - b.tokens) / l.rate * float64(time.Second))
}

// sweep drops buckets that have not been used for ten minutes, so a scanner
// cycling tokens cannot grow the map without bound.
func (l *limiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < time.Minute {
		return
	}
	l.lastSweep = now
	for k, b := range l.buckets {
		if now.Sub(b.last) > 10*time.Minute {
			delete(l.buckets, k)
		}
	}
}
