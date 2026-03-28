import React, { useState, useEffect } from 'react'
import api from '../api/gateway.js'

const CHANNEL_TYPES = ['discord', 'telegram', 'whatsapp', 'slack']

export default function Channels() {
  const [channels, setChannels] = useState([])
  const [loading, setLoading]   = useState(true)
  const [error, setError]       = useState(null)
  const [form, setForm]         = useState({ channel_type: 'discord', plugin_id: '' })
  const [creating, setCreating] = useState(false)
  const [formErr, setFormErr]   = useState(null)

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
    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
      s === 'Active' ? 'bg-green-100 text-green-700' :
      s === 'Error'  ? 'bg-red-100 text-red-600' :
                       'bg-gray-100 text-gray-500'
    }`}>{s || 'Inactive'}</span>
  )

  return (
    <div className="p-8 max-w-4xl">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Channels</h1>

      {/* Create */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6 shadow-sm">
        <h2 className="font-semibold text-gray-700 text-sm mb-3">Add Channel</h2>
        <form onSubmit={handleCreate} className="flex gap-3 flex-wrap items-end">
          <label className="flex flex-col gap-1 text-xs text-gray-600">
            Type
            <select className="border rounded px-3 py-1.5 text-sm" value={form.channel_type}
              onChange={e => setForm(f => ({ ...f, channel_type: e.target.value }))}>
              {CHANNEL_TYPES.map(t => <option key={t}>{t}</option>)}
            </select>
          </label>
          <label className="flex flex-col gap-1 text-xs text-gray-600 flex-1 min-w-40">
            Plugin ID (optional)
            <input className="border rounded px-3 py-1.5 text-sm" placeholder="plugin-001"
              value={form.plugin_id}
              onChange={e => setForm(f => ({ ...f, plugin_id: e.target.value }))} />
          </label>
          <button type="submit" disabled={creating}
            className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm hover:bg-blue-700 disabled:opacity-50">
            {creating ? 'Adding…' : 'Add Channel'}
          </button>
        </form>
        {formErr && <p className="text-red-500 text-xs mt-2">{formErr}</p>}
      </div>

      {/* List */}
      {loading ? (
        <div className="text-gray-400 text-sm">Loading...</div>
      ) : error ? (
        <div className="text-red-500 text-sm">Error: {error}</div>
      ) : channels.length === 0 ? (
        <div className="text-gray-400 text-sm">No channels configured.</div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-gray-500 text-xs">
              <tr>
                <th className="text-left px-5 py-2.5">ID</th>
                <th className="text-left px-5 py-2.5">Type</th>
                <th className="text-left px-5 py-2.5">Plugin</th>
                <th className="text-left px-5 py-2.5">Status</th>
                <th className="px-5 py-2.5" />
              </tr>
            </thead>
            <tbody>
              {channels.map(ch => (
                <tr key={ch.id} className="border-t border-gray-50 hover:bg-gray-50">
                  <td className="px-5 py-3 font-mono text-xs text-gray-500">{ch.id}</td>
                  <td className="px-5 py-3 capitalize">{ch.channel_type}</td>
                  <td className="px-5 py-3 text-gray-500 text-xs">{ch.plugin_id || '-'}</td>
                  <td className="px-5 py-3">{statusBadge(ch.status)}</td>
                  <td className="px-5 py-3 text-right">
                    <button onClick={() => handleDelete(ch.id)}
                      className="text-red-400 hover:text-red-600 text-xs">Delete</button>
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
