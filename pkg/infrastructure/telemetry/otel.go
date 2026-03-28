// @AI_GENERATED
// Package telemetry provides lightweight observability for OpenClaw.
// It uses only the Go standard library (expvar + net/http/pprof) to avoid
// adding heavy OpenTelemetry/Prometheus SDK dependencies.
// Metrics are exposed in a Prometheus-compatible text format at /metrics.
package telemetry

import (
	"expvar"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Metrics holds all application-level counters and gauges.
type Metrics struct {
	// HTTP
	HTTPRequestsTotal   *expvar.Int
	HTTPRequestDuration *expvar.Float // cumulative seconds (for avg latency)

	// Agent
	AgentTurnDuration *expvar.Float // cumulative seconds
	TokenUsageTotal   *expvar.Int

	// Active gauges (atomic int64)
	activeSessions int64
	activeChannels int64
}

var global *Metrics

// Init initialises the global Metrics instance and registers expvar variables.
// Call once at application startup.
func Init() *Metrics {
	m := &Metrics{
		HTTPRequestsTotal:   expvar.NewInt("openclaw_http_requests_total"),
		HTTPRequestDuration: expvar.NewFloat("openclaw_http_request_duration_seconds"),
		AgentTurnDuration:   expvar.NewFloat("openclaw_agent_turn_duration_seconds"),
		TokenUsageTotal:     expvar.NewInt("openclaw_agent_token_usage_total"),
	}

	// Register gauge accessors as expvar.Func so they read the atomic values.
	expvar.Publish("openclaw_active_sessions", expvar.Func(func() any {
		return atomic.LoadInt64(&m.activeSessions)
	}))
	expvar.Publish("openclaw_active_channels", expvar.Func(func() any {
		return atomic.LoadInt64(&m.activeChannels)
	}))

	global = m
	return m
}

// Global returns the global Metrics instance (nil if Init was not called).
func Global() *Metrics { return global }

// IncActiveSessions increments the active sessions gauge.
func (m *Metrics) IncActiveSessions() { atomic.AddInt64(&m.activeSessions, 1) }

// DecActiveSessions decrements the active sessions gauge.
func (m *Metrics) DecActiveSessions() { atomic.AddInt64(&m.activeSessions, -1) }

// IncActiveChannels increments the active channels gauge.
func (m *Metrics) IncActiveChannels() { atomic.AddInt64(&m.activeChannels, 1) }

// DecActiveChannels decrements the active channels gauge.
func (m *Metrics) DecActiveChannels() { atomic.AddInt64(&m.activeChannels, -1) }

// RecordHTTPRequest records a completed HTTP request.
func (m *Metrics) RecordHTTPRequest(duration time.Duration) {
	m.HTTPRequestsTotal.Add(1)
	m.HTTPRequestDuration.Add(duration.Seconds())
}

// RecordAgentTurn records a completed agent turn.
func (m *Metrics) RecordAgentTurn(duration time.Duration, tokens int64) {
	m.AgentTurnDuration.Add(duration.Seconds())
	m.TokenUsageTotal.Add(tokens)
}

// MetricsHandler returns an http.Handler that serves Prometheus-compatible
// text format metrics derived from the expvar registry.
func MetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		expvar.Do(func(kv expvar.KeyValue) {
			// Only export openclaw_ prefixed metrics.
			if len(kv.Key) < 10 || kv.Key[:10] != "openclaw_" {
				return
			}
			fmt.Fprintf(w, "# HELP %s OpenClaw metric\n", kv.Key)
			fmt.Fprintf(w, "# TYPE %s gauge\n", kv.Key)
			fmt.Fprintf(w, "%s %s\n", kv.Key, kv.Value.String())
		})
	})
}

// @AI_GENERATED: end
