import React, { useState } from 'react'

const FIELDS = [
  { key: 'user_id', label: 'User ID',   placeholder: 'user-001', help: 'Used as X-User-ID header for memory APIs', icon: 'person' },
  { key: 'token',   label: 'JWT Token', placeholder: 'Bearer …', help: 'Optional auth token for protected endpoints', icon: 'key' },
]

const ENDPOINTS = [
  'GET  /api/v1/health',
  'GET  /api/v1/agents',
  'GET  /api/v1/sessions',
  'POST /api/v1/sessions',
  'POST /api/v1/sessions/:id/messages',
  'POST /api/v1/sessions/:id/messages/stream',
  'GET  /api/v1/sessions/:id/approvals',
  'POST /api/v1/sessions/:id/approvals/:approval_id',
  'GET  /api/v1/channels',
  'POST /api/v1/channels',
  'DELETE /api/v1/channels/:id',
  'GET  /api/v1/memories',
  'POST /api/v1/memories',
  'POST /api/v1/memories/search',
  'GET  /api/v1/security/audit',
  'POST /rpc  (cron.list, cron.add, cron.update, cron.remove, cron.run, cron.runs)',
]

export default function Config() {
  const [values, setValues] = useState(() =>
    Object.fromEntries(FIELDS.map(f => [f.key, localStorage.getItem(f.key) || '']))
  )
  const [saved, setSaved] = useState(false)

  function handleSave(e) {
    e.preventDefault()
    FIELDS.forEach(f => {
      if (values[f.key]) localStorage.setItem(f.key, values[f.key])
      else localStorage.removeItem(f.key)
    })
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div className="p-10 max-w-3xl">
      <div className="mb-10">
        <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">Config</h2>
        <p className="text-on-surface-variant text-sm mt-1">Client-side settings stored in localStorage</p>
      </div>

      {/* Settings form */}
      <div className="bg-surface-container rounded-xl p-8 mb-8" style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
        <h3 className="font-headline font-bold text-base mb-6">Client Settings</h3>
        <form onSubmit={handleSave} className="space-y-5">
          {FIELDS.map(f => (
            <label key={f.key} className="flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant flex items-center gap-2">
                <span className="material-symbols-outlined text-sm">{f.icon}</span>
                {f.label}
              </span>
              <input
                className="bg-surface-container-low rounded-lg px-4 py-3 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all placeholder:text-on-surface-variant/30"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                placeholder={f.placeholder}
                value={values[f.key]}
                onChange={e => setValues(v => ({ ...v, [f.key]: e.target.value }))}
              />
              <span className="text-xs text-on-surface-variant/50 ml-1">{f.help}</span>
            </label>
          ))}
          <div className="flex items-center gap-4 pt-2">
            <button type="submit"
              className="px-8 py-2.5 text-sm font-bold rounded-lg text-on-primary transition-all hover:brightness-110"
              style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
              Save Settings
            </button>
            {saved && (
              <span className="text-tertiary text-sm flex items-center gap-1">
                <span className="material-symbols-outlined text-sm">check_circle</span> Saved
              </span>
            )}
          </div>
        </form>
      </div>

      {/* API reference */}
      <div className="bg-surface-container-lowest rounded-xl overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
        <div className="px-6 py-4 flex items-center gap-3" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
          <span className="material-symbols-outlined text-secondary">api</span>
          <h3 className="font-headline font-bold text-base">API Reference</h3>
        </div>
        <div className="p-6 space-y-2">
          {ENDPOINTS.map(ep => {
            const [method, ...rest] = ep.split(' ')
            const path = rest.join(' ')
            const methodColor = {
              GET: 'text-tertiary', POST: 'text-secondary', DELETE: 'text-error', PUT: 'text-primary',
            }[method] || 'text-on-surface-variant'
            return (
              <div key={ep} className="flex items-center gap-3 py-1.5 px-3 rounded hover:bg-surface-container transition-colors">
                <span className={`font-mono text-[10px] font-bold uppercase w-14 shrink-0 ${methodColor}`}>{method}</span>
                <code className="font-mono text-xs text-on-surface-variant">{path}</code>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
