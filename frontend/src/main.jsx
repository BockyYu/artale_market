import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import App from './App.jsx'
import MemberAuth from './MemberAuth.jsx'
import Login from './admin/Login.jsx'
import Layout from './admin/Layout.jsx'
import Admins from './admin/Admins.jsx'
import Members from './admin/Members.jsx'
import './App.css'

function RequireAuth({ children }) {
  const token = localStorage.getItem('admin_token')
  if (!token) return <Navigate to="/admin/login" replace />
  return children
}

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        {/* 前台 */}
        <Route path="/*" element={<App />} />
        <Route path="/login" element={<MemberAuth />} />

        {/* 後台登入 */}
        <Route path="/admin/login" element={<Login />} />

        {/* 後台（JWT 保護） */}
        <Route
          path="/admin"
          element={
            <RequireAuth>
              <Layout />
            </RequireAuth>
          }
        >
          <Route index element={<Navigate to="admins" replace />} />
          <Route path="admins" element={<Admins />} />
          <Route path="members" element={<Members />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </React.StrictMode>,
)
