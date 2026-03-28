/**
 * k6 performance test: Session creation and message sending
 *
 * Targets:
 *   POST /api/v1/sessions          — create session
 *   POST /api/v1/sessions/:id/messages — send message
 *
 * SLOs:
 *   Session creation p99 < 100ms
 *   Message send p99 < 500ms
 *
 * Run: k6 run perf/k6/session.js
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

const sessionCreateLatency = new Trend('session_create_latency', true);
const messageSendLatency = new Trend('message_send_latency', true);
const sessionErrors = new Counter('session_errors');

export const options = {
  stages: [
    { duration: '10s', target: 20 },  // ramp up
    { duration: '30s', target: 20 },  // steady state
    { duration: '10s', target: 0 },   // ramp down
  ],
  thresholds: {
    session_create_latency: ['p(99)<100'],
    message_send_latency: ['p(99)<500'],
    session_errors: ['count<5'],
    http_req_failed: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || '';

function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (AUTH_TOKEN) h['Authorization'] = `Bearer ${AUTH_TOKEN}`;
  return h;
}

export default function () {
  // 1. Create a session
  const createRes = http.post(
    `${BASE_URL}/api/v1/sessions`,
    JSON.stringify({ model: 'gpt-4o-mini', provider: 'openai' }),
    { headers: headers() }
  );

  sessionCreateLatency.add(createRes.timings.duration);

  const created = check(createRes, {
    'session created (2xx)': (r) => r.status >= 200 && r.status < 300,
  });

  if (!created) {
    sessionErrors.add(1);
    sleep(1);
    return;
  }

  let sessionID;
  try {
    const body = JSON.parse(createRes.body);
    sessionID = body.id || (body.data && body.data.id);
  } catch {
    sessionErrors.add(1);
    return;
  }

  if (!sessionID) {
    sessionErrors.add(1);
    return;
  }

  // 2. Send a message
  const msgRes = http.post(
    `${BASE_URL}/api/v1/sessions/${sessionID}/messages`,
    JSON.stringify({ content: 'Hello, world!' }),
    { headers: headers() }
  );

  messageSendLatency.add(msgRes.timings.duration);

  check(msgRes, {
    'message sent (2xx)': (r) => r.status >= 200 && r.status < 300,
  });

  sleep(0.5);
}
