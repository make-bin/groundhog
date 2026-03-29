import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/gateway.js'

const PROVIDER_COLORS = {
  openai_compat: 'text-secondary',
  openai:        'text-tertiary',
  anthropic:     'text-primary',
  ollama:        'text-on-surface-variant',
  groq:          'text-primary-container',
}

export default function Agents() {
  const navigate = useNavigate()
  const [agents, setAgents]   = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState(null)

  useEffect(() => {
    api.agents.list()
      .then(d => setAgents(Array.isArray(d) ? d : []))
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  function startChat(agentID) {
    // Navigate to sessions with agent pre-selected
    navigate('/sessions', { state: { agentID } })
  }

  return (
    <div className="p-10 max-w-5xl">
      <div className="mb-10">
        <h2 className="font-headline text-3xl font-black text-on-surface tracking-tighter">Agents</h2>
        <p className="text-on-surface-variant text-sm mt-1">
          Configured AI agents — each with its own model, skills, and routing rules
        </p>
      </div>

      {loading ? (
        <div className="text-on-surface/40 text-sm flex items-center gap-2">
          <span className="material-symbols-outlined animate-spin">autorenew</span> Loading…
        </div>
      ) : error ? (
        <div className="text-error text-sm">{error}</div>
      ) : agents.length === 0 ? (
        <div className="bg-surface-container-low rounded-xl p-16 text-center">
          <span className="material-symbols-outlined text-on-surface-variant/30 mb-4" style={{ fontSize: '3rem' }}>smart_toy</span>
          <p className="text-on-surface-variant text-sm mb-2">No agents configured.</p>
          <p className="text-on-surface-variant/50 text-xs">Add agents to <code className="font-mono bg-surface-container px-1 rounded">configs/config.yaml</code> under the <code className="font-mono bg-surface-container px-1 rounded">agents.list</code> section.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {agents.map(agent => (
            <div key={agent.id}
              className="group bg-surface-container rounded-xl p-6 hover:bg-surface-container-high transition-all relative overflow-hidden"
              style={{ border: '1px solid rgba(70,69,84,0.1)' }}>

              {/* Default badge */}
              {agent.is_default && (
                <div className="absolute top-4 right-4">
                  <span className="text-[10px] font-bold uppercase tracking-widest text-tertiary bg-tertiary-container/20 px-2 py-0.5 rounded-full">
                    Default
                  </span>
                </div>
              )}

              {/* Header */}
              <div className="flex items-start gap-4 mb-5">
                <div className="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center shrink-0">
                  <span className="material-symbols-outlined text-primary" style={{ fontVariationSettings: "'FILL' 1" }}>smart_toy</span>
                </div>
                <div className="min-w-0">
                  <div className="flex items-center gap-2">
                    <h3 className="font-headline font-bold text-on-surface">{agent.name}</h3>
                  </div>
                  <code className="text-[10px] font-mono text-on-surface-variant/50 uppercase tracking-tighter">{agent.id}</code>
                </div>
              </div>

              {/* Description */}
              {agent.description && (
                <p className="text-sm text-on-surface-variant leading-relaxed mb-5">{agent.description}</p>
              )}

              {/* Model info */}
              <div className="space-y-2 mb-5">
                {agent.model && (
                  <div className="flex items-center gap-2 text-xs">
                    <span className="material-symbols-outlined text-on-surface-variant/50 text-sm">model_training</span>
                    <span className={`font-mono font-medium ${PROVIDER_COLORS[agent.provider] ?? 'text-on-surface-variant'}`}>
                      {agent.provider && <span className="text-on-surface-variant/50 mr-1">{agent.provider} /</span>}
                      {agent.model.split('/').pop()}
                    </span>
                  </div>
                )}
                {agent.skills?.length > 0 && (
                  <div className="flex items-center gap-2 text-xs">
                    <span className="material-symbols-outlined text-on-surface-variant/50 text-sm">psychology</span>
                    <div className="flex flex-wrap gap-1">
                      {agent.skills.map(s => (
                        <span key={s} className="bg-surface-container-lowest text-on-surface-variant px-1.5 py-0.5 rounded text-[10px] font-mono"
                          style={{ border: '1px solid rgba(70,69,84,0.2)' }}>
                          {s}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* Action */}
              <button
                onClick={() => startChat(agent.id)}
                className="w-full py-2 text-xs font-bold rounded text-on-primary transition-all hover:brightness-110 flex items-center justify-center gap-2"
                style={{ background: 'linear-gradient(135deg, #c0c1ff, #8083ff)' }}>
                <span className="material-symbols-outlined text-sm">chat</span>
                Start Session
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
