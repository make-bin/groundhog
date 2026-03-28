import React, { useState, useEffect, useRef } from 'react'
import { useParams, Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import api from '../api/gateway.js'

// ── Tool call card (inline in message stream) ─────────────────────────────────
function ToolCallCard({ tool }) {
  const [open, setOpen] = useState(false)
  const running = !tool.done

  const statusCls = running
    ? 'border-blue-200 bg-blue-50'
    : tool.isError
      ? 'border-red-200 bg-red-50'
      : 'border-green-200 bg-green-50'

  const icon = running ? '⚙️' : tool.isError ? '✗' : '✓'
  const iconCls = running
    ? 'text-blue-500'
    : tool.isError ? 'text-red-500' : 'text-green-600'

  return (
    <div className={`rounded-xl border text-xs my-2 overflow-hidden ${statusCls}`}>
      <button
        className="flex items-center gap-2 w-full px-3 py-2 text-left"
        onClick={() => setOpen(o => !o)}
      >
        <span className={`font-bold text-sm ${iconCls}`}>{icon}</span>
        <code className="font-mono font-semibold text-gray-700">{tool.tool_name}</code>
        {tool.duration_ms > 0 && (
          <span className="text-gray-400 ml-1">{tool.duration_ms}ms</span>
        )}
        {running && (
          <span className="ml-1 text-blue-400 animate-pulse">running…</span>
        )}
        <span className="ml-auto text-gray-400">{open ? '▲' : '▼'}</span>
      </button>

      {open && (
        <div className="border-t border-current/10 px-3 py-2 space-y-2">
          {tool.args && Object.keys(tool.args).length > 0 && (
            <div>
              <div className="text-[10px] font-semibold uppercase tracking-wide text-gray-500 mb-1">Input</div>
              <pre className="bg-white/70 rounded p-2 overflow-auto max-h-32 text-[11px] leading-relaxed text-gray-700">
                {JSON.stringify(tool.args, null, 2)}
              </pre>
            </div>
          )}
          {tool.result && (
            <div>
              <div className={`text-[10px] font-semibold uppercase tracking-wide mb-1 ${tool.isError ? 'text-red-500' : 'text-gray-500'}`}>
                {tool.isError ? 'Error' : 'Output'}
              </div>
              <pre className={`rounded p-2 overflow-auto max-h-40 text-[11px] leading-relaxed ${tool.isError ? 'bg-red-50 text-red-700' : 'bg-white/70 text-gray-700'}`}>
                {tool.result.length > 800 ? tool.result.slice(0, 800) + '\n…[truncated]' : tool.result}
              </pre>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

// ── Approval modal ────────────────────────────────────────────────────────────
function ApprovalModal({ approval, sessionId, onDone }) {
  const [deciding, setDeciding] = useState(false)

  async function decide(decision) {
    setDeciding(true)
    try {
      await api.sessions.resolveApproval(sessionId, approval.approval_id, decision)
    } catch (e) {
      console.error('resolveApproval failed', e)
    }
    onDone(decision)
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
        {/* Header */}
        <div className="flex items-center gap-3 px-5 py-4 border-b border-gray-100">
          <div className="w-9 h-9 rounded-full bg-amber-100 flex items-center justify-center text-lg shrink-0">
            🔐
          </div>
          <div>
            <h3 className="font-semibold text-gray-900 text-sm">Approval Required</h3>
            <p className="text-xs text-gray-500 mt-0.5">Agent wants to execute a tool</p>
          </div>
        </div>

        {/* Tool info */}
        <div className="px-5 py-4">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">Tool</span>
            <code className="bg-amber-50 text-amber-800 border border-amber-200 rounded px-2 py-0.5 text-xs font-mono font-semibold">
              {approval.tool_name}
            </code>
          </div>

          {approval.args && Object.keys(approval.args).length > 0 && (
            <div>
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wide block mb-1.5">Arguments</span>
              <pre className="bg-gray-50 border border-gray-200 rounded-lg p-3 text-xs font-mono overflow-auto max-h-48 text-gray-700 leading-relaxed">
                {JSON.stringify(approval.args, null, 2)}
              </pre>
            </div>
          )}
        </div>

        {/* Actions */}
        <div className="flex gap-2 px-5 pb-5">
          <button
            onClick={() => decide('deny')}
            disabled={deciding}
            className="flex-1 px-4 py-2.5 rounded-xl border border-gray-300 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            ✕ Deny
          </button>
          <button
            onClick={() => decide('approve')}
            disabled={deciding}
            className="flex-1 px-4 py-2.5 rounded-xl bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            {deciding ? '…' : '✓ Approve'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ── Single message bubble ─────────────────────────────────────────────────────
const mdComponents = {
  p:          ({ children }) => <p className="mb-2 last:mb-0">{children}</p>,
  h1:         ({ children }) => <h1 className="text-lg font-bold mt-3 mb-1">{children}</h1>,
  h2:         ({ children }) => <h2 className="text-base font-bold mt-3 mb-1">{children}</h2>,
  h3:         ({ children }) => <h3 className="text-sm font-semibold mt-2 mb-1">{children}</h3>,
  ul:         ({ children }) => <ul className="list-disc pl-5 mb-2 space-y-0.5">{children}</ul>,
  ol:         ({ children }) => <ol className="list-decimal pl-5 mb-2 space-y-0.5">{children}</ol>,
  li:         ({ children }) => <li className="leading-relaxed">{children}</li>,
  strong:     ({ children }) => <strong className="font-semibold">{children}</strong>,
  em:         ({ children }) => <em className="italic">{children}</em>,
  blockquote: ({ children }) => <blockquote className="border-l-4 border-gray-300 pl-3 my-2 text-gray-600 italic">{children}</blockquote>,
  hr:         () => <hr className="my-3 border-gray-200" />,
  a:          ({ href, children }) => <a href={href} target="_blank" rel="noopener noreferrer" className="text-blue-600 underline hover:text-blue-800">{children}</a>,
  code:       ({ inline, children }) => inline
    ? <code className="bg-gray-100 text-gray-800 rounded px-1 py-0.5 text-[12px] font-mono">{children}</code>
    : <code className="block bg-gray-900 text-gray-100 rounded-lg p-3 my-2 text-[12px] font-mono overflow-x-auto whitespace-pre">{children}</code>,
  pre:        ({ children }) => <>{children}</>,
  table:      ({ children }) => <div className="overflow-x-auto my-2"><table className="text-xs border-collapse w-full">{children}</table></div>,
  thead:      ({ children }) => <thead className="bg-gray-100">{children}</thead>,
  th:         ({ children }) => <th className="border border-gray-300 px-2 py-1 text-left font-semibold">{children}</th>,
  td:         ({ children }) => <td className="border border-gray-300 px-2 py-1">{children}</td>,
}

function Message({ msg }) {
  const isUser = msg.role === 'user'
  const isSystem = msg.role === 'system'
  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-3 px-2`}>
      {!isUser && (
        <div className="w-7 h-7 rounded-full bg-blue-600 text-white text-xs flex items-center justify-center mr-2 mt-0.5 shrink-0 font-bold">
          AI
        </div>
      )}
      <div className={`max-w-[80%] min-w-0 ${isUser ? '' : 'flex-1'}`}>
        {/* Tool call cards above the text bubble */}
        {!isUser && msg.toolCalls?.length > 0 && (
          <div className="mb-2 space-y-1">
            {msg.toolCalls.map(tc => <ToolCallCard key={tc.tool_call_id} tool={tc} />)}
          </div>
        )}

        {/* Text bubble — only render if there's content */}
        {(msg.content || msg.streaming) && (
          <div className={`rounded-2xl px-4 py-2.5 text-sm leading-relaxed shadow-sm ${
            isUser   ? 'bg-blue-600 text-white rounded-br-sm ml-auto max-w-[75%]' :
            isSystem ? 'bg-red-50 text-red-700 border border-red-200' :
                       'bg-white text-gray-800 border border-gray-200 rounded-bl-sm'
          }`}>
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
        <div className="w-7 h-7 rounded-full bg-gray-200 text-gray-600 text-xs flex items-center justify-center ml-2 mt-0.5 shrink-0 font-bold">
          U
        </div>
      )}
    </div>
  )
}

// ── Main Chat component ───────────────────────────────────────────────────────
export default function Chat() {
  const { id: sessionId } = useParams()
  const [session, setSession]       = useState(null)
  const [messages, setMessages]     = useState([])
  const [approvalQueue, setApprovalQueue] = useState([])  // queue of pending approvals
  const [status, setStatus]         = useState('idle')
  const [input, setInput]           = useState('')
  const [sending, setSending]       = useState(false)

  const bottomRef  = useRef(null)
  const sendingRef = useRef(false)   // avoid stale closure in send()
  const abortRef   = useRef(null)    // AbortController for current stream

  // Load session + history
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

    // Abort any in-flight stream on unmount
    return () => abortRef.current?.abort()
  }, [sessionId])

  // Auto-scroll
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  async function send() {
    const text = input.trim()
    if (!text || sendingRef.current) return

    sendingRef.current = true
    setSending(true)
    setStatus('thinking')
    setInput('')
    setApprovalQueue([])

    // Append user bubble immediately
    setMessages(prev => [...prev, {
      id: `u-${Date.now()}`, role: 'user', content: text, streaming: false,
    }])

    // Placeholder for AI response
    const aiId = `ai-${Date.now()}`
    setMessages(prev => [...prev, {
      id: aiId, role: 'assistant', content: '', streaming: true,
    }])

    const controller = new AbortController()
    abortRef.current = controller

    let res
    try {
      res = await fetch(`/api/v1/sessions/${sessionId}/messages/stream`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-User-ID': localStorage.getItem('user_id') || 'user-001',
        },
        body: JSON.stringify({ user_input: text }),
        signal: controller.signal,
      })
    } catch (err) {
      if (err.name !== 'AbortError') {
        setMessages(prev => prev.map(m =>
          m.id === aiId ? { ...m, content: `Network error: ${err.message}`, streaming: false, role: 'system' } : m
        ))
      }
      setStatus('idle')
      setSending(false)
      sendingRef.current = false
      return
    }

    if (!res.ok) {
      const errText = await res.text().catch(() => res.statusText)
      setMessages(prev => prev.map(m =>
        m.id === aiId ? { ...m, content: `Error ${res.status}: ${errText}`, streaming: false, role: 'system' } : m
      ))
      setStatus('idle')
      setSending(false)
      sendingRef.current = false
      return
    }

    // Read SSE stream
    const reader  = res.body.getReader()
    const decoder = new TextDecoder()
    let buf = ''

    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buf += decoder.decode(value, { stream: true })
        const lines = buf.split('\n')
        buf = lines.pop() // keep incomplete last line

        for (const line of lines) {
          if (!line.startsWith('data: ')) continue
          let evt
          try { evt = JSON.parse(line.slice(6)) } catch { continue }

          switch (evt.type) {
            case 'chunk':
              setMessages(prev => prev.map(m =>
                m.id === aiId ? { ...m, content: m.content + evt.chunk } : m
              ))
              break

            case 'tool_start':
              // Add a running tool card to the current AI message
              setMessages(prev => prev.map(m => {
                if (m.id !== aiId) return m
                const existing = m.toolCalls || []
                return { ...m, toolCalls: [...existing, { ...evt.tool, done: false }] }
              }))
              setStatus('executing')
              break

            case 'tool_done':
              // Update the matching tool card with result
              setMessages(prev => prev.map(m => {
                if (m.id !== aiId) return m
                const updated = (m.toolCalls || []).map(tc =>
                  tc.tool_call_id === evt.tool.tool_call_id
                    ? { ...tc, ...evt.tool, done: true }
                    : tc
                )
                return { ...m, toolCalls: updated }
              }))
              setStatus('thinking')
              break

            case 'approval_required':
              setApprovalQueue(q => [...q, evt.approval])
              setStatus('waiting')
              break

            case 'done':
              setMessages(prev => prev.map(m =>
                m.id === aiId ? { ...m, streaming: false } : m
              ))
              setApprovalQueue([])
              setStatus('idle')
              setSending(false)
              sendingRef.current = false
              return   // stream finished cleanly

            case 'error':
              setMessages(prev => prev.map(m =>
                m.id === aiId
                  ? { ...m, content: m.content || evt.error, streaming: false, role: 'system' }
                  : m
              ))
              setApprovalQueue([])
              setStatus('idle')
              setSending(false)
              sendingRef.current = false
              return
          }
        }
      }
    } catch (err) {
      if (err.name !== 'AbortError') {
        setMessages(prev => prev.map(m =>
          m.id === aiId ? { ...m, streaming: false } : m
        ))
      }
    }

    setStatus('idle')
    setSending(false)
    sendingRef.current = false
  }

  function onKey(e) {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send() }
  }

  const statusLabel = {
    idle:     { text: 'Ready',      cls: 'bg-green-100 text-green-700' },
    thinking: { text: 'Thinking…',  cls: 'bg-blue-100 text-blue-700'  },
    waiting:  { text: 'Waiting…',   cls: 'bg-amber-100 text-amber-700' },
    executing:{ text: 'Executing…', cls: 'bg-purple-100 text-purple-700' },
  }[status] || { text: status, cls: 'bg-gray-100 text-gray-600' }

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-5 py-3 flex items-center gap-3 shrink-0 shadow-sm">
        <Link to="/sessions" className="text-gray-400 hover:text-gray-600 text-sm transition-colors">
          ← Sessions
        </Link>
        <div className="h-4 w-px bg-gray-200" />
        <div className="min-w-0">
          <div className="text-sm font-semibold text-gray-800 truncate">
            {session?.agent_id || 'Chat'}
          </div>
          <div className="text-xs text-gray-400 font-mono truncate">{sessionId}</div>
        </div>
        <div className="ml-auto flex items-center gap-3 shrink-0">
          {session?.active_model && (
            <span className="text-xs text-gray-400 hidden sm:block">
              {session.active_model.split('/').pop()}
            </span>
          )}
          <span className={`text-xs px-2.5 py-1 rounded-full font-medium ${statusLabel.cls}`}>
            {statusLabel.text}
          </span>
        </div>
      </div>

      {/* Messages area */}
      <div className="flex-1 overflow-y-auto py-4">
        {messages.length === 0 && !sending && (
          <div className="flex flex-col items-center justify-center h-full text-gray-400 select-none">
            <div className="text-5xl mb-3">🦔</div>
            <p className="text-sm">Start a conversation</p>
          </div>
        )}

        {messages.map(msg => <Message key={msg.id} msg={msg} />)}

        {/* Thinking indicator when no AI bubble yet */}
        {sending && messages[messages.length - 1]?.role === 'user' && (
          <div className="flex justify-start mb-3 px-2">
            <div className="w-7 h-7 rounded-full bg-blue-600 text-white text-xs flex items-center justify-center mr-2 mt-0.5 shrink-0 font-bold">AI</div>
            <div className="bg-white border border-gray-200 rounded-2xl rounded-bl-sm px-4 py-2.5 text-sm text-gray-400 shadow-sm">
              <span className="cursor" aria-hidden />
            </div>
          </div>
        )}

        {/* Waiting for approval inline banner */}
        {approvalQueue.length > 0 && (
          <div className="mx-4 mb-3 flex items-center gap-3 bg-amber-50 border border-amber-200 rounded-xl px-4 py-3 text-sm">
            <span className="text-amber-500 text-base shrink-0">🔐</span>
            <div className="flex-1 min-w-0">
              <span className="font-medium text-amber-800">Waiting for your approval</span>
              <span className="text-amber-600 ml-2 text-xs">
                Tool: <code className="font-mono">{approvalQueue[0].tool_name}</code>
                {approvalQueue.length > 1 && ` (+${approvalQueue.length - 1} more)`}
              </span>
            </div>
          </div>
        )}

        <div ref={bottomRef} />
      </div>

      {/* Input bar */}
      <div className="bg-white border-t border-gray-200 px-4 py-3 shrink-0">
        <div className="flex gap-2 items-end max-w-4xl mx-auto">
          <textarea
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={onKey}
            placeholder="Type a message… (Enter to send, Shift+Enter for newline)"
            rows={2}
            disabled={sending}
            className="flex-1 resize-none border border-gray-300 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent disabled:bg-gray-50 disabled:text-gray-400 transition-colors"
          />
          <button
            onClick={send}
            disabled={!input.trim() || sending}
            className="bg-blue-600 text-white px-5 py-2.5 rounded-xl text-sm font-medium hover:bg-blue-700 active:bg-blue-800 disabled:opacity-40 disabled:cursor-not-allowed transition-colors shrink-0"
          >
            {sending ? '…' : 'Send'}
          </button>
        </div>
      </div>

      {approvalQueue[0] && (
        <ApprovalModal
          approval={approvalQueue[0]}
          sessionId={sessionId}
          onDone={(decision) => {
            // If denied, insert a system notice into the chat
            if (decision === 'deny') {
              setMessages(prev => [...prev, {
                id: `deny-${Date.now()}`,
                role: 'system',
                content: `⛔ Tool "${approvalQueue[0].tool_name}" was denied.`,
                streaming: false,
              }])
            }
            // Pop the front of the queue; if more pending, stay in waiting
            setApprovalQueue(q => {
              const next = q.slice(1)
              if (next.length === 0) setStatus('thinking')
              return next
            })
          }}
        />
      )}
    </div>
  )
}
