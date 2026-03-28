// @AI_GENERATED
// Package pprof provides HTTP profiling endpoints for development use.
// It exposes the standard net/http/pprof handlers on a separate port so
// they are never accidentally exposed on the main application port.
package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* handlers on http.DefaultServeMux
)

// EnableHTTP starts a pprof HTTP server on the given address in a background goroutine.
// addr should be a local address such as "localhost:6060".
// The server exposes the standard pprof endpoints:
//
//	/debug/pprof/          — index
//	/debug/pprof/profile   — CPU profile
//	/debug/pprof/heap      — heap profile
//	/debug/pprof/goroutine — goroutine profile
//	/debug/pprof/block     — block profile
//	/debug/pprof/mutex     — mutex profile
//	/debug/pprof/trace     — execution trace
func EnableHTTP(addr string) {
	go func() {
		fmt.Printf("pprof server listening on http://%s/debug/pprof/\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Printf("pprof server error: %v\n", err)
		}
	}()
}

// @AI_GENERATED: end
