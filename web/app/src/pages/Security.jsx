import React, { useState, useEffect } from 'react'
import api from '../api/gateway.js'

const PAGE_SIZE = 20

export default function Security() {
  const [logs, setLogs]     = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError]   = useState(null)
  const [page, setPage]     = useState(1)
  const [total, setTotal]   = useState(0)

  useEffect(() => {
    setLoading(true)
    api.security.audit({ page, page_size: PAGE_SIZE })
      .then(d => {
        setLogs(Array.isArray(d) ? d : (d?.items ?? []))
        setTotal(d?.total ?? 0)
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [page])

  return (
    <div className="p-8 max-w-5xl">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Security — Audit Logs</h1>

      {loading ? (
        <div className="text-gray-400 text-sm">Loading...</div>
      ) : error ? (
        <div className="text-red-500 text-sm">Error: {error}</div>
      ) : logs.length === 0 ? (
        <div className="bg-white rounded-xl border border-gray-100 p-8 text-center text-gray-400 text-sm">
          No audit logs found.
        </div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-gray-500 text-xs">
              <tr>
                <th className="text-left px-5 py-2.5">Action</th>
                <th className="text-left px-5 py-2.5">Principal</th>
                <th className="text-left px-5 py-2.5">Resource</th>
                <th className="text-left px-5 py-2.5">Resource ID</th>
                <th className="text-left px-5 py-2.5">Time</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((log, i) => (
                <tr key={log.id ?? i} className="border-t border-gray-50 hover:bg-gray-50">
                  <td className="px-5 py-3 font-medium text-gray-800">{log.action}</td>
                  <td className="px-5 py-3 text-gray-600">{log.principal_id || '-'}</td>
                  <td className="px-5 py-3 text-gray-500">{log.resource_type || '-'}</td>
                  <td className="px-5 py-3 font-mono text-xs text-gray-400">{log.resource_id || '-'}</td>
                  <td className="px-5 py-3 text-gray-400 text-xs">
                    {log.created_at ? new Date(log.created_at).toLocaleString() : '-'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="flex items-center gap-3 mt-4">
        <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1 || loading}
          className="px-3 py-1.5 rounded border text-sm disabled:opacity-40 hover:bg-gray-100">← Prev</button>
        <span className="text-sm text-gray-600">Page {page}{total > 0 ? ` · ${total} total` : ''}</span>
        <button onClick={() => setPage(p => p + 1)} disabled={logs.length < PAGE_SIZE || loading}
          className="px-3 py-1.5 rounded border text-sm disabled:opacity-40 hover:bg-gray-100">Next →</button>
      </div>
    </div>
  )
}
