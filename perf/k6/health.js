/**
 * k6 performance test: Health check endpoint
 *
 * Target: GET /api/v1/health
 * SLO:    p99 < 10ms at 100 VUs
 *
 * Run: k6 run perf/k6/health.js
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

const healthLatency = new Trend('health_latency', true);
const healthErrors = new Counter('health_errors');

export const options = {
  vus: 100,
  duration: '30s',
  thresholds: {
    // p99 must be below 10ms
    health_latency: ['p(99)<10'],
    health_errors: ['count<1'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const res = http.get(`${BASE_URL}/api/v1/health`);

  healthLatency.add(res.timings.duration);

  const ok = check(res, {
    'status is 200': (r) => r.status === 200,
    'body has status ok': (r) => {
      try {
        return JSON.parse(r.body).status === 'ok';
      } catch {
        return false;
      }
    },
  });

  if (!ok) {
    healthErrors.add(1);
  }

  sleep(0.01);
}
