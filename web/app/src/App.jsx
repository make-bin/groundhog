import React from 'react'
import { BrowserRouter, Routes, Route, NavLink, Navigate } from 'react-router-dom'

const Dashboard  = React.lazy(() => import('./pages/Dashboard.jsx'))
const Sessions   = React.lazy(() => import('./pages/Sessions.jsx'))
const Chat       = React.lazy(() => import('./pages/Chat.jsx'))
const Channels   = React.lazy(() => import('./pages/Channels.jsx'))
const Memory     = React.lazy(() => import('./pages/Memory.jsx'))
const Security   = React.lazy(() => import('./pages/Security.jsx'))
const Config     = React.lazy(() => import('./pages/Config.jsx'))
const Cron       = React.lazy(() => import('./pages/Cron.jsx'))

const NAV = [
  { to: '/',         label: '📊 Dashboard'  },
  { to: '/sessions', label: '💬 Sessions'   },
  { to: '/channels', label: '📡 Channels'   },
  { to: '/cron',     label: '⏰ Cron Jobs'  },
  { to: '/memory',   label: '🧠 Memory'     },
  { to: '/security', label: '🔒 Security'   },
  { to: '/config',   label: '⚙️  Config'    },
]

function Sidebar() {
  return (
    <nav className="w-52 shrink-0 bg-gray-900 text-white flex flex-col">
      <div className="px-5 py-4 border-b border-gray-700">
        <span className="text-lg font-bold tracking-tight">🦔 Groundhog</span>
      </div>
      <div className="flex flex-col gap-0.5 p-3 flex-1">
        {NAV.map(({ to, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              `px-3 py-2 rounded-md text-sm transition-colors ${
                isActive ? 'bg-blue-600 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white'
              }`
            }
          >
            {label}
          </NavLink>
        ))}
      </div>
      <div className="px-4 py-3 border-t border-gray-700 text-xs text-gray-500">
        v1.0 · <a href="/api/v1/health" target="_blank" className="hover:text-gray-300">health</a>
      </div>
    </nav>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <div className="flex h-screen overflow-hidden">
        <Sidebar />
        <main className="flex-1 overflow-auto bg-gray-50">
          <React.Suspense fallback={
            <div className="flex items-center justify-center h-full text-gray-400">Loading...</div>
          }>
            <Routes>
              <Route path="/"              element={<Dashboard />} />
              <Route path="/sessions"      element={<Sessions />} />
              <Route path="/sessions/:id/chat" element={<Chat />} />
              <Route path="/channels"      element={<Channels />} />
              <Route path="/cron"          element={<Cron />} />
              <Route path="/memory"        element={<Memory />} />
              <Route path="/security"      element={<Security />} />
              <Route path="/config"        element={<Config />} />
              <Route path="*"              element={<Navigate to="/" replace />} />
            </Routes>
          </React.Suspense>
        </main>
      </div>
    </BrowserRouter>
  )
}
