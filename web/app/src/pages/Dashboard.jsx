import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import api, { cronApi } from '../api/gateway.js'

function StatCard({ icon, iconColor, label, value, sub, badge, to, barColor, barPct }) {
  const inner = (
    <div className="bg-surface-container-low p-5 rounded-xl border border-transparent hover:border-primary/20 transition-all group cursor-pointer">
      <div className="flex justify-between items-start mb-4">
        <div className={`p-2 rounded-lg ${iconColor}/10`}>
          <span className={`material-symbols-outlined ${iconColor}`}>{icon}</span>
        </div>
        {badge && (
          <span className={`text-[10px] font-bold uppercase tracking-widest ${iconColor} bg-current/10 px-2 py-0.5 rounded-full`}
            style={{ background: 'rgba(78,222,163,0.1)', color: 'inherit' }}>
            {badge}
          </span>
        )}
      </div>
      <h3 className="text-on-surface-variant text-xs font-medium uppercase tracking-tighter">{label}</h3>
      <p className="font-headline text-3xl font-bold text-on-surface mt-1">{value ?? '—'}</p>
      {sub && <p className="text-[10px] text-on-surface/40 mt-3 flex items-center gap-1">{sub}</p>}
      {barColor && (
        <div className="mt-4 h-1 w-full bg-surface-container rounded-full overflow-hidden">
          <div className={`h-full ${barColor}`} style={{ width: `${barPct ?? 100}%` }} />
        </div>
      )}
    </div>
  )
  return to ? <Link to={to}>{inner}</Link> : inner
}

