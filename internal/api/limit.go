package api

import (
	"hash/fnv"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-freelru"
	"golang.org/x/time/rate"
)

const (
	requestLimit = 1 // per second
	burst        = 3
	capacity     = 8192
	lifeTime     = 3 * time.Minute
)

var mu sync.Mutex
var visitors = newCache()

func hashFunc(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func newCache() *freelru.SyncedLRU[string, *rate.Limiter] {
	cache, err := freelru.NewSynced[string, *rate.Limiter](capacity, hashFunc)
	if err != nil {
		panic(err)
	}
	cache.SetLifetime(lifeTime)
	return cache
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, ok := visitors.Get(ip)
	if !ok {
		limiter = rate.NewLimiter(requestLimit, burst)
		visitors.Add(ip, limiter)
	}
	return limiter
}

func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		}

		limiter := getLimiter(ip)
		if limiter.Allow() == false {

			// блокирует вызовы, пока они не прекратятся на 1 / requestLimit секунд
			// при повторении запроса во время блокировки обновляет время ожидания
			limiter.Wait(r.Context())

			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
