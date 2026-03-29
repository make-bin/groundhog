import React, { useState, useEffect } from 'react'
import api from '../api/gateway.js'

const PAGE_SIZE = 20

export default function Security() {
  const [logs, setLogs]       = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState(null)
  const [page, setPage]       = useState(1)
  const [total, setTotal]     = useState(0)
  const [filterAction, setFilterAction]   = useState('')
  const [filterPrincipal, setFilterPrincipal] = useState('')

  useEffect(() => {
    setLoading(true)
    const params = { page, page_size: PAGE_SIZE }
    if (filterAction) params.action = filterAction
    if (filterPrincipal) params.principal_id = filterPrincipal
    api.security.audit(params)
      .then(d => {
        setLogs(Array.isArray(d) ? d : (d?.items ?? []))
        setTotal(d?.total ?? 0)
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [page, filterAction, filterPrincipal])

  return (
    <div className="p-10 max-w-6xl">
      <div className="mb-10">
        <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">Security</h2>
        <p className="text-on-surface-variant text-sm mt-1">Audit log of all system operations</p>
      </div>

      {/* Filters */}
      <div className="flex gap-4 mb-8">
        <div className="relative flex-1 max-w-xs">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-sm">filter_list</span>
          <input
            className="w-full bg-surface-container-low rounded-lg py-2.5 pl-10 pr-4 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all placeholder:text-on-surface-variant/30"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}
            placeholder="Filter by action…"
            value={filterAction}
            onChange={e => { setFilterAction(e.target.value); setPage(1) }}
          />
        </div>
        <div className="relative flex-1 max-w-xs">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-sm">person</span>
          <input
            className="w-full bg-surface-container-low rounded-lg py-2.5 pl-10 pr-4 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all placeholder:text-on-surface-variant/30"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}
            placeholder="Filter by principal ID…"
            value={filterPrincipal}
            onChange={e => { setFilterPrincipal(e.target.value); setPage(1) }}
          />
        </div>
        {total > 0 && (
          <div className="flex items-center text-xs text-on-surface-variant ml-auto">
            {total} total records
          </div>
        )}
      </div>

      {loading ? (
        <div className="text-on-surface/40 text-sm flex items-center gap-2">
          <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
        </div>
      ) : error ? (
        <div className="text-error text-sm">{error}</div>
      ) : logs.length === 0 ? (
        <div className="bg-surface-container-low rounded-xl p-16 text-center">
          <span className="material-symbols-outlined text-on-surface-variant/30 mb-4" style={{ fontSize: '3rem' }}>shield</span>
          <p className="text-on-surface-variant text-sm">No audit logs found.</p>
        </div>
      ) : (
        <div className="bg-surface-container-lowest rounded-xl overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
          <table className="w-full text-sm">
            <thead className="bg-surface-container-low/50">
              <tr>
                {['Action', 'Principal', 'Resource', 'Resource ID', 'Time'].map(h => (
                  <th key={h} className="text-left px-6 py-3 text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y" style={{ borderColor: 'rgba(70,69,84,0.05)' }}>
              {logs.map((log, i) => (
                <tr key={log.id ?? i} className="hover:bg-surface-container/40 transition-colors">
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded bg-primary/10 text-primary text-[10px] font-bold uppercase tracking-wide">
                      {log.action}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{log.principal_id || '—'}</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{log.resource_type || '—'}</td>
                  <td className="px-6 py-4 font-mono text-xs text-on-surface-variant/50">{log.resource_id || '—'}</td>
                  <td className="px-6 py-4 text-on-surface-variant/50 text-xs">
                    {log.created_at ? new Date(log.created_at).toLocaleString() : '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Pagination */}
      <div className="flex items-center gap-3 mt-6">
        <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1 || loading}
          className="px-4 py-2 rounded-lg text-sm font-medium text-on-surface-variant hover:bg-surface-container transition-colors disabled:opacity-40 flex items-center gap-1"
          style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
          <span className="material-symbols-outlined text-sm">arrow_back</span> Prev
        </button>
        <span className="text-sm text-on-surface-variant">Page {page}{total > 0 ? ` · ${total} total` : ''}</span>
        <button onClick={() => setPage(p => p + 1)} disabled={logs.length < PAGE_SIZE || loading}
          className="px-4 py-2 rounded-lg text-sm font-medium text-on-surface-variant hover:bg-surface-container transition-colors disabled:opacity-40 flex items-center gap-1"
          style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
          Next <span className="material-symbols-outlined text-sm">arrow_forward</span>
        </button>
      </div>
    </div>
  )
}
