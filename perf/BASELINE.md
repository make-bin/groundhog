# Performance Baseline

## SLO Targets

| Endpoint | Metric | Target |
|---|---|---|
| `GET /api/v1/health` | p99 latency | < 10ms |
| `POST /api/v1/sessions` | p99 latency | < 100ms |
| `POST /api/v1/sessions/:id/messages` | p99 latency | < 500ms |
| WebSocket connect | p99 latency | < 500ms |
| WebSocket round-trip | p99 latency | < 200ms |

## Running Tests

```bash
# All tests
make perf-test

# Individual tests
make perf-test-health
make perf-test-session
make perf-test-ws

# With custom target
BASE_URL=http://staging:8080 k6 run perf/k6/health.js
```

## Prerequisites

- [k6](https://k6.io/docs/getting-started/installation/) installed
- OpenClaw server running (`make run` or `docker compose up`)
