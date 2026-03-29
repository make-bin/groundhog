import React, { useState, useEffect } from 'react'
import api from '../api/gateway.js'

export default function Memory() {
  const [memories, setMemories]   = useState([])
  const [loading, setLoading]     = useState(true)
  const [error, setError]         = useState(null)
  const [content, setContent]     = useState('')
  const [creating, setCreating]   = useState(false)
  const [createErr, setCreateErr] = useState(null)
  const [query, setQuery]         = useState('')
  const [results, setResults]     = useState(null)
  const [searching, setSearching] = useState(false)
  const [editId, setEditId]       = useState(null)
  const [editContent, setEditContent] = useState('')

  function load() {
    setLoading(true)
    api.memories.list()
      .then(d => setMemories(d?.memories ?? (Array.isArray(d) ? d : [])))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  async function handleCreate(e) {
    e.preventDefault()
    if (!content.trim()) return
    setCreateErr(null)
    setCreating(true)
    try { await api.memories.create(content.trim()); setContent(''); load() }
    catch (e) { setCreateErr(e.message) }
    finally { setCreating(false) }
  }

  async function handleDelete(id) {
    if (!confirm('Delete this memory?')) return
    await api.memories.delete(id).catch(() => {})
    load()
    if (results) setResults(r => r.filter(x => (x.memory?.id ?? x.id) !== id))
  }

  async function handleSearch(e) {
    e.preventDefault()
    if (!query.trim()) return
    setSearching(true)
    try { const res = await api.memories.search(query.trim(), 10); setResults(Array.isArray(res) ? res : []) }
    catch { setResults([]) }
    finally { setSearching(false) }
  }

  async function handleUpdate(id) {
    await api.memories.update(id, editContent).catch(() => {})
    setEditId(null)
    load()
  }

  return (
    <div className="p-10 max-w-6xl">
      {/* Header */}
      <div className="flex items-end justify-between mb-10">
        <div>
          <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">
            Memory <span className="text-secondary">Vault</span>
          </h2>
          <p className="text-on-surface-variant text-sm mt-1 flex items-center gap-2">
            <span className="w-1.5 h-1.5 rounded-full bg-tertiary pulse-dot" />
            Vector indexing active
          </p>
        </div>
      </div>

      {/* Search */}
      <section className="mb-10">
        <div className="bg-surface-container-low rounded-xl p-8" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
          <label className="font-headline text-sm font-medium text-on-surface-variant block mb-2 ml-1">Vector Query</label>
          <form onSubmit={handleSearch} className="flex gap-4 items-end">
            <div className="flex-1 relative group">
              <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant group-focus-within:text-secondary transition-colors">search</span>
              <input
                className="w-full bg-surface-container-lowest rounded-lg py-4 pl-12 pr-4 text-on-surface focus:ring-2 focus:ring-secondary/50 focus:outline-none transition-all placeholder:text-on-surface-variant/30 text-sm"
                placeholder="Describe the concept or information you're looking for..."
                value={query}
                onChange={e => setQuery(e.target.value)}
              />
            </div>
            <button type="submit" disabled={searching || !query.trim()}
              className="px-8 py-4 text-sm font-bold rounded-lg text-on-primary disabled:opacity-50 transition-all hover:brightness-110 flex items-center gap-2 shrink-0"
              style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
              <span className="material-symbols-outlined text-sm">explore</span>
              {searching ? 'Searching…' : 'Vector Search'}
            </button>
          </form>
          {results !== null && (
            <div className="mt-6">
              {results.length === 0 ? (
                <p className="text-on-surface-variant/50 text-sm">No results found.</p>
              ) : (
                <div className="space-y-3">
                  <p className="text-[10px] uppercase tracking-widest font-bold text-on-surface-variant mb-3">
                    {results.length} result{results.length !== 1 ? 's' : ''} found
                  </p>
                  {results.map((r, i) => {
                    const m = r.memory ?? r
                    return (
                      <div key={m.id ?? i} className="bg-surface-container rounded-xl p-4" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
                        <p className="text-sm text-on-surface leading-relaxed">{m.content}</p>
                        <div className="mt-3 flex items-center gap-3 text-[10px] text-on-surface-variant/50">
                          <span>Match score:</span>
                          <span className="text-secondary font-mono font-bold">{(r.score ?? r.hybrid_score ?? 0).toFixed(3)}</span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </div>
          )}
        </div>
      </section>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-10">
        {/* Ingest panel */}
        <div className="lg:col-span-1">
          <div className="sticky top-8">
            <h2 className="font-headline text-xl font-bold mb-4">Ingest Knowledge</h2>
            <div className="bg-surface-container p-6 rounded-xl" style={{ border: '1px solid rgba(70,69,84,0.1)' }}>
              <form onSubmit={handleCreate} className="flex flex-col gap-4">
                <textarea rows={6}
                  className="w-full bg-surface-container-lowest rounded-lg p-4 text-sm text-on-surface focus:ring-1 focus:ring-primary/50 focus:outline-none resize-none transition-all placeholder:text-on-surface-variant/20"
                  style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                  placeholder="Paste new knowledge snippet here..."
                  value={content}
                  onChange={e => setContent(e.target.value)}
                />
                {createErr && <p className="text-error text-xs">{createErr}</p>}
                <button type="submit" disabled={creating || !content.trim()}
                  className="w-full py-2.5 text-sm font-bold rounded-lg text-on-secondary disabled:opacity-50 transition-all hover:brightness-110 flex items-center justify-center gap-2"
                  style={{ background: '#4cd7f6', color: '#003640' }}>
                  <span className="material-symbols-outlined text-sm">add</span>
                  {creating ? 'Committing…' : 'Commit to Memory'}
                </button>
              </form>
              <div className="mt-6 p-4 rounded-lg border-l-2 border-secondary/50 bg-surface-container-lowest">
                <p className="text-[11px] text-on-surface-variant leading-relaxed">
                  <span className="text-secondary font-bold">Pro-tip:</span> Use natural language for better semantic retrieval.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Memory list */}
        <div className="lg:col-span-2 space-y-6">
          <div className="flex justify-between items-end pb-4" style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
            <h2 className="font-headline text-xl font-bold">
              Recent Entries
              <span className="text-on-surface-variant/40 font-normal text-sm ml-2">({memories.length} total)</span>
            </h2>
          </div>

          {loading ? (
            <div className="text-on-surface/40 text-sm flex items-center gap-2">
              <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
            </div>
          ) : error ? (
            <div className="text-error text-sm">{error}</div>
          ) : memories.length === 0 ? (
            <div className="bg-surface-container-low rounded-xl p-12 text-center">
              <span className="material-symbols-outlined text-on-surface-variant/30 mb-3" style={{ fontSize: '2.5rem' }}>memory</span>
              <p className="text-on-surface-variant text-sm">No memories stored yet.</p>
            </div>
          ) : (
            <div className="space-y-4">
              {memories.map(m => (
                <div key={m.id} className="group bg-surface-container hover:bg-surface-container-high transition-all p-6 rounded-xl"
                  style={{ border: '1px solid rgba(70,69,84,0.05)' }}>
                  {editId === m.id ? (
                    <div className="flex flex-col gap-3">
                      <textarea rows={3}
                        className="w-full bg-surface-container-lowest rounded-lg p-3 text-sm text-on-surface focus:ring-1 focus:ring-primary/50 focus:outline-none resize-none"
                        style={{ border: '1px solid rgba(70,69,84,0.2)' }}
                        value={editContent}
                        onChange={e => setEditContent(e.target.value)}
                      />
                      <div className="flex gap-2">
                        <button onClick={() => handleUpdate(m.id)}
                          className="px-4 py-1.5 text-xs font-bold rounded text-on-primary"
                          style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>Save</button>
                        <button onClick={() => setEditId(null)}
                          className="px-4 py-1.5 text-xs font-bold rounded text-on-surface-variant hover:text-on-surface transition-colors"
                          style={{ border: '1px solid rgba(70,69,84,0.2)' }}>Cancel</button>
                      </div>
                    </div>
                  ) : (
                    <>
                      <div className="flex justify-between items-start mb-4">
                        <div className="flex items-center gap-3">
                          <span className="flex items-center gap-1.5 px-2 py-0.5 rounded bg-tertiary-container/30 text-tertiary text-[10px] font-bold uppercase tracking-wider">
                            <span className="w-1.5 h-1.5 rounded-full bg-tertiary" /> Embedded
                          </span>
                          <span className="text-[10px] text-on-surface-variant/50 font-mono">{m.id?.slice(0, 16)}…</span>
                        </div>
                        <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button onClick={() => { setEditId(m.id); setEditContent(m.content) }}
                            className="p-1.5 hover:bg-surface-variant rounded text-on-surface-variant transition-colors">
                            <span className="material-symbols-outlined text-sm">edit</span>
                          </button>
                          <button onClick={() => handleDelete(m.id)}
                            className="p-1.5 hover:bg-surface-variant rounded text-error/70 transition-colors">
                            <span className="material-symbols-outlined text-sm">delete</span>
                          </button>
                        </div>
                      </div>
                      <p className="text-on-surface leading-relaxed text-[0.9375rem]">{m.content}</p>
                      <div className="mt-5 flex items-center justify-between text-[11px] text-on-surface-variant/60 font-medium">
                        <div className="flex items-center gap-4">
                          <span className="flex items-center gap-1">
                            <span className="material-symbols-outlined text-sm">schedule</span>
                            {m.created_at ? new Date(m.created_at).toLocaleString() : ''}
                          </span>
                          {m.tags?.length > 0 && (
                            <span className="flex items-center gap-1">
                              <span className="material-symbols-outlined text-sm">label</span>
                              {m.tags.map(t => `#${t}`).join(' ')}
                            </span>
                          )}
                        </div>
                      </div>
                    </>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
