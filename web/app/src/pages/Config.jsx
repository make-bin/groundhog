import React, { useState } from 'react'

const FIELDS = [
  { key: 'user_id',   label: 'User ID',   placeholder: 'user-001', help: 'Used as X-User-ID header for memory APIs' },
  { key: 'token',     label: 'JWT Token', placeholder: 'Bearer …', help: 'Optional auth token' },
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
    <div className="p-8 max-w-2xl">
      <h1 className="text-2xl font-bold text-gray-800 mb-2">Config</h1>
      <p className="text-gray-500 text-sm mb-6">Client-side settings stored in localStorage.</p>

      <div className="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
        <form onSubmit={handleSave} className="space-y-4">
          {FIELDS.map(f => (
            <label key={f.key} className="flex flex-col gap-1">
              <span className="text-sm font-medium text-gray-700">{f.label}</span>
              <input
                className="border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
                placeholder={f.placeholder}
                value={values[f.key]}
                onChange={e => setValues(v => ({ ...v, [f.key]: e.target.value }))}
              />
              <span className="text-xs text-gray-400">{f.help}</span>
            </label>
          ))}

          <div className="flex items-center gap-3 pt-2">
            <button type="submit"
              className="bg-blue-600 text-white px-5 py-2 rounded-lg text-sm hover:bg-blue-700 transition-colors">
              Save
            </button>
            {saved && <span className="text-green-600 text-sm">✓ Saved</span>}
          </div>
        </form>
      </div>

      {/* API info */}
      <div className="mt-6 bg-white rounded-xl border border-gray-200 p-5 shadow-sm">
        <h2 className="font-semibold text-gray-700 text-sm mb-3">API Endpoints</h2>
        <div className="space-y-1 text-xs font-mono text-gray-500">
          {[
            'GET  /api/v1/health',
            'GET  /api/v1/sessions',
            'POST /api/v1/sessions',
            'POST /api/v1/sessions/:id/messages',
            'POST /api/v1/sessions/:id/messages/stream',
            'GET  /api/v1/channels',
            'POST /api/v1/channels',
            'GET  /api/v1/memories',
            'POST /api/v1/memories',
            'POST /api/v1/memories/search',
            'GET  /api/v1/security/audit',
          ].map(e => <div key={e}>{e}</div>)}
        </div>
      </div>
    </div>
  )
}
