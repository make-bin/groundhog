import React, { useState, useEffect } from 'react'
import api from '../api/gateway.js'

const CHANNEL_TYPES = ['discord', 'telegram', 'whatsapp', 'slack', 'signal']

const TYPE_ICONS = {
  discord: 'chat_bubble', telegram: 'send', whatsapp: 'phone_iphone',
  slack: 'workspaces', signal: 'lock',
}

export default function Channels() {
  const [channels, setChannels] = useState([])
  const [loading, setLoading]   = useState(true)
  const [error, setError]       = useState(null)
  const [form, setForm]         = useState({ channel_type: 'discord', plugin_id: '' })
  const [creating, setCreating] = useState(false)
  const [formErr, setFormErr]   = useState(null)
  const [showForm, setShowForm] = useState(false)

  function load() {
    setLoading(true)
    api.channels.list()
      .then(d => setChannels(Array.isArray(d) ? d : []))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  async function handleCreate(e) {
    e.preventDefault()
    setFormErr(null)
    setCreating(true)
    try {
      await api.channels.create(form)
      setForm({ channel_type: 'discord', plugin_id: '' })
      setShowForm(false)
      load()
    } catch (e) {
      setFormErr(e.message)
    } finally {
      setCreating(false)
    }
  }

  async function handleDelete(id) {
    if (!confirm('Delete this channel?')) return
    await api.channels.delete(id).catch(() => {})
    load()
  }

  const statusBadge = (s) => (
    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wide ${
      s === 'Active' ? 'bg-tertiary-container/20 text-tertiary' :
      s === 'Error'  ? 'bg-error-container/20 text-error' :
                       'bg-surface-variant text-on-surface-variant'
    }`}>
      <span className={`w-1 h-1 rounded-full ${s === 'Active' ? 'bg-tertiary pulse-dot' : s === 'Error' ? 'bg-error' : 'bg-on-surface-variant'}`} />
      {s || 'Inactive'}
    </span>
  )

  return (
    <div className="p-10 max-w-5xl">
      <div className="flex items-end justify-between mb-10">
        <div>
          <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">Channels</h2>
          <p className="text-on-surface-variant text-sm mt-1">Manage external messaging channel integrations</p>
        </div>
        <button onClick={() => setShowForm(v => !v)}
          className="px-6 py-2.5 text-sm font-bold rounded text-on-primary flex items-center gap-2 shadow-lg transition-all hover:brightness-110"
          style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
          <span className="material-symbols-outlined text-sm">add</span> Add Channel
        </button>
      </div>

      {/* Add form */}
      {showForm && (
        <div className="bg-surface-container rounded-xl p-6 mb-8" style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
          <h3 className="font-headline font-bold text-base mb-5">New Channel Configuration</h3>
          <form onSubmit={handleCreate} className="grid grid-cols-2 gap-5">
            <label className="flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Channel Type</span>
              <select
                className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                value={form.channel_type}
                onChange={e => setForm(f => ({ ...f, channel_type: e.target.value }))}>
                {CHANNEL_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </label>
            <label className="flex flex-col gap-1.5">
              <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Plugin ID (optional)</span>
              <input
                className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all placeholder:text-on-surface-variant/30"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                placeholder="plugin-001"
                value={form.plugin_id}
                onChange={e => setForm(f => ({ ...f, plugin_id: e.target.value }))}
              />
            </label>
            {formErr && <p className="col-span-2 text-error text-xs">{formErr}</p>}
            <div className="col-span-2 flex justify-end gap-3">
              <button type="button" onClick={() => setShowForm(false)}
                className="px-5 py-2 text-sm font-bold text-on-surface-variant hover:text-on-surface transition-colors">
                Discard
              </button>
              <button type="submit" disabled={creating}
                className="px-8 py-2 text-sm font-bold rounded text-on-primary disabled:opacity-50 transition-all hover:brightness-110"
                style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
                {creating ? 'Adding…' : 'Deploy Channel'}
              </button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <div className="text-on-surface/40 text-sm flex items-center gap-2">
          <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
        </div>
      ) : error ? (
        <div className="text-error text-sm">{error}</div>
      ) : channels.length === 0 ? (
        <div className="bg-surface-container-low rounded-xl p-16 text-center">
          <span className="material-symbols-outlined text-on-surface-variant/30 mb-4" style={{ fontSize: '3rem' }}>hub</span>
          <p className="text-on-surface-variant text-sm">No channels configured.</p>
        </div>
      ) : (
        <div className="bg-surface-container-lowest rounded-xl overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
          <table className="w-full text-sm">
            <thead className="bg-surface-container-low/50">
              <tr>
                {['Channel', 'Type', 'Plugin', 'Status', ''].map(h => (
                  <th key={h} className="text-left px-6 py-3 text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y" style={{ borderColor: 'rgba(70,69,84,0.05)' }}>
              {channels.map(ch => (
                <tr key={ch.id} className="hover:bg-surface-container/40 transition-colors">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
                        <span className="material-symbols-outlined text-primary text-sm">{TYPE_ICONS[ch.channel_type] || 'hub'}</span>
                      </div>
                      <span className="font-mono text-xs text-on-surface-variant">{ch.id}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 capitalize text-on-surface text-sm font-medium">{ch.channel_type}</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs font-mono">{ch.plugin_id || '—'}</td>
                  <td className="px-6 py-4">{statusBadge(ch.status)}</td>
                  <td className="px-6 py-4 text-right">
                    <button onClick={() => handleDelete(ch.id)}
                      className="text-on-surface-variant/40 hover:text-error transition-colors">
                      <span className="material-symbols-outlined text-sm">delete</span>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
