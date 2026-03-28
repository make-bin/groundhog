import React, { useState, useEffect, useCallback } from 'react'
import { cronApi } from '../api/gateway.js'

// ── helpers ──────────────────────────────────────────────────────────────────

function fmtMs(ms) {
  if (!ms) return '—'
  return new Date(ms).toLocaleString()
}

function fmtDuration(ms) {
  if (!ms) return '—'
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
}

function StatusBadge({ status }) {
  const map = {
    ok:      'bg-green-100 text-green-700',
    error:   'bg-red-100 text-red-600',
    running: 'bg-blue-100 text-blue-700',
    skipped: 'bg-gray-100 text-gray-500',
  }
  return (
    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${map[status] ?? 'bg-gray-100 text-gray-500'}`}>
      {status || '—'}
    </span>
  )
}

function EnabledBadge({ enabled }) {
  return (
    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
      enabled ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-400'
    }`}>
      {enabled ? 'enabled' : 'disabled'}
    </span>
  )
}

function scheduleLabel(s) {
  if (!s) return '—'
  if (s.kind === 'at')    return `at ${s.at}`
  if (s.kind === 'every') return `every ${s.every_ms / 1000}s`
  if (s.kind === 'cron')  return `cron: ${s.expr}${s.tz ? ` (${s.tz})` : ''}`
  return s.kind
}

// ── ScheduleForm ─────────────────────────────────────────────────────────────

function ScheduleForm({ value, onChange }) {
  const set = (k, v) => onChange({ ...value, [k]: v })
  return (
    <div className="flex flex-col gap-2">
      <label className="text-xs text-gray-600">
        Kind
        <select className="ml-2 border rounded px-2 py-1 text-sm"
          value={value.kind} onChange={e => onChange({ kind: e.target.value })}>
          <option value="every">every</option>
          <option value="cron">cron</option>
          <option value="at">at (one-shot)</option>
        </select>
      </label>
      {value.kind === 'every' && (
        <label className="text-xs text-gray-600">
          Interval (seconds)
          <input type="number" min={1} className="ml-2 border rounded px-2 py-1 text-sm w-28"
            value={(value.every_ms ?? 60000) / 1000}
            onChange={e => set('every_ms', Number(e.target.value) * 1000)} />
        </label>
      )}
      {value.kind === 'cron' && (
        <>
          <label className="text-xs text-gray-600">
            Cron expr
            <input className="ml-2 border rounded px-2 py-1 text-sm w-40"
              placeholder="* * * * *" value={value.expr ?? ''}
              onChange={e => set('expr', e.target.value)} />
          </label>
          <label className="text-xs text-gray-600">
            Timezone
            <input className="ml-2 border rounded px-2 py-1 text-sm w-40"
              placeholder="Asia/Shanghai" value={value.tz ?? ''}
              onChange={e => set('tz', e.target.value)} />
          </label>
        </>
      )}
      {value.kind === 'at' && (
        <label className="text-xs text-gray-600">
          Run at (RFC3339)
          <input className="ml-2 border rounded px-2 py-1 text-sm w-56"
            placeholder="2099-12-31T23:59:59Z" value={value.at ?? ''}
            onChange={e => set('at', e.target.value)} />
        </label>
      )}
    </div>
  )
}

// ── PayloadForm ───────────────────────────────────────────────────────────────

function PayloadForm({ value, onChange, sessionTarget }) {
  const set = (k, v) => onChange({ ...value, [k]: v })
  const isMain = sessionTarget === 'main'
  return (
    <div className="flex flex-col gap-2">
      <label className="text-xs text-gray-600">
        Kind
        <select className="ml-2 border rounded px-2 py-1 text-sm"
          value={value.kind}
          onChange={e => onChange({ kind: e.target.value })}>
          {isMain
            ? <option value="systemEvent">systemEvent</option>
            : <option value="agentTurn">agentTurn</option>}
        </select>
      </label>
      {value.kind === 'systemEvent' && (
        <label className="text-xs text-gray-600">
          Text
          <input className="ml-2 border rounded px-2 py-1 text-sm w-64"
            value={value.text ?? ''} onChange={e => set('text', e.target.value)} />
        </label>
      )}
      {value.kind === 'agentTurn' && (
        <>
          <label className="text-xs text-gray-600">
            Message
            <textarea rows={2} className="ml-2 border rounded px-2 py-1 text-sm w-64 resize-none align-top"
              value={value.message ?? ''} onChange={e => set('message', e.target.value)} />
          </label>
          <label className="text-xs text-gray-600">
            Timeout (s)
            <input type="number" min={1} className="ml-2 border rounded px-2 py-1 text-sm w-20"
              value={value.timeout_seconds ?? 60}
              onChange={e => set('timeout_seconds', Number(e.target.value))} />
          </label>
          <label className="text-xs text-gray-600 flex items-center gap-1">
            <input type="checkbox" checked={!!value.light_context}
              onChange={e => set('light_context', e.target.checked)} />
            Light context
          </label>
        </>
      )}
    </div>
  )
}