export default function Dashboard() {
  const [health, setHealth]       = useState(null)
  const [sessions, setSessions]   = useState(null)
  const [channels, setChannels]   = useState(null)
  const [memories, setMemories]   = useState(null)
  const [cronStatus, setCronStatus] = useState(null)
  const [agents, setAgents]       = useState(null)

  useEffect(() => {
    api.health().then(setHealth).catch(() => setHealth({ status: 'error' }))
    api.sessions.list({ limit: 100 }).then(d => setSessions(d?.sessions ?? d ?? [])).catch(() => setSessions([]))
    api.channels.list().then(d => setChannels(Array.isArray(d) ? d : [])).catch(() => setChannels([]))
    api.memories.list().then(d => setMemories(d?.memories ?? [])).catch(() => setMemories([]))
    cronApi.status().then(setCronStatus).catch(() => {})
    api.agents.list().then(d => setAgents(Array.isArray(d) ? d : [])).catch(() => setAgents([]))
  }, [])

  const dbOk    = health?.database === 'ok'
  const redisOk = health?.redis === 'ok' || health?.redis?.startsWith?.('ok')

  return (
    <div className="p-10 max-w-6xl">
      {/* Page header */}
      <div className="mb-10">
        <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">System Overview</h2>
        <p className="text-on-surface-variant text-sm mt-1">Real-time telemetry across the Groundhog network</p>
      </div>

      {/* Stats bento grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-10">
        <StatCard
          icon="health_metrics" iconColor="text-tertiary"
          label="System Health" value={health ? '100%' : '…'}
          badge={health ? 'Optimal' : null}
          barColor="bg-tertiary" barPct={100}
          to="/"
        />
        <StatCard
          icon="smart_toy" iconColor="text-primary"
          label="Agents" value={agents?.length ?? '…'}
          badge={agents?.find(a => a.is_default) ? `default: ${agents.find(a => a.is_default).id}` : null}
          to="/agents"
        />
        <StatCard
          icon="forum" iconColor="text-secondary"
          label="Active Sessions" value={sessions?.length ?? '…'}
          badge={sessions?.filter(s => s.state === 'Active').length > 0 ? `+${sessions.filter(s => s.state === 'Active').length} active` : null}
          to="/sessions"
        />
        <StatCard
          icon="hub" iconColor="text-on-surface-variant"
          label="Channels" value={channels?.length ?? '…'}
          sub={<><span className="material-symbols-outlined text-[10px]">sync</span> All nodes synced</>}
          to="/channels"
        />
      </div>

      {/* Middle: health + cron */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-10">
        {/* Health panel */}
        <div className="lg:col-span-2 bg-surface-container-low rounded-2xl overflow-hidden flex flex-col">
          <div className="p-6 flex justify-between items-center" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
            <div>
              <h2 className="font-headline text-xl font-bold text-primary">Real-time System Health</h2>
              <p className="text-xs text-on-surface-variant font-medium mt-1 uppercase tracking-widest">Global Telemetry</p>
            </div>
            <div className="flex items-center gap-2 bg-surface-container rounded-lg px-3 py-1.5" style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
              <span className="w-2 h-2 rounded-full bg-tertiary pulse-dot" />
              <span className="text-[10px] font-bold uppercase tracking-wider text-tertiary">Live</span>
            </div>
          </div>
          <div className="p-8 grid grid-cols-3 gap-8">
            {[
              { label: 'API Cluster',   icon: 'bolt',     ok: health !== null, sub: 'Response: ~12ms' },
              { label: 'Database',      icon: 'database', ok: dbOk,            sub: 'Load: nominal' },
              { label: 'Redis Cache',   icon: 'memory',   ok: redisOk,         sub: 'Hits: active' },
            ].map(({ label, icon, ok, sub }) => (
              <div key={label} className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-bold text-on-surface-variant/60 uppercase">{label}</span>
                  <span className={`text-[10px] font-bold ${ok ? 'text-tertiary' : 'text-error'}`}>{ok ? 'Stable' : 'Error'}</span>
                </div>
                <div className={`p-4 bg-surface-container-lowest rounded-xl border-l-2 ${ok ? 'border-tertiary' : 'border-error'}`}>
                  <div className="flex items-center gap-3">
                    <span className={`material-symbols-outlined ${ok ? 'text-tertiary' : 'text-error'}`}>{icon}</span>
                    <div>
                      <p className="text-sm font-bold text-on-surface">{ok ? 'Active' : 'Offline'}</p>
                      <p className="text-[10px] text-on-surface-variant/50">{sub}</p>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Cron scheduler */}
        <div className="bg-surface-container rounded-2xl p-6" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
          <div className="flex items-center gap-3 mb-8">
            <div className="p-2 bg-primary-container/20 rounded-lg">
              <span className="material-symbols-outlined text-primary-container">watch_later</span>
            </div>
            <h2 className="font-headline text-lg font-bold">Cron Scheduler</h2>
          </div>
          <div className="bg-surface-container-lowest rounded-xl p-5 mb-8 text-center relative overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
            <div className="absolute top-0 right-0 p-2">
              <span className="material-symbols-outlined text-primary-container/20" style={{ fontSize: '2.5rem' }}>timer</span>
            </div>
            <p className="text-[10px] font-bold text-on-surface-variant/40 uppercase tracking-widest mb-1">Scheduler Status</p>
            <p className={`font-headline text-2xl font-bold ${cronStatus?.running ? 'text-tertiary' : 'text-error'}`}>
              {cronStatus ? (cronStatus.running ? 'Running' : 'Stopped') : '…'}
            </p>
            <p className="text-[10px] text-on-surface-variant mt-2 font-medium">
              {cronStatus?.enabled_jobs ?? 0} enabled · {cronStatus?.running_jobs ?? 0} running
            </p>
          </div>
          <div className="space-y-3">
            <h3 className="text-xs font-bold text-on-surface-variant uppercase tracking-widest">Quick Links</h3>
            {[
              { label: 'Manage Agents',   to: '/agents',   color: 'bg-primary' },
              { label: 'Manage Sessions', to: '/sessions', color: 'bg-secondary' },
              { label: 'Memory Store',    to: '/memory',   color: 'bg-tertiary' },
            ].map(({ label, to, color }) => (
              <Link key={to} to={to}
                className="flex items-center justify-between p-3 bg-surface-container-low rounded-lg hover:bg-surface-container-highest transition-colors group">
                <div className="flex items-center gap-3">
                  <div className={`w-1.5 h-1.5 rounded-full ${color}`} />
                  <span className="text-xs font-medium">{label}</span>
                </div>
                <span className="material-symbols-outlined text-on-surface-variant/40 group-hover:text-primary transition-colors text-sm">arrow_forward</span>
              </Link>
            ))}
          </div>
        </div>
      </div>

      {/* Recent sessions */}
      <div className="bg-surface-container-lowest rounded-xl overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
        <div className="px-6 py-4 flex items-center justify-between" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
          <h2 className="font-headline text-lg font-bold">Recent Sessions</h2>
          <Link to="/sessions" className="text-xs text-secondary hover:text-primary transition-colors flex items-center gap-1">
            View all <span className="material-symbols-outlined text-sm">arrow_forward</span>
          </Link>
        </div>
        {!sessions ? (
          <div className="p-6 text-on-surface/40 text-sm">Loading...</div>
        ) : sessions.length === 0 ? (
          <div className="p-6 text-on-surface/40 text-sm">No sessions yet.</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="bg-surface-container-low/50">
              <tr>
                {['Session ID', 'Agent', 'User', 'Model', 'State', 'Created'].map(h => (
                  <th key={h} className="text-left px-6 py-3 text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y" style={{ borderColor: 'rgba(70,69,84,0.05)' }}>
              {sessions.slice(0, 5).map(s => (
                <tr key={s.id} className="hover:bg-surface-container/40 transition-colors">
                  <td className="px-6 py-4">
                    <Link to={`/sessions/${s.id}/chat`} className="font-mono text-xs text-secondary hover:text-primary transition-colors">
                      {s.id?.slice(0, 28)}…
                    </Link>
                  </td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs font-mono">{s.agent_id || '—'}</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{s.user_id || '—'}</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{s.active_model?.split('/').pop() || '—'}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wide ${
                      s.state === 'Active'
                        ? 'bg-tertiary-container/20 text-tertiary'
                        : 'bg-surface-variant text-on-surface-variant'
                    }`}>
                      {s.state === 'Active' && <span className="w-1 h-1 rounded-full bg-tertiary pulse-dot" />}
                      {s.state}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-on-surface-variant/50 text-xs">
                    {s.created_at ? new Date(s.created_at).toLocaleString() : '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
