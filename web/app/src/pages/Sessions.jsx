import React, { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../api/gateway.js'

const PROVIDERS = ['openai_compat', 'openai', 'ollama', 'groq', 'mistral']
const DEFAULT_MODEL = 'deepseek-ai/DeepSeek-V3.1-Terminus'

export default function Sessions() {
  const navigate = useNavigate()
  const [sessions, setSessions]   = useState([])
  const [loading, setLoading]     = useState(true)
  const [error, setError]         = useState(null)
  const [showCreate, setShowCreate] = useState(false)
  const [form, setForm] = useState({
    user_id: localStorage.getItem('user_id') || 'user-001',
    agent_id: 'agent-001',
    provider: 'openai_compat',
    model_name: DEFAULT_MODEL,
    system_prompt: '',
  })
  const [creating, setCreating] = useState(false)
  const [formErr, setFormErr]   = useState(null)

  function load() {
    setLoading(true)
    api.sessions.list()
      .then(d => setSessions(d?.sessions ?? (Array.isArray(d) ? d : [])))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  async function handleCreate(e) {
    e.preventDefault()
    setFormErr(null)
    setCreating(true)
    try {
      const sess = await api.sessions.create(form)
      setShowCreate(false)
      load()
      navigate(`/sessions/${sess.id}/chat`)
    } catch (e) {
      setFormErr(e.message)
    } finally {
      setCreating(false)
    }
  }

  async function handleDelete(id, ev) {
    ev.stopPropagation()
    if (!confirm('Delete this session?')) return
    await api.sessions.delete(id).catch(() => {})
    load()
  }

  return (
    <div className="p-8 max-w-5xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Sessions</h1>
        <button
          onClick={() => setShowCreate(v => !v)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-blue-700 transition-colors"
        >
          + New Session
        </button>
      </div>

      {/* Create form */}
      {showCreate && (
        <div className="bg-white border border-gray-200 rounded-xl p-5 mb-6 shadow-sm">
          <h2 className="font-semibold text-gray-700 mb-4 text-sm">Create Session</h2>
          <form onSubmit={handleCreate} className="grid grid-cols-2 gap-3">
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              User ID
              <input className="border rounded px-3 py-1.5 text-sm" value={form.user_id}
                onChange={e => setForm(f => ({ ...f, user_id: e.target.value }))} />
            </label>
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              Agent ID
              <input className="border rounded px-3 py-1.5 text-sm" value={form.agent_id}
                onChange={e => setForm(f => ({ ...f, agent_id: e.target.value }))} />
            </label>
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              Provider
              <select className="border rounded px-3 py-1.5 text-sm" value={form.provider}
                onChange={e => setForm(f => ({ ...f, provider: e.target.value }))}>
                {PROVIDERS.map(p => <option key={p}>{p}</option>)}
              </select>
            </label>
            <label className="flex flex-col gap-1 text-xs text-gray-600">
              Model Name
              <input className="border rounded px-3 py-1.5 text-sm" value={form.model_name}
                onChange={e => setForm(f => ({ ...f, model_name: e.target.value }))} />
            </label>
            <label className="col-span-2 flex flex-col gap-1 text-xs text-gray-600">
              System Prompt (optional)
              <textarea rows={2} className="border rounded px-3 py-1.5 text-sm resize-none"
                value={form.system_prompt}
                onChange={e => setForm(f => ({ ...f, system_prompt: e.target.value }))} />
            </label>
            {formErr && <p className="col-span-2 text-red-500 text-xs">{formErr}</p>}
            <div className="col-span-2 flex gap-2">
              <button type="submit" disabled={creating}
                className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm hover:bg-blue-700 disabled:opacity-50">
                {creating ? 'Creating…' : 'Create & Open Chat'}
              </button>
              <button type="button" onClick={() => setShowCreate(false)}
                className="px-4 py-1.5 rounded text-sm border hover:bg-gray-50">
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* List */}
      {loading ? (
        <div className="text-gray-400 text-sm">Loading...</div>
      ) : error ? (
        <div className="text-red-500 text-sm">Error: {error}</div>
      ) : sessions.length === 0 ? (
        <div className="text-gray-400 text-sm">No sessions yet. Create one to start chatting.</div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-gray-500 text-xs">
              <tr>
                <th className="text-left px-5 py-2.5">Session ID</th>
                <th className="text-left px-5 py-2.5">User</th>
                <th className="text-left px-5 py-2.5">Model</th>
                <th className="text-left px-5 py-2.5">Turns</th>
                <th className="text-left px-5 py-2.5">State</th>
                <th className="text-left px-5 py-2.5">Created</th>
                <th className="px-5 py-2.5" />
              </tr>
            </thead>
            <tbody>
              {sessions.map(s => (
                <tr key={s.id} className="border-t border-gray-50 hover:bg-blue-50 cursor-pointer"
                  onClick={() => navigate(`/sessions/${s.id}/chat`)}>
                  <td className="px-5 py-3 font-mono text-xs text-blue-600">{s.id?.slice(0, 28)}…</td>
                  <td className="px-5 py-3 text-gray-600">{s.user_id}</td>
                  <td className="px-5 py-3 text-gray-500 text-xs">{s.active_model?.split('/').pop()}</td>
                  <td className="px-5 py-3 text-gray-500">{s.turns?.length ?? 0}</td>
                  <td className="px-5 py-3">
                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                      s.state === 'Active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'
                    }`}>{s.state}</span>
                  </td>
                  <td className="px-5 py-3 text-gray-400 text-xs">
                    {s.created_at ? new Date(s.created_at).toLocaleString() : '-'}
                  </td>
                  <td className="px-5 py-3 text-right">
                    <button onClick={e => handleDelete(s.id, e)}
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
