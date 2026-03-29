import React, { useState, useEffect, useRef } from 'react'
import { useParams, Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import api from '../api/gateway.js'

function ToolCallCard({ tool }) {
  const [open, setOpen] = useState(false)
  const running = !tool.done
  const border = running ? 'border-secondary' : tool.isError ? 'border-error' : 'border-tertiary'
  const iconColor = running ? 'text-secondary' : tool.isError ? 'text-error' : 'text-tertiary'
  const icon = running ? 'construction' : tool.isError ? 'cancel' : 'check_circle'

  return (
    <div className={`bg-surface-container-lowest rounded overflow-hidden my-2 border-l-2 ${border}`}
      style={{ border: `1px solid rgba(70,69,84,0.1)`, borderLeft: undefined }}>
      <button className="flex items-center gap-2 w-full px-4 py-2.5 text-left bg-surface-container/30"
        onClick={() => setOpen(o => !o)}>
        <span className={`material-symbols-outlined text-sm ${iconColor}`}>{icon}</span>
        <code className="font-mono text-xs font-bold text-on-surface-variant">{tool.tool_name}</code>
        {tool.duration_ms > 0 && <span className="text-on-surface-variant/40 text-[10px] ml-1">{tool.duration_ms}ms</span>}
        {running && <span className="text-secondary text-[10px] ml-1 animate-pulse">running…</span>}
        <span className="ml-auto text-on-surface-variant/40 text-[10px]">{open ? '▲' : '▼'}</span>
      </button>
      {open && (
        <div className="px-4 py-3 space-y-2" style={{ borderTop: '1px solid rgba(70,69,84,0.1)' }}>
          {tool.args && Object.keys(tool.args).length > 0 && (
            <div>
              <div className="text-[10px] font-bold uppercase tracking-wide text-on-surface-variant/50 mb-1">Input</div>
              <pre className="bg-surface/50 rounded p-2 overflow-auto max-h-32 text-[11px] leading-relaxed text-on-surface-variant">
                {JSON.stringify(tool.args, null, 2)}
              </pre>
            </div>
          )}
          {tool.result && (
            <div>
              <div className={`text-[10px] font-bold uppercase tracking-wide mb-1 ${tool.isError ? 'text-error' : 'text-on-surface-variant/50'}`}>
                {tool.isError ? 'Error' : 'Output'}
              </div>
              <pre className={`rounded p-2 overflow-auto max-h-40 text-[11px] leading-relaxed ${tool.isError ? 'bg-error-container/10 text-error' : 'bg-surface/50 text-on-surface-variant'}`}>
                {tool.result.length > 800 ? tool.result.slice(0, 800) + '\n…[truncated]' : tool.result}
              </pre>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function ApprovalModal({ approval, sessionId, onDone }) {
  const [deciding, setDeciding] = useState(false)
  async function decide(decision) {
    setDeciding(true)
    try { await api.sessions.resolveApproval(sessionId, approval.approval_id, decision) } catch {}
    onDone(decision)
  }
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4" style={{ background: 'rgba(11,19,38,0.85)', backdropFilter: 'blur(8px)' }}>
      <div className="bg-surface-container w-full max-w-md rounded-2xl shadow-2xl overflow-hidden" style={{ border: '1px solid rgba(70,69,84,0.3)' }}>
        <div className="bg-error-container/10 px-6 py-5 flex items-center gap-3" style={{ borderBottom: '1px solid rgba(255,180,171,0.15)' }}>
          <div className="w-10 h-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
            <span className="material-symbols-outlined text-error">lock_person</span>
          </div>
          <div>
            <h3 className="font-headline font-bold text-error text-sm uppercase tracking-wide">Approval Required</h3>
            <p className="text-xs text-on-surface/60 mt-0.5">Agent wants to execute a tool</p>
          </div>
        </div>
        <div className="px-6 py-5">
          <div className="flex items-center gap-2 mb-4">
            <span className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wide">Tool</span>
            <code className="bg-error-container/20 text-error px-2 py-0.5 rounded text-xs font-mono font-bold"
              style={{ border: '1px solid rgba(255,180,171,0.2)' }}>
              {approval.tool_name}
            </code>
          </div>
          {approval.args && Object.keys(approval.args).length > 0 && (
            <div>
              <span className="text-[10px] font-bold text-on-surface-variant uppercase tracking-wide block mb-2">Arguments</span>
              <pre className="bg-surface-container-lowest rounded-lg p-3 text-xs font-mono overflow-auto max-h-48 text-on-surface-variant leading-relaxed"
                style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
                {JSON.stringify(approval.args, null, 2)}
              </pre>
            </div>
          )}
        </div>
        <div className="flex gap-3 px-6 pb-6">
          <button onClick={() => decide('deny')} disabled={deciding}
            className="flex-1 py-2.5 rounded-xl text-sm font-bold text-on-surface-variant hover:bg-surface-container-high transition-colors disabled:opacity-50"
            style={{ border: '1px solid rgba(70,69,84,0.3)' }}>
            ✕ Deny
          </button>
          <button onClick={() => decide('approve')} disabled={deciding}
            className="flex-1 py-2.5 rounded-xl text-sm font-bold text-on-primary disabled:opacity-50 transition-all hover:brightness-110"
            style={{ background: 'linear-gradient(135deg, #4edea3, #00885d)' }}>
            {deciding ? '…' : '✓ Approve'}
          </button>
        </div>
      </div>
    </div>
  )
}

const mdComponents = {
  p:    ({ children }) => <p className="mb-2 last:mb-0">{children}</p>,
  h1:   ({ children }) => <h1 className="text-lg font-bold mt-3 mb-1 font-headline">{children}</h1>,
  h2:   ({ children }) => <h2 className="text-base font-bold mt-3 mb-1 font-headline">{children}</h2>,
  h3:   ({ children }) => <h3 className="text-sm font-semibold mt-2 mb-1">{children}</h3>,
  ul:   ({ children }) => <ul className="list-disc pl-5 mb-2 space-y-0.5">{children}</ul>,
  ol:   ({ children }) => <ol className="list-decimal pl-5 mb-2 space-y-0.5">{children}</ol>,
  li:   ({ children }) => <li className="leading-relaxed">{children}</li>,
  strong: ({ children }) => <strong className="font-semibold text-on-surface">{children}</strong>,
  blockquote: ({ children }) => <blockquote className="border-l-2 border-secondary pl-3 my-2 text-on-surface-variant italic">{children}</blockquote>,
  a:    ({ href, children }) => <a href={href} target="_blank" rel="noopener noreferrer" className="text-secondary underline hover:text-primary">{children}</a>,
  code: ({ inline, children }) => inline
    ? <code className="bg-surface-container-high text-secondary rounded px-1 py-0.5 text-[12px] font-mono">{children}</code>
    : <code className="block bg-surface-container-lowest text-on-surface-variant rounded-lg p-3 my-2 text-[12px] font-mono overflow-x-auto whitespace-pre">{children}</code>,
  pre:  ({ children }) => <>{children}</>,
  table: ({ children }) => <div className="overflow-x-auto my-2"><table className="text-xs border-collapse w-full">{children}</table></div>,
  thead: ({ children }) => <thead className="bg-surface-container-low">{children}</thead>,
  th:   ({ children }) => <th className="px-2 py-1 text-left font-semibold text-on-surface-variant" style={{ border: '1px solid rgba(70,69,84,0.2)' }}>{children}</th>,
  td:   ({ children }) => <td className="px-2 py-1 text-on-surface-variant" style={{ border: '1px solid rgba(70,69,84,0.2)' }}>{children}</td>,
}

function Message({ msg }) {
  const isUser   = msg.role === 'user'
  const isSystem = msg.role === 'system'
  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-6 px-2`}>
      {!isUser && (
        <div className="w-7 h-7 rounded-full bg-secondary flex items-center justify-center mr-3 mt-0.5 shrink-0">
          <span className="material-symbols-outlined text-on-secondary" style={{ fontSize: '14px', fontVariationSettings: "'FILL' 1" }}>smart_toy</span>
        </div>
      )}
      <div className={`max-w-[80%] min-w-0 ${!isUser ? 'flex-1' : ''}`}>
        {!isUser && msg.toolCalls?.length > 0 && (
          <div className="mb-2 space-y-1">
            {msg.toolCalls.map(tc => <ToolCallCard key={tc.tool_call_id} tool={tc} />)}
          </div>
        )}
        {(msg.content || msg.streaming) && (
          <div className={`rounded-xl px-5 py-3 text-sm leading-relaxed ${
            isUser
              ? 'bg-surface-container-highest text-on-surface rounded-tr-none ml-auto max-w-[75%]'
              : isSystem
                ? 'bg-error-container/10 text-error'
                : 'glass-panel text-on-surface rounded-tl-none'
          }`} style={!isUser && !isSystem ? { borderLeft: '2px solid rgba(76,215,246,0.4)' } : {}}>
            {isUser || isSystem ? (
              <span style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                {msg.content}
                {msg.streaming && <span className="cursor" aria-hidden />}
              </span>
            ) : (
              <div className="prose-sm max-w-none">
                <ReactMarkdown remarkPlugins={[remarkGfm]} components={mdComponents}>
                  {msg.content}
                </ReactMarkdown>
                {msg.streaming && <span className="cursor" aria-hidden />}
              </div>
            )}
          </div>
        )}
      </div>
      {isUser && (
        <div className="w-7 h-7 rounded-full bg-surface-container-highest flex items-center justify-center ml-3 mt-0.5 shrink-0">
          <span className="material-symbols-outlined text-on-surface-variant" style={{ fontSize: '14px' }}>person</span>
        </div>
      )}
    </div>
  )
}

export default function Chat() {
  const { id: sessionId } = useParams()
  const [session, setSession]           = useState(null)
  const [messages, setMessages]         = useState([])
  const [approvalQueue, setApprovalQueue] = useState([])
  const [status, setStatus]             = useState('idle')
  const [input, setInput]               = useState('')
  const [sending, setSending]           = useState(false)
  const bottomRef  = useRef(null)
  const sendingRef = useRef(false)
  const abortRef   = useRef(null)

  useEffect(() => {
    api.sessions.get(sessionId).then(s => {
      setSession(s)
      const hist = []
      for (const t of s.turns ?? []) {
        hist.push({ id: `u-${t.id}`, role: 'user',      content: t.user_input, streaming: false })
        hist.push({ id: `a-${t.id}`, role: 'assistant', content: t.response,   streaming: false })
      }
      setMessages(hist)
    }).catch(() => {})
    return () => abortRef.current?.abort()
  }, [sessionId])

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  async function send() {
    const text = input.trim()
    if (!text || sendingRef.current) return
    sendingRef.current = true
    setSending(true)
    setStatus('thinking')
    setInput('')
    setApprovalQueue([])
    setMessages(prev => [...prev, { id: `u-${Date.now()}`, role: 'user', content: text, streaming: false }])
    const aiId = `ai-${Date.now()}`
    setMessages(prev => [...prev, { id: aiId, role: 'assistant', content: '', streaming: true }])
    const controller = new AbortController()
    abortRef.current = controller
    let res
    try {
      res = await fetch(`/api/v1/sessions/${sessionId}/messages/stream`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'X-User-ID': localStorage.getItem('user_id') || 'user-001' },
        body: JSON.stringify({ user_input: text }),
        signal: controller.signal,
      })
    } catch (err) {
      if (err.name !== 'AbortError') {
        setMessages(prev => prev.map(m => m.id === aiId ? { ...m, content: `Network error: ${err.message}`, streaming: false, role: 'system' } : m))
      }
      setStatus('idle'); setSending(false); sendingRef.current = false; return
    }
    if (!res.ok) {
      const errText = await res.text().catch(() => res.statusText)
      setMessages(prev => prev.map(m => m.id === aiId ? { ...m, content: `Error ${res.status}: ${errText}`, streaming: false, role: 'system' } : m))
      setStatus('idle'); setSending(false); sendingRef.current = false; return
    }
    const reader = res.body.getReader()
    const decoder = new TextDecoder()
    let buf = ''
    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        buf += decoder.decode(value, { stream: true })
        const lines = buf.split('\n')
        buf = lines.pop()
        for (const line of lines) {
          if (!line.startsWith('data: ')) continue
          let evt; try { evt = JSON.parse(line.slice(6)) } catch { continue }
          switch (evt.type) {
            case 'chunk':
              setMessages(prev => prev.map(m => m.id === aiId ? { ...m, content: m.content + evt.chunk } : m)); break
            case 'tool_start':
              setMessages(prev => prev.map(m => m.id !== aiId ? m : { ...m, toolCalls: [...(m.toolCalls || []), { ...evt.tool, done: false }] }))
              setStatus('executing'); break
            case 'tool_done':
              setMessages(prev => prev.map(m => m.id !== aiId ? m : { ...m, toolCalls: (m.toolCalls || []).map(tc => tc.tool_call_id === evt.tool.tool_call_id ? { ...tc, ...evt.tool, done: true } : tc) }))
              setStatus('thinking'); break
            case 'approval_required':
              setApprovalQueue(q => [...q, evt.approval]); setStatus('waiting'); break
            case 'done':
              setMessages(prev => prev.map(m => m.id === aiId ? { ...m, streaming: false } : m))
              setApprovalQueue([]); setStatus('idle'); setSending(false); sendingRef.current = false; return
            case 'error':
              setMessages(prev => prev.map(m => m.id === aiId ? { ...m, content: m.content || evt.error, streaming: false, role: 'system' } : m))
              setApprovalQueue([]); setStatus('idle'); setSending(false); sendingRef.current = false; return
          }
        }
      }
    } catch {}
    setStatus('idle'); setSending(false); sendingRef.current = false
  }

  function onKey(e) { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send() } }

  const statusMap = {
    idle:      { text: 'Ready',      cls: 'bg-tertiary-container/20 text-tertiary' },
    thinking:  { text: 'Thinking…',  cls: 'bg-secondary/10 text-secondary' },
    waiting:   { text: 'Waiting…',   cls: 'bg-error-container/20 text-error' },
    executing: { text: 'Executing…', cls: 'bg-primary/10 text-primary' },
  }
  const statusLabel = statusMap[status] || { text: status, cls: 'bg-surface-variant text-on-surface-variant' }

  return (
    <div className="flex flex-col h-screen bg-surface">
      {/* Header */}
      <header className="bg-surface-container-low shrink-0 px-8 h-16 flex items-center justify-between"
        style={{ borderBottom: '1px solid rgba(70,69,84,0.1)' }}>
        <div className="flex items-center gap-4">
          <Link to="/sessions" className="text-on-surface-variant hover:text-on-surface transition-colors flex items-center gap-1 text-sm">
            <span className="material-symbols-outlined text-sm">arrow_back</span> Sessions
          </Link>
          <div className="h-4 w-px bg-outline-variant/20" />
          <div className="w-9 h-9 rounded bg-primary/10 flex items-center justify-center text-primary">
            <span className="material-symbols-outlined" style={{ fontVariationSettings: "'FILL' 1" }}>terminal</span>
          </div>
          <div>
            <h2 className="font-headline font-bold text-base tracking-tight flex items-center gap-2">
              {session?.agent_id || 'Chat'}
              <span className="bg-surface-container-highest text-[10px] px-2 py-0.5 rounded-full font-mono text-secondary">LIVE</span>
            </h2>
            <p className="text-[10px] text-on-surface-variant font-mono">{sessionId}</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          {session?.active_model && (
            <span className="text-xs text-on-surface-variant hidden sm:block">{session.active_model.split('/').pop()}</span>
          )}
          <span className={`text-[10px] px-2.5 py-1 rounded-full font-bold uppercase tracking-wide ${statusLabel.cls}`}>
            {statusLabel.text}
          </span>
        </div>
      </header>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto py-8">
        {messages.length === 0 && !sending && (
          <div className="flex flex-col items-center justify-center h-full text-on-surface/30 select-none">
            <span className="material-symbols-outlined mb-3" style={{ fontSize: '3rem', fontVariationSettings: "'FILL' 1" }}>smart_toy</span>
            <p className="text-sm">Start a conversation</p>
          </div>
        )}
        {messages.map(msg => <Message key={msg.id} msg={msg} />)}
        {sending && messages[messages.length - 1]?.role === 'user' && (
          <div className="flex justify-start mb-6 px-2">
            <div className="w-7 h-7 rounded-full bg-secondary flex items-center justify-center mr-3 mt-0.5 shrink-0">
              <span className="material-symbols-outlined text-on-secondary" style={{ fontSize: '14px', fontVariationSettings: "'FILL' 1" }}>smart_toy</span>
            </div>
            <div className="glass-panel rounded-xl rounded-tl-none px-5 py-3 text-sm text-on-surface-variant shadow-xl">
              <span className="cursor" aria-hidden />
            </div>
          </div>
        )}
        {approvalQueue.length > 0 && (
          <div className="mx-4 mb-4 flex items-center gap-3 rounded-xl px-5 py-3 text-sm bg-error-container/10"
            style={{ border: '1px solid rgba(255,180,171,0.2)' }}>
            <span className="material-symbols-outlined text-error shrink-0">lock_person</span>
            <div className="flex-1 min-w-0">
              <span className="font-bold text-error">Waiting for approval</span>
              <span className="text-error/70 ml-2 text-xs">
                Tool: <code className="font-mono">{approvalQueue[0].tool_name}</code>
                {approvalQueue.length > 1 && ` (+${approvalQueue.length - 1} more)`}
              </span>
            </div>
          </div>
        )}
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <div className="bg-surface-container shrink-0 px-6 py-4" style={{ borderTop: '1px solid rgba(70,69,84,0.1)' }}>
        <div className="flex gap-3 items-end max-w-4xl mx-auto">
          <div className="flex-1 bg-surface-container-low rounded-xl flex items-end px-4 py-2.5 gap-3"
            style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
            <textarea
              value={input}
              onChange={e => setInput(e.target.value)}
              onKeyDown={onKey}
              placeholder="Type a message… (Enter to send, Shift+Enter for newline)"
              rows={2}
              disabled={sending}
              className="flex-1 resize-none bg-transparent border-none focus:ring-0 text-sm text-on-surface placeholder:text-on-surface-variant/30 disabled:text-on-surface/40"
            />
          </div>
          <div className="flex gap-2 shrink-0">
            {sending && (
              <button onClick={() => abortRef.current?.abort()}
                className="px-4 py-2.5 rounded-xl text-sm font-bold text-error flex items-center gap-1 hover:bg-error-container/10 transition-colors"
                style={{ border: '1px solid rgba(255,180,171,0.3)' }}>
                <span className="material-symbols-outlined text-sm">stop</span> Stop
              </button>
            )}
            <button onClick={send} disabled={!input.trim() || sending}
              className="px-5 py-2.5 rounded-xl text-sm font-bold text-on-primary disabled:opacity-40 transition-all hover:brightness-110 flex items-center gap-2"
              style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
              Send <span className="material-symbols-outlined text-sm">send</span>
            </button>
          </div>
        </div>
      </div>

      {approvalQueue[0] && (
        <ApprovalModal
          approval={approvalQueue[0]}
          sessionId={sessionId}
          onDone={(decision) => {
            if (decision === 'deny') {
              setMessages(prev => [...prev, {
                id: `deny-${Date.now()}`, role: 'system',
                content: `⛔ Tool "${approvalQueue[0].tool_name}" was denied.`, streaming: false,
              }])
            }
            setApprovalQueue(q => { const next = q.slice(1); if (next.length === 0) setStatus('thinking'); return next })
          }}
        />
      )}
    </div>
  )
}