// ── AddJobModal ───────────────────────────────────────────────────────────────

const DEFAULT_FORM = {
  name: '',
  description: '',
  session_target: 'isolated',
  agent_id: '',
  wake_mode: 'next-heartbeat',
  enabled: true,
  delete_after_run: false,
  schedule: { kind: 'every', every_ms: 60000 },
  payload: { kind: 'agentTurn', message: '', timeout_seconds: 60 },
}

function AddJobModal({ onClose, onCreated }) {
  const [form, setForm] = useState(DEFAULT_FORM)
  const [saving, setSaving] = useState(false)
  const [err, setErr] = useState(null)

  const setField = (k, v) => setForm(f => ({ ...f, [k]: v }))

  // keep payload kind in sync with session_target
  useEffect(() => {
    if (form.session_target === 'main') {
      setField('payload', { kind: 'systemEvent', text: '' })
    } else if (form.payload.kind === 'systemEvent') {
      setField('payload', { kind: 'agentTurn', message: '', timeout_seconds: 60 })
    }
  }, [form.session_target])

  async function submit(e) {
    e.preventDefault()
    setErr(null)
    setSaving(true)
    try {
      const req = {
        name: form.name,
        description: form.description || undefined,
        session_target: form.session_target,
        agent_id: form.agent_id || undefined,
        wake_mode: form.wake_mode,
        enabled: form.enabled,
        delete_after_run: form.delete_after_run || undefined,
        schedule: form.schedule,
        payload: form.payload,
      }
      const job = await cronApi.add(req)
      onCreated(job)
    } catch (e) {
      setErr(e.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-lg p-6 overflow-y-auto max-h-[90vh]">
        <h2 className="font-semibold text-gray-800 mb-4">New Cron Job</h2>
        <form onSubmit={submit} className="flex flex-col gap-3">
          <label className="flex flex-col gap-1 text-xs text-gray-600">
            Name *
            <input required className="border rounded px-3 py-1.5 text-sm"
              value={form.name} onChange={e => setField('name', e.target.value)} />
          </label>
          <label className="flex flex-col gap-1 text-xs text-gray-600">
            Description
            <input className="border rounded px-3 py-1.5 text-sm"
              value={form.description} onChange={e => setField('description', e.target.value)} />
          </label>
          <div className="grid grid-cols-2 gap-3">
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              Session Target *
              <select className="border rounded px-3 py-1.5 text-sm"
                value={form.session_target} onChange={e => setField('session_target', e.target.value)}>
                <option value="isolated">isolated</option>
                <option value="current">current</option>
                <option value="main">main</option>
              </select>
            </label>
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              Agent ID
              <input className="border rounded px-3 py-1.5 text-sm"
                placeholder="agent-001" value={form.agent_id}
                onChange={e => setField('agent_id', e.target.value)} />
            </label>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              Wake Mode
              <select className="border rounded px-3 py-1.5 text-sm"
                value={form.wake_mode} onChange={e => setField('wake_mode', e.target.value)}>
                <option value="next-heartbeat">next-heartbeat</option>
                <option value="now">now</option>
              </select>
            </label>
            <div className="flex flex-col gap-2 pt-4">
              <label className="flex items-center gap-2 text-xs text-gray-600">
                <input type="checkbox" checked={form.enabled}
                  onChange={e => setField('enabled', e.target.checked)} />
                Enabled
              </label>
              <label className="flex items-center gap-2 text-xs text-gray-600">
                <input type="checkbox" checked={form.delete_after_run}
                  onChange={e => setField('delete_after_run', e.target.checked)} />
                Delete after run
              </label>
            </div>
          </div>

          <div className="border rounded-lg p-3 bg-gray-50">
            <div className="text-xs font-medium text-gray-600 mb-2">Schedule</div>
            <ScheduleForm value={form.schedule} onChange={v => setField('schedule', v)} />
          </div>

          <div className="border rounded-lg p-3 bg-gray-50">
            <div className="text-xs font-medium text-gray-600 mb-2">Payload</div>
            <PayloadForm value={form.payload} onChange={v => setField('payload', v)}
              sessionTarget={form.session_target} />
          </div>

          {err && <p className="text-red-500 text-xs">{err}</p>}

          <div className="flex gap-2 pt-1">
            <button type="submit" disabled={saving}
              className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm hover:bg-blue-700 disabled:opacity-50">
              {saving ? 'Creating…' : 'Create'}
            </button>
            <button type="button" onClick={onClose}
              className="px-4 py-1.5 rounded text-sm border hover:bg-gray-50">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ── RunLogsPanel ──────────────────────────────────────────────────────────────

function RunLogsPanel({ job, onClose }) {
  const [logs, setLogs] = useState([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    cronApi.runs({ scope: 'job', job_id: job.id, limit: 20, sort_dir: 'DESC' })
      .then(r => { setLogs(r.logs ?? []); setTotal(r.total ?? 0) })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [job.id])

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-3xl p-6 max-h-[85vh] flex flex-col">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="font-semibold text-gray-800">Run Logs</h2>
            <p className="text-xs text-gray-500 mt-0.5">{job.name} · {total} total</p>
          </div>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-xl leading-none">×</button>
        </div>
        <div className="overflow-auto flex-1">
          {loading ? (
            <div className="text-gray-400 text-sm p-4">Loading…</div>
          ) : logs.length === 0 ? (
            <div className="text-gray-400 text-sm p-4">No runs yet.</div>
          ) : (
            <table className="w-full text-xs">
              <thead className="bg-gray-50 text-gray-500 sticky top-0">
                <tr>
                  <th className="text-left px-3 py-2">Time</th>
                  <th className="text-left px-3 py-2">Status</th>
                  <th className="text-left px-3 py-2">Duration</th>
                  <th className="text-left px-3 py-2">Delivery</th>
                  <th className="text-left px-3 py-2">Error</th>
                </tr>
              </thead>
              <tbody>
                {logs.map(l => (
                  <tr key={l.id} className="border-t border-gray-50 hover:bg-gray-50">
                    <td className="px-3 py-2 text-gray-500 whitespace-nowrap">{fmtMs(l.run_at_ms)}</td>
                    <td className="px-3 py-2"><StatusBadge status={l.status} /></td>
                    <td className="px-3 py-2 text-gray-500">{fmtDuration(l.duration_ms)}</td>
                    <td className="px-3 py-2"><StatusBadge status={l.delivery_status} /></td>
                    <td className="px-3 py-2 text-red-500 max-w-xs truncate">{l.error || '—'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  )
}

// ── JobRow ────────────────────────────────────────────────────────────────────

function JobRow({ job, onRefresh, onViewLogs }) {
  const [running, setRunning] = useState(false)
  const [toggling, setToggling] = useState(false)

  async function handleRun() {
    setRunning(true)
    try { await cronApi.run(job.id, 'force') } catch (e) { alert(e.message) }
    finally { setRunning(false); onRefresh() }
  }

  async function handleToggle() {
    setToggling(true)
    try { await cronApi.update(job.id, { enabled: !job.enabled }) } catch (e) { alert(e.message) }
    finally { setToggling(false); onRefresh() }
  }

  async function handleDelete() {
    if (!confirm(`Delete "${job.name}"?`)) return
    try { await cronApi.remove(job.id) } catch (e) { alert(e.message) }
    onRefresh()
  }

  const isRunningNow = !!job.state?.running_at_ms

  return (
    <tr className="border-t border-gray-50 hover:bg-gray-50 text-sm">
      <td className="px-4 py-3">
        <div className="font-medium text-gray-800">{job.name}</div>
        {job.description && <div className="text-xs text-gray-400 mt-0.5">{job.description}</div>}
      </td>
      <td className="px-4 py-3 text-xs text-gray-500">{scheduleLabel(job.schedule)}</td>
      <td className="px-4 py-3 text-xs text-gray-500">{job.session_target}</td>
      <td className="px-4 py-3"><EnabledBadge enabled={job.enabled} /></td>
      <td className="px-4 py-3">
        <StatusBadge status={isRunningNow ? 'running' : job.state?.last_run_status} />
        {job.state?.consecutive_errors > 0 && (
          <span className="ml-1 text-xs text-red-400">×{job.state.consecutive_errors}</span>
        )}
      </td>
      <td className="px-4 py-3 text-xs text-gray-400">{fmtMs(job.state?.next_run_at_ms)}</td>
      <td className="px-4 py-3 text-xs text-gray-400">{fmtMs(job.state?.last_run_at_ms)}</td>
      <td className="px-4 py-3">
        <div className="flex gap-1.5 justify-end">
          <button onClick={() => onViewLogs(job)}
            className="text-xs px-2 py-1 rounded border hover:bg-gray-100 text-gray-600">
            Logs
          </button>
          <button onClick={handleRun} disabled={running || isRunningNow}
            className="text-xs px-2 py-1 rounded border hover:bg-blue-50 text-blue-600 disabled:opacity-40">
            {running ? '…' : 'Run'}
          </button>
          <button onClick={handleToggle} disabled={toggling}
            className="text-xs px-2 py-1 rounded border hover:bg-gray-100 text-gray-600 disabled:opacity-40">
            {job.enabled ? 'Disable' : 'Enable'}
          </button>
          <button onClick={handleDelete}
            className="text-xs px-2 py-1 rounded border hover:bg-red-50 text-red-500">
            Delete
          </button>
        </div>
      </td>
    </tr>
  )
}

// ── Main Page ─────────────────────────────────────────────────────────────────

export default function Cron() {
  const [jobs, setJobs]           = useState([])
  const [status, setStatus]       = useState(null)
  const [loading, setLoading]     = useState(true)
  const [showAdd, setShowAdd]     = useState(false)
  const [logsJob, setLogsJob]     = useState(null)
  const [includeDisabled, setIncludeDisabled] = useState(true)

  const load = useCallback(() => {
    setLoading(true)
    Promise.all([
      cronApi.list({ include_disabled: includeDisabled, limit: 100 }),
      cronApi.status(),
    ])
      .then(([listRes, statusRes]) => {
        setJobs(listRes.jobs ?? [])
        setStatus(statusRes)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [includeDisabled])

  useEffect(() => { load() }, [load])

  // auto-refresh every 10s
  useEffect(() => {
    const t = setInterval(load, 10000)
    return () => clearInterval(t)
  }, [load])

  return (
    <div className="p-8 max-w-7xl">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">Cron Jobs</h1>
          <p className="text-sm text-gray-500 mt-0.5">Scheduled tasks managed by the cron scheduler</p>
        </div>
        <button onClick={() => setShowAdd(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-blue-700 transition-colors">
          + New Job
        </button>
      </div>

      {/* Scheduler status bar */}
      {status && (
        <div className="flex gap-3 mb-6 flex-wrap">
          {[
            { label: 'Scheduler', ok: status.running, val: status.running ? 'running' : 'stopped' },
            { label: 'Enabled jobs', ok: true, val: status.enabled_jobs },
            { label: 'Running now', ok: status.running_jobs === 0, val: status.running_jobs },
            { label: 'Next run', ok: true, val: status.next_run_at_ms ? fmtMs(status.next_run_at_ms) : '—' },
          ].map(({ label, ok, val }) => (
            <div key={label} className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium border ${
              ok ? 'bg-green-50 border-green-200 text-green-700' : 'bg-amber-50 border-amber-200 text-amber-700'
            }`}>
              <span>{label}</span>
              <span className="opacity-70">{val}</span>
            </div>
          ))}
          <button onClick={load}
            className="px-3 py-1.5 rounded-full text-xs border border-gray-200 text-gray-500 hover:bg-gray-50">
            ↻ Refresh
          </button>
        </div>
      )}

      {/* Filter */}
      <div className="flex items-center gap-3 mb-4">
        <label className="flex items-center gap-2 text-xs text-gray-600 cursor-pointer">
          <input type="checkbox" checked={includeDisabled}
            onChange={e => setIncludeDisabled(e.target.checked)} />
          Show disabled jobs
        </label>
        <span className="text-xs text-gray-400">{jobs.length} job{jobs.length !== 1 ? 's' : ''}</span>
      </div>

      {/* Table */}
      {loading ? (
        <div className="text-gray-400 text-sm">Loading…</div>
      ) : jobs.length === 0 ? (
        <div className="text-gray-400 text-sm">No cron jobs yet. Create one to get started.</div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-gray-500 text-xs">
              <tr>
                <th className="text-left px-4 py-2.5">Name</th>
                <th className="text-left px-4 py-2.5">Schedule</th>
                <th className="text-left px-4 py-2.5">Target</th>
                <th className="text-left px-4 py-2.5">Status</th>
                <th className="text-left px-4 py-2.5">Last Run</th>
                <th className="text-left px-4 py-2.5">Next Run</th>
                <th className="text-left px-4 py-2.5">Last Run At</th>
                <th className="px-4 py-2.5" />
              </tr>
            </thead>
            <tbody>
              {jobs.map(job => (
                <JobRow key={job.id} job={job} onRefresh={load} onViewLogs={setLogsJob} />
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showAdd && (
        <AddJobModal
          onClose={() => setShowAdd(false)}
          onCreated={() => { setShowAdd(false); load() }}
        />
      )}

      {logsJob && (
        <RunLogsPanel job={logsJob} onClose={() => setLogsJob(null)} />
      )}
    </div>
  )
}
