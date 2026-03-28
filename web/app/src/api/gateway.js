const BASE = '/api/v1'

function getUserID() {
  return localStorage.getItem('user_id') || 'user-001'
}

async function request(method, path, body, extraHeaders = {}) {
  const headers = {
    'Content-Type': 'application/json',
    'X-User-ID': getUserID(),
    ...extraHeaders,
  }
  const token = localStorage.getItem('token')
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  if (res.status === 204) return null

  const text = await res.text()
  if (!res.ok) throw new Error(`${method} ${path} → ${res.status}: ${text}`)

  try {
    const json = JSON.parse(text)
    // Unwrap standard envelope { code, data }
    return json.data !== undefined ? json.data : json
  } catch {
    return text
  }
}

const get  = (path, h) => request('GET',    path, undefined, h)
const post = (path, b, h) => request('POST', path, b, h)
const put  = (path, b) => request('PUT',    path, b)
const del  = (path)    => request('DELETE', path)

const api = {
  health: () => get('/health'),

  sessions: {
    list:        (params) => get('/sessions' + (params ? '?' + new URLSearchParams(params) : '')),
    get:         (id)     => get(`/sessions/${id}`),
    create:      (data)   => post('/sessions', data),
    delete:      (id)     => del(`/sessions/${id}`),
    sendMessage: (id, data) => post(`/sessions/${id}/messages`, data),
    resolveApproval: (id, approvalId, decision) =>
      post(`/sessions/${id}/approvals/${approvalId}`, { decision }),
    listApprovals: (id) => get(`/sessions/${id}/approvals`),
  },

  channels: {
    list:   ()     => get('/channels'),
    create: (data) => post('/channels', data),
    delete: (id)   => del(`/channels/${id}`),
    status: (id)   => get(`/channels/${id}/status`),
  },

  memories: {
    list:   ()           => get('/memories'),
    create: (content)    => post('/memories', { content }),
    get:    (id)         => get(`/memories/${id}`),
    update: (id, content) => put(`/memories/${id}`, { content }),
    delete: (id)         => del(`/memories/${id}`),
    search: (query, limit = 10) => post('/memories/search', { query, limit }),
  },

  security: {
    audit: (params) => get('/security/audit' + (params ? '?' + new URLSearchParams(params) : '')),
  },
}

// JSON-RPC client for /rpc endpoint
async function rpc(method, params) {
  const res = await fetch('/rpc', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ method, params }),
  })
  const json = await res.json()
  if (json.error) throw new Error(json.error)
  return json.result
}

export const cronApi = {
  status:  ()           => rpc('cron.status'),
  list:    (params)     => rpc('cron.list', params ?? {}),
  add:     (params)     => rpc('cron.add', params),
  update:  (id, patch)  => rpc('cron.update', { id, patch }),
  remove:  (id)         => rpc('cron.remove', { id }),
  run:     (id, mode)   => rpc('cron.run', { id, mode: mode ?? 'force' }),
  runs:    (params)     => rpc('cron.runs', params ?? {}),
}

export default api
