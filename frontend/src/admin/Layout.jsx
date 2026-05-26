import { NavLink, useNavigate, Outlet } from 'react-router-dom'
import { logout, currentUser } from './api'
import './admin.css'

const NAV = [
  { to: '/admin/admins', label: '管理員帳號' },
  { to: '/admin/members', label: '會員列表' },
  { to: '/admin/items', label: '道具列表' },
  { to: '/admin/bots', label: '通知機器人' },
  { to: '/admin/alerts', label: '價格提醒' },
]

export default function Layout() {
  const navigate = useNavigate()
  const user = currentUser()

  function handleLogout() {
    logout()
    navigate('/admin/login')
  }

  return (
    <div className="admin-layout">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h2>後台管理</h2>
          <p>{user?.username} · {user?.role}</p>
        </div>
        <nav className="sidebar-nav">
          {NAV.map(n => (
            <NavLink
              key={n.to}
              to={n.to}
              className={({ isActive }) => `nav-item${isActive ? ' active' : ''}`}
            >
              {n.label}
            </NavLink>
          ))}
        </nav>
        <div className="sidebar-footer">
          <button className="btn-logout" onClick={handleLogout}>登出</button>
        </div>
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  )
}
