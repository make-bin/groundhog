import React, { useState, useEffect } from 'react'
import api from '../api/gateway.js'

export default function Memory() {
  const [memories, setMemories] = useState([])
  const [loading, setLoading]   = useState(true)
  const [error, setError]       = useState(null)
  const [content, setContent]   = useState('')
  const [creating, setCreating] = useState(false)
  const [createErr, setCreateErr] = useState(null)
  const [query, setQuery]       = useState('')
  const [results, setResults]   = useState(null)
  const [searching, setSearching] = useState(false)
  const [editId, setEditId]     = useState(null)
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
    try {
      await api.memories.create(content.trim())
      setContent('')
      load()
    } catch (e) {
      setCreateErr(e.message)
    } finally {
      setCreating(false)
    }
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
    try {
      const res = await api.memories.search(query.trim(), 10)
      setResults(Array.isArray(res) ? res : [])
    } catch (e) {
      setResults([])
    } finally {
      setSearching(false)
    }
  }

  async function handleUpdate(id) {
    await api.memories.update(id, editContent).catch(() => {})
    setEditId(null)
    load()
  }

  return (
    <div className="p-8 max-w-4xl">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Memory</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 mb-6">
        {/* Create */}
        <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm">
          <h2 className="font-semibold text-gray-700 text-sm mb-3">Save Memory</h2>
          <form onSubmit={handleCreate} className="flex flex-col gap-2">
            <textarea rows={3} className="border rounded-lg px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="Enter memory content…"
              value={content} onChange={e => setContent(e.target.value)} />
            {createErr && <p className="text-red-500 text-xs">{createErr}</p>}
            <button type="submit" disabled={creating || !content.trim()}
              className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm hover:bg-blue-700 disabled:opacity-50 self-start">
              {creating ? 'Saving…' : 'Save'}
            </button>
          </form>
        </div>

        {/* Search */}
        <div className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm">
          <h2 className="font-semibold text-gray-700 text-sm mb-3">Search Memory</h2>
          <form onSubmit={handleSearch} className="flex gap-2 mb-3">
            <input className="flex-1 border rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="Search query…"
              value={query} onChange={e => setQuery(e.target.value)} />
            <button type="submit" disabled={searching || !query.trim()}
              className="bg-purple-600 text-white px-4 py-1.5 rounded text-sm hover:bg-purple-700 disabled:opacity-50">
              {searching ? '…' : 'Search'}
            </button>
          </form>
          {results !== null && (
            results.length === 0 ? (
              <p className="text-gray-400 text-xs">No results.</p>
            ) : (
              <div className="space-y-2 max-h-48 overflow-y-auto">
                {results.map((r, i) => {
                  const m = r.memory ?? r
                  return (
                    <div key={m.id ?? i} className="bg-purple-50 rounded-lg p-2.5 text-xs">
                      <div className="text-gray-700 mb-1">{m.content}</div>
                      <div className="text-purple-500">score: {(r.score ?? r.hybrid_score ?? 0).toFixed(3)}</div>
                    </div>
                  )
                })}
              </div>
            )
          )}
        </div>
      </div>

      {/* Memory list */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <div className="px-5 py-3 border-b border-gray-100 text-sm font-semibold text-gray-700">
          All Memories ({memories.length})
        </div>
        {loading ? (
          <div className="p-5 text-gray-400 text-sm">Loading...</div>
        ) : error ? (
          <div className="p-5 text-red-500 text-sm">{error}</div>
        ) : memories.length === 0 ? (
          <div className="p-5 text-gray-400 text-sm">No memories stored.</div>
        ) : (
          <div className="divide-y divide-gray-50">
            {memories.map(m => (
              <div key={m.id} className="px-5 py-3 hover:bg-gray-50 group">
                {editId === m.id ? (
                  <div className="flex gap-2 items-start">
                    <textarea rows={2} className="flex-1 border rounded px-2 py-1 text-sm resize-none"
                      value={editContent} onChange={e => setEditContent(e.target.value)} />
                    <button onClick={() => handleUpdate(m.id)}
                      className="bg-blue-600 text-white px-3 py-1 rounded text-xs">Save</button>
                    <button onClick={() => setEditId(null)}
                      className="border px-3 py-1 rounded text-xs">Cancel</button>
                  </div>
                ) : (
                  <div className="flex items-start gap-3">
                    <div className="flex-1">
                      <p className="text-sm text-gray-700">{m.content}</p>
                      <p className="text-xs text-gray-400 mt-0.5">
                        {m.created_at ? new Date(m.created_at).toLocaleString() : ''}
                        {m.tags?.length > 0 && (
                          <span className="ml-2">{m.tags.map(t => `#${t}`).join(' ')}</span>
                        )}
                      </p>
                    </div>
                    <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button onClick={() => { setEditId(m.id); setEditContent(m.content) }}
                        className="text-blue-400 hover:text-blue-600 text-xs">Edit</button>
                      <button onClick={() => handleDelete(m.id)}
                        className="text-red-400 hover:text-red-600 text-xs">Delete</button>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
