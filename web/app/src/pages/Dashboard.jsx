import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import api, { cronApi } from '../api/gateway.js'

function StatCard({ title, value, sub, to, color = 'blue' }) {
  const colors = {
    blue:  'bg-blue-50 border-blue-200 text-blue-700',
    green: 'bg-green-50 border-green-200 text-green-700',
    purple:'bg-purple-50 border-purple-200 text-purple-700',
    amber: 'bg-amber-50 border-amber-200 text-amber-700',
  }
  const card = (
    <div className={`border rounded-xl p-5 ${colors[color]} transition-shadow hover:shadow-md`}>
      <div className="text-sm font-medium opacity-70 mb-1">{title}</div>
      <div className="text-3xl font-bold">{value ?? '—'}</div>
      {sub && <div className="text-xs mt-1 opacity-60">{sub}</div>}
    </div>
  )
  return to ? <Link to={to}>{card}</Link> : card
}

export default function Dashboard() {
  const [health, setHealth]     = useState(null)
  const [sessions, setSessions] = useState(null)
  const [channels, setChannels] = useState(null)
  const [memories, setMemories] = useState(null)
  const [cronStatus, setCronStatus] = useState(null)

  useEffect(() => {
    api.health().then(setHealth).catch(() => setHealth({ status: 'error' }))
    api.sessions.list({ limit: 100 }).then(d => setSessions(d?.sessions ?? d ?? [])).catch(() => setSessions([]))
    api.channels.list().then(d => setChannels(Array.isArray(d) ? d : [])).catch(() => setChannels([]))
    api.memories.list().then(d => setMemories(d?.memories ?? [])).catch(() => setMemories([]))
    cronApi.status().then(setCronStatus).catch(() => {})
  }, [])

  const dbOk    = health?.database === 'ok'
  const redisOk = health?.redis === 'ok' || health?.redis?.startsWith?.('ok')

  return (
    <div className="p-8 max-w-5xl">
      <h1 className="text-2xl font-bold text-gray-800 mb-2">Dashboard</h1>
      <p className="text-gray-500 mb-8 text-sm">Groundhog Gateway overview</p>

      {/* Status bar */}
      <div className="flex gap-3 mb-8 flex-wrap">
        {[
          { label: 'API',      ok: health !== null,  val: health?.status },
          { label: 'Database', ok: dbOk,             val: dbOk ? 'ok' : health?.database ?? '…' },
          { label: 'Redis',    ok: redisOk,          val: redisOk ? 'ok' : (health?.redis ?? '…') },
        ].map(({ label, ok, val }) => (
          <div key={label} className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium border ${
            ok ? 'bg-green-50 border-green-200 text-green-700' : 'bg-red-50 border-red-200 text-red-600'
          }`}>
            <span>{ok ? '●' : '○'}</span>
            <span>{label}</span>
            <span className="opacity-60">{val}</span>
          </div>
        ))}
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10">
        <StatCard title="Sessions"  value={sessions?.length}  sub="total"  to="/sessions"  color="blue"   />
        <StatCard title="Channels"  value={channels?.length}  sub="total"  to="/channels"  color="green"  />
        <StatCard title="Memories"  value={memories?.length}  sub="stored" to="/memory"    color="purple" />
        <StatCard title="Cron Jobs" value={cronStatus?.enabled_jobs ?? '…'} sub={cronStatus?.running ? 'scheduler running' : 'scheduler stopped'} to="/cron" color="amber" />
      </div>

      {/* Recent sessions */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <div className="px-5 py-3 border-b border-gray-100 flex items-center justify-between">
          <h2 className="font-semibold text-gray-700 text-sm">Recent Sessions</h2>
          <Link to="/sessions" className="text-xs text-blue-600 hover:underline">View all →</Link>
        </div>
        {!sessions ? (
          <div className="p-5 text-gray-400 text-sm">Loading...</div>
        ) : sessions.length === 0 ? (
          <div className="p-5 text-gray-400 text-sm">No sessions yet.</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-gray-500 text-xs">
              <tr>
                <th className="text-left px-5 py-2">Session ID</th>
                <th className="text-left px-5 py-2">User</th>
                <th className="text-left px-5 py-2">Model</th>
                <th className="text-left px-5 py-2">State</th>
                <th className="text-left px-5 py-2">Created</th>
              </tr>
            </thead>
            <tbody>
              {sessions.slice(0, 5).map(s => (
                <tr key={s.id} className="border-t border-gray-50 hover:bg-gray-50">
                  <td className="px-5 py-2.5 font-mono text-xs text-gray-500">
                    <Link to={`/sessions/${s.id}/chat`} className="text-blue-600 hover:underline">
                      {s.id?.slice(0, 24)}…
                    </Link>
                  </td>
                  <td className="px-5 py-2.5 text-gray-600">{s.user_id || '-'}</td>
                  <td className="px-5 py-2.5 text-gray-500 text-xs">{s.active_model?.split('/').pop() || '-'}</td>
                  <td className="px-5 py-2.5">
                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                      s.state === 'Active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'
                    }`}>{s.state}</span>
                  </td>
                  <td className="px-5 py-2.5 text-gray-400 text-xs">
                    {s.created_at ? new Date(s.created_at).toLocaleString() : '-'}
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
