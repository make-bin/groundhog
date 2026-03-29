import React, { useState, useEffect } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import api from '../api/gateway.js'

const PROVIDERS = ['openai_compat', 'openai', 'anthropic', 'ollama', 'groq', 'mistral']

function CreateModal({ onClose, onCreated, defaultAgentID }) {
  const navigate = useNavigate()
  const [agents, setAgents]   = useState([])
  const [form, setForm] = useState({
    user_id:      localStorage.getItem('user_id') || 'user-001',
    agent_id:     defaultAgentID || '',
    provider:     '',
    model_name:   '',
    system_prompt: '',
  })
  const [creating, setCreating] = useState(false)
  const [err, setErr]           = useState(null)
  const [showAdvanced, setShowAdvanced] = useState(false)

  // Load agents for the picker
  useEffect(() => {
    api.agents.list()
      .then(d => {
        const list = Array.isArray(d) ? d : []
        setAgents(list)
        // Auto-select default agent if none pre-selected
        if (!form.agent_id && list.length > 0) {
          const def = list.find(a => a.is_default) ?? list[0]
          setForm(f => ({ ...f, agent_id: def.id }))
        }
      })
      .catch(() => {})
  }, [])

  async function handleCreate(e) {
    e.preventDefault()
    setErr(null)
    setCreating(true)
    try {
      // Only send provider/model if explicitly overridden (registry fills defaults server-side)
      const payload = {
        agent_id:      form.agent_id,
        user_id:       form.user_id,
        system_prompt: form.system_prompt || undefined,
      }
      if (showAdvanced && form.provider) payload.provider   = form.provider
      if (showAdvanced && form.model_name) payload.model_name = form.model_name
      const sess = await api.sessions.create(payload)
      onCreated()
      navigate(`/sessions/${sess.id}/chat`)
    } catch (e) {
      setErr(e.message)
    } finally {
      setCreating(false)
    }
  }

  const selectedAgent = agents.find(a => a.id === form.agent_id)

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ background: 'rgba(11,19,38,0.85)', backdropFilter: 'blur(8px)' }}>
      <div className="bg-surface-container w-full max-w-lg rounded-2xl overflow-hidden shadow-2xl" style={{ border: '1px solid rgba(70,69,84,0.3)' }}>
        <div className="bg-surface-container-high px-8 py-6 flex justify-between items-center" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
          <h3 className="font-headline text-xl font-bold">New Session</h3>
          <button onClick={onClose} className="text-on-surface-variant hover:text-on-surface transition-colors">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>
        <form onSubmit={handleCreate} className="p-8 flex flex-col gap-5">

          {/* Agent picker */}
          <div className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Agent</span>
            {agents.length > 0 ? (
              <div className="grid grid-cols-1 gap-2 max-h-48 overflow-y-auto pr-1">
                {agents.map(a => (
                  <button key={a.id} type="button"
                    onClick={() => setForm(f => ({ ...f, agent_id: a.id }))}
                    className={`flex items-center gap-3 px-4 py-3 rounded-lg text-left transition-all ${
                      form.agent_id === a.id
                        ? 'bg-primary/10 border-primary/40'
                        : 'bg-surface-container-low hover:bg-surface-container-high'
                    }`}
                    style={{ border: `1px solid ${form.agent_id === a.id ? 'rgba(192,193,255,0.4)' : 'rgba(70,69,84,0.2)'}` }}>
                    <div className={`w-8 h-8 rounded-lg flex items-center justify-center shrink-0 ${form.agent_id === a.id ? 'bg-primary/20' : 'bg-surface-container'}`}>
                      <span className={`material-symbols-outlined text-sm ${form.agent_id === a.id ? 'text-primary' : 'text-on-surface-variant'}`}
                        style={{ fontVariationSettings: "'FILL' 1" }}>smart_toy</span>
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-on-surface">{a.name}</span>
                        {a.is_default && <span className="text-[10px] text-tertiary font-bold uppercase">default</span>}
                      </div>
                      {a.model && (
                        <span className="text-[10px] font-mono text-on-surface-variant/50">{a.model.split('/').pop()}</span>
                      )}
                    </div>
                    {form.agent_id === a.id && (
                      <span className="material-symbols-outlined text-primary text-sm shrink-0">check_circle</span>
                    )}
                  </button>
                ))}
              </div>
            ) : (
              <input
                required
                className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all placeholder:text-on-surface-variant/30"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                placeholder="agent-id"
                value={form.agent_id}
                onChange={e => setForm(f => ({ ...f, agent_id: e.target.value }))}
              />
            )}
          </div>

          {/* User ID */}
          <label className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">User ID</span>
            <input
              required
              className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all"
              style={{ border: '1px solid rgba(70,69,84,0.2)' }}
              value={form.user_id}
              onChange={e => setForm(f => ({ ...f, user_id: e.target.value }))}
            />
          </label>

          {/* System prompt override */}
          <label className="flex flex-col gap-1.5">
            <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">
              System Prompt Override <span className="text-on-surface-variant/40 normal-case font-normal">(optional)</span>
            </span>
            <textarea rows={2}
              className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none resize-none transition-all placeholder:text-on-surface-variant/30"
              style={{ border: '1px solid rgba(70,69,84,0.2)' }}
              placeholder={selectedAgent?.model ? `Using ${selectedAgent.name} defaults…` : 'Leave blank to use agent defaults'}
              value={form.system_prompt}
              onChange={e => setForm(f => ({ ...f, system_prompt: e.target.value }))}
            />
          </label>

          {/* Advanced toggle */}
          <button type="button" onClick={() => setShowAdvanced(v => !v)}
            className="flex items-center gap-1.5 text-xs text-on-surface-variant/60 hover:text-on-surface-variant transition-colors self-start">
            <span className="material-symbols-outlined text-sm">{showAdvanced ? 'expand_less' : 'expand_more'}</span>
            {showAdvanced ? 'Hide' : 'Show'} model overrides
          </button>

          {showAdvanced && (
            <div className="grid grid-cols-2 gap-4 bg-surface-container-lowest rounded-xl p-4"
              style={{ border: '1px solid rgba(70,69,84,0.15)' }}>
              <label className="flex flex-col gap-1.5">
                <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Provider</span>
                <select
                  className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all"
                  style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                  value={form.provider}
                  onChange={e => setForm(f => ({ ...f, provider: e.target.value }))}>
                  <option value="">— agent default —</option>
                  {PROVIDERS.map(p => <option key={p} value={p}>{p}</option>)}
                </select>
              </label>
              <label className="flex flex-col gap-1.5">
                <span className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">Model</span>
                <input
                  className="bg-surface-container-low rounded px-3 py-2 text-sm text-on-surface focus:ring-1 focus:ring-primary focus:outline-none transition-all placeholder:text-on-surface-variant/30"
                  style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                  placeholder="agent default"
                  value={form.model_name}
                  onChange={e => setForm(f => ({ ...f, model_name: e.target.value }))}
                />
              </label>
            </div>
          )}

          {err && <p className="text-error text-xs">{err}</p>}

          <div className="flex justify-end gap-3 pt-1">
            <button type="button" onClick={onClose}
              className="px-5 py-2 text-sm font-bold text-on-surface-variant hover:text-on-surface transition-colors">
              Discard
            </button>
            <button type="submit" disabled={creating || !form.agent_id}
              className="px-8 py-2 text-sm font-bold rounded text-on-primary disabled:opacity-50 transition-all hover:brightness-110"
              style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
              {creating ? 'Creating…' : 'Create & Open Chat'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default function Sessions() {
  const navigate = useNavigate()
  const location = useLocation()
  const [sessions, setSessions] = useState([])
  const [loading, setLoading]   = useState(true)
  const [error, setError]       = useState(null)
  // Pre-select agent if navigated from Agents page
  const [showCreate, setShowCreate] = useState(!!location.state?.agentID)
  const [defaultAgentID] = useState(location.state?.agentID ?? null)

  function load() {
    setLoading(true)
    api.sessions.list()
      .then(d => setSessions(d?.sessions ?? (Array.isArray(d) ? d : [])))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  async function handleDelete(id, ev) {
    ev.stopPropagation()
    if (!confirm('Delete this session?')) return
    await api.sessions.delete(id).catch(() => {})
    load()
  }

  return (
    <div className="p-10 max-w-6xl">
      <div className="flex items-end justify-between mb-10">
        <div>
          <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">Sessions</h2>
          <p className="text-on-surface-variant text-sm mt-1">Manage AI agent conversation sessions</p>
        </div>
        <button onClick={() => setShowCreate(true)}
          className="px-6 py-2.5 text-sm font-bold rounded text-on-primary flex items-center gap-2 shadow-lg transition-all hover:brightness-110"
          style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
          <span className="material-symbols-outlined text-sm">add</span> New Session
        </button>
      </div>

      {loading ? (
        <div className="text-on-surface/40 text-sm flex items-center gap-2">
          <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
        </div>
      ) : error ? (
        <div className="text-error text-sm">{error}</div>
      ) : sessions.length === 0 ? (
        <div className="bg-surface-container-low rounded-xl p-16 text-center">
          <span className="material-symbols-outlined text-on-surface-variant/30 mb-4" style={{ fontSize: '3rem' }}>forum</span>
          <p className="text-on-surface-variant text-sm">No sessions yet. Create one to start chatting.</p>
        </div>
      ) : (
        <div className="bg-surface-container-lowest rounded-xl overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
          <table className="w-full text-sm">
            <thead className="bg-surface-container-low/50">
              <tr>
                {['Session ID', 'User', 'Model', 'Turns', 'State', 'Created', ''].map(h => (
                  <th key={h} className="text-left px-6 py-3 text-[10px] uppercase tracking-widest font-bold text-on-surface-variant">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y" style={{ borderColor: 'rgba(70,69,84,0.05)' }}>
              {sessions.map(s => (
                <tr key={s.id}
                  className="hover:bg-surface-container/40 transition-colors cursor-pointer"
                  onClick={() => navigate(`/sessions/${s.id}/chat`)}>
                  <td className="px-6 py-4 font-mono text-xs text-secondary">{s.id?.slice(0, 26)}…</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{s.user_id}</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{s.active_model?.split('/').pop()}</td>
                  <td className="px-6 py-4 text-on-surface-variant text-xs">{s.turns?.length ?? 0}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wide ${
                      s.state === 'Active' ? 'bg-tertiary-container/20 text-tertiary' : 'bg-surface-variant text-on-surface-variant'
                    }`}>
                      {s.state === 'Active' && <span className="w-1 h-1 rounded-full bg-tertiary pulse-dot" />}
                      {s.state}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-on-surface-variant/50 text-xs">
                    {s.created_at ? new Date(s.created_at).toLocaleString() : '—'}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button onClick={e => handleDelete(s.id, e)}
                      className="text-on-surface-variant/40 hover:text-error transition-colors text-xs">
                      <span className="material-symbols-outlined text-sm">delete</span>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showCreate && (
        <CreateModal onClose={() => setShowCreate(false)} onCreated={load} defaultAgentID={defaultAgentID} />
      )}
    </div>
  )
}
