import React, { useState, useEffect, useCallback } from 'react'
import { cronApi } from '../api/gateway.js'

function fmtMs(ms) {
  if (!ms) return '—'
  return new Date(ms).toLocaleString()
}
function fmtDuration(ms) {
  if (!ms) return '—'
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
}
function scheduleLabel(s) {
  if (!s) return '—'
  if (s.kind === 'at')    return `at ${s.at}`
  if (s.kind === 'every') return `every ${s.every_ms / 1000}s`
  if (s.kind === 'cron')  return `cron: ${s.expr}${s.tz ? ` (${s.tz})` : ''}`
  return s.kind
}

function StatusChip({ status }) {
  const map = {
    ok:      'bg-tertiary-container/20 text-tertiary',
    error:   'bg-error-container/20 text-error',
    running: 'bg-secondary/10 text-secondary',
    skipped: 'bg-surface-variant text-on-surface-variant',
  }
  const dot = { ok: 'bg-tertiary pulse-dot', error: 'bg-error', running: 'bg-secondary pulse-dot', skipped: 'bg-on-surface-variant' }
  const cls = map[status] ?? 'bg-surface-variant text-on-surface-variant'
  return (
    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wide ${cls}`}>
      <span className={`w-1 h-1 rounded-full ${dot[status] ?? 'bg-on-surface-variant'}`} />
      {status || '—'}
    </span>
  )
}

function ScheduleForm({ value, onChange }) {
  const set = (k, v) => onChange({ ...value, [k]: v })
  return (
    <div className="space-y-3">
      <label className="flex flex-col gap-1.5">
        <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Schedule Type</span>
        <select className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
          style={{ border: '1px solid rgba(70,69,84,0.2)' }}
          value={value.kind} onChange={e => onChange({ kind: e.target.value })}>
          <option value="every">Interval (Every)</option>
          <option value="cron">Cron Expression</option>
          <option value="at">Specific Date (At)</option>
        </select>
      </label>
      {value.kind === 'every' && (
        <label className="flex flex-col gap-1.5">
          <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Interval (seconds)</span>
          <input type="number" min={1}
            className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}
            value={(value.every_ms ?? 60000) / 1000}
            onChange={e => set('every_ms', Number(e.target.value) * 1000)} />
        </label>
      )}
      {value.kind === 'cron' && (
        <>
          <label className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Cron Expression</span>
            <input className="bg-surface-container-lowest rounded px-3 py-2 text-sm text-secondary font-mono focus:ring-1 focus:ring-primary focus:outline-none"
              style={{ border: '1px solid rgba(70,69,84,0.2)' }}
              placeholder="0 * * * *" value={value.expr ?? ''}
              onChange={e => set('expr', e.target.value)} />
          </label>
          <label className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Timezone</span>
            <input className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
              style={{ border: '1px solid rgba(70,69,84,0.2)' }}
              placeholder="Asia/Shanghai" value={value.tz ?? ''}
              onChange={e => set('tz', e.target.value)} />
          </label>
        </>
      )}
      {value.kind === 'at' && (
        <label className="flex flex-col gap-1.5">
          <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Run At (RFC3339)</span>
          <input className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}
            placeholder="2099-12-31T23:59:59Z" value={value.at ?? ''}
            onChange={e => set('at', e.target.value)} />
        </label>
      )}
    </div>
  )
}

function PayloadForm({ value, onChange, sessionTarget }) {
  const set = (k, v) => onChange({ ...value, [k]: v })
  return (
    <div className="space-y-3">
      <label className="flex flex-col gap-1.5">
        <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Payload Type</span>
        <select className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
          style={{ border: '1px solid rgba(70,69,84,0.2)' }}
          value={value.kind} onChange={e => onChange({ kind: e.target.value })}>
          {sessionTarget === 'main'
            ? <option value="systemEvent">systemEvent</option>
            : <option value="agentTurn">agentTurn</option>}
        </select>
      </label>
      {value.kind === 'systemEvent' && (
        <label className="flex flex-col gap-1.5">
          <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Event Text</span>
          <input className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}
            value={value.text ?? ''} onChange={e => set('text', e.target.value)} />
        </label>
      )}
      {value.kind === 'agentTurn' && (
        <>
          <label className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Message</span>
            <textarea rows={3} className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none resize-none"
              style={{ border: '1px solid rgba(70,69,84,0.2)' }}
              value={value.message ?? ''} onChange={e => set('message', e.target.value)} />
          </label>
          <label className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Timeout (seconds)</span>
            <input type="number" min={1} className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
              style={{ border: '1px solid rgba(70,69,84,0.2)' }}
              value={value.timeout_seconds ?? 60} onChange={e => set('timeout_seconds', Number(e.target.value))} />
          </label>
          <label className="flex items-center gap-2 text-sm text-on-surface-variant cursor-pointer">
            <input type="checkbox" className="accent-secondary" checked={!!value.light_context}
              onChange={e => set('light_context', e.target.checked)} />
            Light context
          </label>
        </>
      )}
    </div>
  )
}

const DEFAULT_FORM = {
  name: '', description: '', session_target: 'isolated', agent_id: '',
  wake_mode: 'next-heartbeat', enabled: true, delete_after_run: false,
  schedule: { kind: 'every', every_ms: 60000 },
  payload: { kind: 'agentTurn', message: '', timeout_seconds: 60 },
}

function AddJobModal({ onClose, onCreated }) {
  const [form, setForm] = useState(DEFAULT_FORM)
  const [saving, setSaving] = useState(false)
  const [err, setErr] = useState(null)
  const setField = (k, v) => setForm(f => ({ ...f, [k]: v }))

  useEffect(() => {
    if (form.session_target === 'main') setField('payload', { kind: 'systemEvent', text: '' })
    else if (form.payload.kind === 'systemEvent') setField('payload', { kind: 'agentTurn', message: '', timeout_seconds: 60 })
  }, [form.session_target])

  async function submit(e) {
    e.preventDefault(); setErr(null); setSaving(true)
    try {
      await cronApi.add({ name: form.name, description: form.description || undefined,
        session_target: form.session_target, agent_id: form.agent_id || undefined,
        wake_mode: form.wake_mode, enabled: form.enabled,
        delete_after_run: form.delete_after_run || undefined,
        schedule: form.schedule, payload: form.payload })
      onCreated()
    } catch (e) { setErr(e.message) }
    finally { setSaving(false) }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ background: 'rgba(11,19,38,0.85)', backdropFilter: 'blur(8px)' }}>
      <div className="bg-surface-container w-full max-w-2xl rounded-2xl overflow-hidden shadow-2xl max-h-[90vh] flex flex-col" style={{ border: '1px solid rgba(70,69,84,0.3)' }}>
        <div className="bg-surface-container-high px-8 py-6 flex justify-between items-center shrink-0" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
          <h3 className="font-headline text-xl font-bold">New Task Configuration</h3>
          <button onClick={onClose} className="text-on-surface-variant hover:text-on-surface transition-colors">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>
        <form onSubmit={submit} className="p-8 overflow-y-auto flex-1">
          <div className="grid grid-cols-2 gap-5 mb-6">
            <label className="col-span-2 flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Name *</span>
              <input required className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                value={form.name} onChange={e => setField('name', e.target.value)} />
            </label>
            <label className="col-span-2 flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Description</span>
              <input className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                value={form.description} onChange={e => setField('description', e.target.value)} />
            </label>
            <label className="flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Session Target</span>
              <select className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                value={form.session_target} onChange={e => setField('session_target', e.target.value)}>
                <option value="isolated">isolated</option>
                <option value="current">current</option>
                <option value="main">main</option>
              </select>
            </label>
            <label className="flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Agent ID</span>
              <input className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                placeholder="agent-001" value={form.agent_id} onChange={e => setField('agent_id', e.target.value)} />
            </label>
            <label className="flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Wake Mode</span>
              <select className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                value={form.wake_mode} onChange={e => setField('wake_mode', e.target.value)}>
                <option value="next-heartbeat">next-heartbeat</option>
                <option value="now">now</option>
              </select>
            </label>
            <div className="flex flex-col gap-3 pt-5">
              <label className="flex items-center gap-2 text-sm text-on-surface-variant cursor-pointer">
                <input type="checkbox" className="accent-secondary" checked={form.enabled}
                  onChange={e => setField('enabled', e.target.checked)} /> Enabled
              </label>
              <label className="flex items-center gap-2 text-sm text-on-surface-variant cursor-pointer">
                <input type="checkbox" className="accent-secondary" checked={form.delete_after_run}
                  onChange={e => setField('delete_after_run', e.target.checked)} /> Delete after run
              </label>
            </div>
          </div>
          <div className="bg-surface-container-low rounded-xl p-5 mb-5" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
            <p className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant mb-4">Schedule</p>
            <ScheduleForm value={form.schedule} onChange={v => setField('schedule', v)} />
          </div>
          <div className="bg-surface-container-low rounded-xl p-5 mb-5" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
            <p className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant mb-4">Payload</p>
            <PayloadForm value={form.payload} onChange={v => setField('payload', v)} sessionTarget={form.session_target} />
          </div>
          {err && <p className="text-error text-xs mb-4">{err}</p>}
          <div className="flex justify-end gap-3">
            <button type="button" onClick={onClose}
              className="px-5 py-2 text-sm font-bold text-on-surface-variant hover:text-on-surface transition-colors">Discard</button>
            <button type="submit" disabled={saving}
              className="px-8 py-2 text-sm font-bold rounded text-on-primary disabled:opacity-50 transition-all hover:brightness-110"
              style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
              {saving ? 'Deploying…' : 'Deploy Task'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

function RunLogsModal({ job, onClose }) {
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
    <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ background: 'rgba(11,19,38,0.85)', backdropFilter: 'blur(8px)' }}>
      <div className="bg-surface-container w-full max-w-3xl rounded-2xl shadow-2xl max-h-[85vh] flex flex-col" style={{ border: '1px solid rgba(70,69,84,0.3)' }}>
        <div className="px-8 py-6 flex items-center justify-between shrink-0" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
          <div>
            <h3 className="font-headline font-bold text-lg">Execution Logs</h3>
            <p className="text-xs text-on-surface-variant mt-0.5">{job.name} · {total} total runs</p>
          </div>
          <button onClick={onClose} className="text-on-surface-variant hover:text-on-surface transition-colors">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>
        <div className="overflow-auto flex-1 p-6">
          {loading ? (
            <div className="text-on-surface/40 text-sm flex items-center gap-2">
              <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
            </div>
          ) : logs.length === 0 ? (
            <div className="text-on-surface-variant text-sm text-center py-8">No runs yet.</div>
          ) : (
            <table className="w-full text-xs">
              <thead className="bg-surface-container-low/50 sticky top-0">
                <tr>
                  {['Time', 'Status', 'Duration', 'Delivery', 'Error'].map(h => (
                    <th key={h} className="text-left px-4 py-2.5 text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y" style={{ borderColor: 'rgba(70,69,84,0.05)' }}>
                {logs.map(l => (
                  <tr key={l.id} className="hover:bg-surface-container/40 transition-colors">
                    <td className="px-4 py-3 text-on-surface-variant whitespace-nowrap">{fmtMs(l.run_at_ms)}</td>
                    <td className="px-4 py-3"><StatusChip status={l.status} /></td>
                    <td className="px-4 py-3 text-on-surface-variant font-mono">{fmtDuration(l.duration_ms)}</td>
                    <td className="px-4 py-3"><StatusChip status={l.delivery_status} /></td>
                    <td className="px-4 py-3 text-error max-w-xs truncate">{l.error || '—'}</td>
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

function JobRow({ job, onRefresh, onViewLogs }) {
  const [running, setRunning]   = useState(false)
  const [toggling, setToggling] = useState(false)
  const isRunningNow = !!job.state?.running_at_ms

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

  return (
    <tr className="hover:bg-surface-container/40 transition-colors group">
      <td className="px-6 py-5">
        <div className="font-bold text-on-surface text-sm">{job.name}</div>
        {job.description && <div className="text-[10px] text-on-surface-variant/50 font-mono mt-0.5 uppercase tracking-tighter">{job.id?.slice(0, 16)}</div>}
      </td>
      <td className="px-6 py-5">
        <code className="text-xs bg-surface-container-high px-2 py-1 rounded text-secondary font-mono">{scheduleLabel(job.schedule)}</code>
      </td>
      <td className="px-6 py-5 text-xs text-on-surface-variant">{fmtMs(job.state?.next_run_at_ms)}</td>
      <td className="px-6 py-5">
        <StatusChip status={isRunningNow ? 'running' : job.state?.last_run_status} />
        {job.state?.consecutive_errors > 0 && (
          <span className="ml-1 text-[10px] text-error">×{job.state.consecutive_errors}</span>
        )}
      </td>
      <td className="px-6 py-5">
        <button onClick={handleToggle} disabled={toggling}
          className={`w-10 h-5 rounded-full relative transition-colors ${job.enabled ? 'bg-secondary/30' : 'bg-outline-variant/20'}`}>
          <div className={`absolute top-1 w-3 h-3 rounded-full transition-all ${job.enabled ? 'right-1 bg-secondary' : 'left-1 bg-outline-variant'}`} />
        </button>
      </td>
      <td className="px-6 py-5 text-right">
        <div className="flex gap-1.5 justify-end opacity-0 group-hover:opacity-100 transition-opacity">
          <button onClick={() => onViewLogs(job)}
            className="text-xs px-2.5 py-1 rounded text-on-surface-variant hover:bg-surface-container-high transition-colors"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}>Logs</button>
          <button onClick={handleRun} disabled={running || isRunningNow}
            className="text-xs px-2.5 py-1 rounded text-secondary hover:bg-secondary/10 transition-colors disabled:opacity-40"
            style={{ border: '1px solid rgba(76,215,246,0.2)' }}>
            {running ? '…' : 'Run'}
          </button>
          <button onClick={handleDelete}
            className="text-xs px-2.5 py-1 rounded text-error hover:bg-error-container/10 transition-colors"
            style={{ border: '1px solid rgba(255,180,171,0.2)' }}>Del</button>
        </div>
      </td>
    </tr>
  )
}

export default function Cron() {
  const [jobs, setJobs]       = useState([])
  const [status, setStatus]   = useState(null)
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  const [logsJob, setLogsJob] = useState(null)
  const [includeDisabled, setIncludeDisabled] = useState(true)
  const [activeTab, setActiveTab] = useState('jobs')

  const load = useCallback(() => {
    setLoading(true)
    Promise.all([
      cronApi.list({ include_disabled: includeDisabled, limit: 100 }),
      cronApi.status(),
    ])
      .then(([listRes, statusRes]) => { setJobs(listRes.jobs ?? []); setStatus(statusRes) })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [includeDisabled])

  useEffect(() => { load() }, [load])
  useEffect(() => { const t = setInterval(load, 10000); return () => clearInterval(t) }, [load])

  return (
    <div className="p-10 max-w-7xl">
      {/* Header */}
      <div className="flex items-end justify-between mb-10">
        <div>
          <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">Task Orchestration</h2>
          <p className="text-on-surface-variant text-sm mt-1 max-w-md">Manage autonomous agent triggers and system event intervals</p>
        </div>
        <div className="flex gap-3">
          <button onClick={load}
            className="px-5 py-2.5 text-sm font-semibold rounded text-on-surface-variant hover:bg-surface-container transition-all flex items-center gap-2"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
            <span className="material-symbols-outlined text-sm">refresh</span> Refresh
          </button>
          <button onClick={() => setShowAdd(true)}
            className="px-6 py-2.5 text-sm font-bold rounded text-on-primary flex items-center gap-2 shadow-lg transition-all hover:brightness-110"
            style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
            <span className="material-symbols-outlined text-sm">add</span> Create Task
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10">
        {[
          { label: 'Total Tasks',   value: jobs.length,                    color: 'text-on-surface' },
          { label: 'Enabled',       value: status?.enabled_jobs ?? '…',    color: 'text-secondary' },
          { label: 'Running Now',   value: status?.running_jobs ?? 0,      color: 'text-tertiary' },
          { label: 'Next Run',      value: status?.next_run_at_ms ? fmtMs(status.next_run_at_ms) : '—', color: 'text-primary', small: true },
        ].map(({ label, value, color, small }) => (
          <div key={label} className="bg-surface-container-low p-6 rounded-xl" style={{ border: '1px solid rgba(70,69,84,0.05)' }}>
            <p className="text-[10px] uppercase tracking-widest text-on-surface-variant font-bold mb-2">{label}</p>
            <h3 className={`font-headline font-bold ${small ? 'text-lg mt-2' : 'text-4xl'} ${color}`}>{value}</h3>
          </div>
        ))}
      </div>

      {/* Tabs */}
      <div className="flex gap-8 mb-6" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
        {[['jobs', 'Jobs List'], ['logs', 'Execution Logs']].map(([id, label]) => (
          <button key={id} onClick={() => setActiveTab(id)}
            className={`pb-4 text-sm font-bold transition-colors ${activeTab === id ? 'border-b-2 border-primary text-on-surface' : 'text-on-surface-variant/60 hover:text-on-surface'}`}>
            {label}
          </button>
        ))}
        <label className="ml-auto flex items-center gap-2 text-xs text-on-surface-variant cursor-pointer pb-4">
          <input type="checkbox" className="accent-secondary" checked={includeDisabled}
            onChange={e => setIncludeDisabled(e.target.checked)} />
          Show disabled
        </label>
      </div>

      {activeTab === 'jobs' && (
        loading ? (
          <div className="text-on-surface/40 text-sm flex items-center gap-2">
            <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
          </div>
        ) : jobs.length === 0 ? (
          <div className="bg-surface-container-low rounded-xl p-16 text-center">
            <span className="material-symbols-outlined text-on-surface-variant/30 mb-4" style={{ fontSize: '3rem' }}>schedule</span>
            <p className="text-on-surface-variant text-sm">No cron jobs yet. Create one to get started.</p>
          </div>
        ) : (
          <div className="bg-surface-container-lowest rounded-xl overflow-x-auto" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
            <table className="w-full text-sm">
              <thead className="bg-surface-container-low/50">
                <tr>
                  {['Job Name', 'Schedule', 'Next Run', 'Last Status', 'Enabled', ''].map(h => (
                    <th key={h} className="text-left px-6 py-3 text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y" style={{ borderColor: 'rgba(70,69,84,0.05)' }}>
                {jobs.map(job => (
                  <JobRow key={job.id} job={job} onRefresh={load} onViewLogs={setLogsJob} />
                ))}
              </tbody>
            </table>
          </div>
        )
      )}

      {activeTab === 'logs' && <AllLogsTab />}

      {showAdd && <AddJobModal onClose={() => setShowAdd(false)} onCreated={() => { setShowAdd(false); load() }} />}
      {logsJob && <RunLogsModal job={logsJob} onClose={() => setLogsJob(null)} />}
    </div>
  )
}

function AllLogsTab() {
  const [logs, setLogs] = useState([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    cronApi.runs({ limit: 30, sort_dir: 'DESC' })
      .then(r => { setLogs(r.logs ?? []); setTotal(r.total ?? 0) })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-on-surface/40 text-sm flex items-center gap-2"><span className="material-symbols-outlined animate-spin">autorenew</span> Loading…</div>

  return (
    <div>
      <div className="flex items-center gap-2 mb-4">
        <span className="w-1.5 h-1.5 rounded-full bg-secondary" />
        <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Live Execution Stream · {total} total</span>
      </div>
      <div className="bg-surface-container-lowest rounded-xl p-4 font-mono text-xs leading-relaxed max-h-96 overflow-y-auto relative"
        style={{ border: '1px solid rgba(70,69,84,0.2)', borderLeft: '2px solid rgba(76,215,246,0.3)' }}>
        {logs.length === 0 ? (
          <p className="text-on-surface-variant/40">No execution logs yet.</p>
        ) : (
          <div className="space-y-2 pl-2">
            {logs.map(l => (
              <div key={l.id} className="flex gap-3 text-on-surface-variant/60">
                <span className="text-on-surface-variant/30 shrink-0">[{fmtMs(l.run_at_ms)}]</span>
                <span className={l.status === 'ok' ? 'text-tertiary' : l.status === 'error' ? 'text-error' : 'text-secondary'}>
                  {l.status?.toUpperCase()}:
                </span>
                <span>{l.job_id}</span>
                {l.duration_ms && <span className="text-on-surface-variant/30">{fmtDuration(l.duration_ms)}</span>}
                {l.error && <span className="text-error truncate max-w-xs">{l.error}</span>}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
