//go:build !solution

package httpgauge

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Gauge struct {
	mu      sync.Mutex
	metrics map[string]int
}

func New() *Gauge {
	return &Gauge{
		metrics: make(map[string]int),
	}
}

func (g *Gauge) Snapshot() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()
	copyMap := make(map[string]int, len(g.metrics))
	for k, v := range g.metrics {
		copyMap[k] = v
	}
	return copyMap
}

// ServeHTTP returns accumulated statistics in text format ordered by pattern.
//
// For example:
//
//	/a 10
//	/b 5
//	/c/{id} 7
func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	type entry struct {
		pattern string
		count   int
	}
	var entries []entry
	for pattern, count := range g.metrics {
		entries = append(entries, entry{pattern: pattern, count: count})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].pattern < entries[j].pattern
	})
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "%s %d\n", e.pattern, e.count)
		if err != nil {
			return
		}
	}
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				routeCtx := chi.RouteContext(r.Context())
				if routeCtx != nil {
					pattern := routeCtx.RoutePattern()
					if pattern == "" && len(routeCtx.RoutePatterns) > 0 {
						pattern = routeCtx.RoutePatterns[len(routeCtx.RoutePatterns)-1]
					}

					if pattern != "" {
						g.mu.Lock()
						g.metrics[pattern]++
						g.mu.Unlock()
					}
				}
				panic(rec)
			}
		}()
		next.ServeHTTP(w, r)
		routeCtx := chi.RouteContext(r.Context())
		if routeCtx != nil {
			pattern := routeCtx.RoutePattern()
			if pattern == "" && len(routeCtx.RoutePatterns) > 0 {
				pattern = routeCtx.RoutePatterns[len(routeCtx.RoutePatterns)-1]
			}

			if pattern != "" {
				g.mu.Lock()
				g.metrics[pattern]++
				g.mu.Unlock()
			}
		}
	})
}
