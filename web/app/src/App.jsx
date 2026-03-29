import React, { useEffect, useState } from 'react'
import { BrowserRouter, Routes, Route, NavLink, Navigate } from 'react-router-dom'
import api from './api/gateway.js'

const Dashboard = React.lazy(() => import('./pages/Dashboard.jsx'))
const Sessions  = React.lazy(() => import('./pages/Sessions.jsx'))
const Chat      = React.lazy(() => import('./pages/Chat.jsx'))
const Channels  = React.lazy(() => import('./pages/Channels.jsx'))
const Agents    = React.lazy(() => import('./pages/Agents.jsx'))
const Cron      = React.lazy(() => import('./pages/Cron.jsx'))
const Memory    = React.lazy(() => import('./pages/Memory.jsx'))
const Security  = React.lazy(() => import('./pages/Security.jsx'))
const Config    = React.lazy(() => import('./pages/Config.jsx'))

const NAV = [
  { to: '/',         icon: 'dashboard',  label: 'Dashboard'  },
  { to: '/agents',   icon: 'smart_toy',  label: 'Agents'     },
  { to: '/sessions', icon: 'forum',      label: 'Sessions'   },
  { to: '/channels', icon: 'hub',        label: 'Channels'   },
  { to: '/cron',     icon: 'schedule',   label: 'Cron Jobs'  },
  { to: '/memory',   icon: 'memory',     label: 'Memory'     },
  { to: '/security', icon: 'shield',     label: 'Security'   },
  { to: '/config',   icon: 'settings',   label: 'Config'     },
]

function Sidebar() {
  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-surface-container-low flex flex-col z-50"
      style={{ borderRight: '1px solid rgba(70,69,84,0.2)' }}>
      <div className="px-6 py-8">
        <h1 className="font-headline text-primary font-bold text-xl tracking-tighter">Groundhog AI</h1>
        <p className="text-on-surface/40 text-[10px] uppercase tracking-widest mt-1">v1.0 · stable</p>
      </div>
      <nav className="flex-1 px-2 space-y-0.5">
        {NAV.map(({ to, icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              isActive
                ? 'flex items-center gap-3 py-3 px-4 text-secondary border-l-2 border-secondary bg-surface-container text-sm font-medium transition-all'
                : 'flex items-center gap-3 py-3 px-4 text-on-surface/50 hover:bg-surface-container hover:text-on-surface text-sm font-medium transition-all'
            }
          >
            <span className="material-symbols-outlined">{icon}</span>
            <span className="font-body tracking-wide">{label}</span>
          </NavLink>
        ))}
      </nav>
      <div className="p-4 space-y-1" style={{ borderTop: '1px solid rgba(70,69,84,0.15)' }}>
        <a href="/api/v1/health" target="_blank"
          className="flex items-center gap-3 py-2 px-4 text-on-surface/40 hover:text-on-surface text-xs transition-colors">
          <span className="material-symbols-outlined text-sm">health_metrics</span>
          <span>Health Check</span>
        </a>
      </div>
    </aside>
  )
}

function StatusBar({ health }) {
  const checks = [
    { label: 'API',      ok: health?.status === 'ok' || health !== null },
    { label: 'Database', ok: health?.database === 'ok' },
    { label: 'Redis',    ok: health?.redis === 'ok' || health?.redis?.startsWith?.('ok') },
  ]
  return (
    <footer className="fixed bottom-0 left-0 w-full flex justify-between items-center px-6 py-2 z-40"
      style={{ background: '#060e20', borderTop: '1px solid rgba(70,69,84,0.2)' }}>
      <div className="flex items-center gap-6">
        {checks.map(({ label, ok }) => (
          <div key={label} className="flex items-center gap-2">
            <span className={`w-1.5 h-1.5 rounded-full ${ok ? 'bg-tertiary pulse-dot' : 'bg-error'}`} />
            <span className={`font-body text-[0.6875rem] uppercase tracking-widest font-semibold ${ok ? 'text-tertiary' : 'text-error'}`}>
              {label}: {ok ? 'Online' : 'Offline'}
            </span>
          </div>
        ))}
      </div>
      <div className="flex items-center gap-6">
        <span className="font-body text-[0.6875rem] uppercase tracking-widest font-semibold text-on-surface/40">
          WebSocket: Connected
        </span>
      </div>
    </footer>
  )
}

export default function App() {
  const [health, setHealth] = useState(null)

  useEffect(() => {
    api.health().then(setHealth).catch(() => setHealth({ status: 'error' }))
    const t = setInterval(() => {
      api.health().then(setHealth).catch(() => {})
    }, 30000)
    return () => clearInterval(t)
  }, [])

  return (
    <BrowserRouter>
      <div className="flex h-screen overflow-hidden bg-surface">
        <Sidebar />
        <main className="flex-1 ml-64 overflow-auto pb-10">
          <React.Suspense fallback={
            <div className="flex items-center justify-center h-full text-on-surface/40 text-sm">
              <span className="material-symbols-outlined animate-spin mr-2">autorenew</span>
              Loading...
            </div>
          }>
            <Routes>
              <Route path="/"                  element={<Dashboard />} />
              <Route path="/agents"            element={<Agents />} />
              <Route path="/sessions"          element={<Sessions />} />
              <Route path="/sessions/:id/chat" element={<Chat />} />
              <Route path="/channels"          element={<Channels />} />
              <Route path="/cron"              element={<Cron />} />
              <Route path="/memory"            element={<Memory />} />
              <Route path="/security"          element={<Security />} />
              <Route path="/config"            element={<Config />} />
              <Route path="*"                  element={<Navigate to="/" replace />} />
            </Routes>
          </React.Suspense>
        </main>
        <StatusBar health={health} />
      </div>
    </BrowserRouter>
  )
}
