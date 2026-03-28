/**
 * k6 performance test: WebSocket connection and message round-trip
 *
 * Target: ws://localhost:8080/ws
 * SLO:    p99 round-trip < 200ms at 50 concurrent connections
 *
 * Run: k6 run perf/k6/websocket.js
 */
import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

const wsConnectLatency = new Trend('ws_connect_latency', true);
const wsRoundTripLatency = new Trend('ws_roundtrip_latency', true);
const wsErrors = new Counter('ws_errors');

export const options = {
  vus: 50,
  duration: '30s',
  thresholds: {
    ws_connect_latency: ['p(99)<500'],
    ws_roundtrip_latency: ['p(99)<200'],
    ws_errors: ['count<10'],
  },
};

const BASE_URL = __ENV.WS_URL || 'ws://localhost:8080/ws';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || '';

export default function () {
  const params = {};
  if (AUTH_TOKEN) {
    params.headers = { Authorization: `Bearer ${AUTH_TOKEN}` };
  }

  const connectStart = Date.now();

  const res = ws.connect(BASE_URL, params, function (socket) {
    wsConnectLatency.add(Date.now() - connectStart);

    socket.on('open', function () {
      // Send a ping message and measure round-trip
      const pingStart = Date.now();
      socket.send(JSON.stringify({ type: 'ping', timestamp: pingStart }));

      socket.on('message', function (data) {
        try {
          const msg = JSON.parse(data);
          if (msg.type === 'pong' || msg.type === 'ping') {
            wsRoundTripLatency.add(Date.now() - pingStart);
          }
        } catch {
          // Non-JSON message; ignore
        }
      });
    });

    socket.on('error', function (e) {
      wsErrors.add(1);
    });

    // Keep connection open briefly then close
    socket.setTimeout(function () {
      socket.close();
    }, 2000);
  });

  check(res, {
    'WebSocket connected successfully': (r) => r && r.status === 101,
  });

  sleep(0.5);
}
