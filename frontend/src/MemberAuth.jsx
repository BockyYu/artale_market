import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { memberLogin, memberRegister } from './member-api'

export default function MemberAuth() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const [tab, setTab] = useState(params.get('tab') === 'register' ? 'register' : 'login')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)

  const [loginForm, setLoginForm] = useState({ username: '', password: '' })
  const [regForm, setRegForm] = useState({ nickname: '', username: '', password: '', email: '', invite_code: '' })

  async function handleLogin(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await memberLogin(loginForm.username, loginForm.password)
      navigate('/')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function handleRegister(e) {
    e.preventDefault()
    setError('')
    setSuccess('')
    setLoading(true)
    try {
      await memberRegister(regForm)
      setSuccess('註冊成功！請使用帳號登入')
      setTab('login')
      setLoginForm(f => ({ ...f, username: regForm.username }))
      setRegForm({ nickname: '', username: '', password: '', email: '', invite_code: '' })
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.page}>
      <div style={styles.card}>
        <h1 style={styles.title}>🏪 Artale Market</h1>

        <div style={styles.tabs}>
          <button
            style={{ ...styles.tab, ...(tab === 'login' ? styles.tabActive : {}) }}
            onClick={() => { setTab('login'); setError(''); setSuccess('') }}
          >
            登入
          </button>
          <button
            style={{ ...styles.tab, ...(tab === 'register' ? styles.tabActive : {}) }}
            onClick={() => { setTab('register'); setError(''); setSuccess('') }}
          >
            註冊
          </button>
        </div>

        {error && <div style={styles.error}>{error}</div>}
        {success && <div style={styles.successMsg}>{success}</div>}

        {tab === 'login' ? (
          <form onSubmit={handleLogin}>
            <Field label="帳號" value={loginForm.username}
              onChange={v => setLoginForm(f => ({ ...f, username: v }))} />
            <Field label="密碼" type="password" value={loginForm.password}
              onChange={v => setLoginForm(f => ({ ...f, password: v }))} />
            <button style={styles.btn} type="submit" disabled={loading}>
              {loading ? '登入中...' : '登入'}
            </button>
          </form>
        ) : (
          <div style={styles.devNotice}>
            <p style={styles.devIcon}>🚧</p>
            <p style={styles.devTitle}>Registration Unavailable</p>
            <p style={styles.devDesc}>
              Member registration is not available during the development phase.
              Please check back later.
            </p>
          </div>
        )}

        <button style={styles.backLink} onClick={() => navigate('/')}>← 回到市場</button>
      </div>
    </div>
  )
}

function Field({ label, type = 'text', value, onChange }) {
  return (
    <div style={styles.field}>
      <label style={styles.label}>{label}</label>
      <input
        style={styles.input}
        type={type}
        value={value}
        onChange={e => onChange(e.target.value)}
        required
      />
    </div>
  )
}

const styles = {
  page: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: '#f0f2f5',
  },
  card: {
    background: '#fff',
    padding: '36px 32px',
    borderRadius: 12,
    boxShadow: '0 4px 24px rgba(0,0,0,.1)',
    width: 380,
  },
  title: {
    fontSize: 22,
    fontWeight: 700,
    color: '#1a1a2e',
    marginBottom: 20,
    textAlign: 'center',
  },
  tabs: {
    display: 'flex',
    borderBottom: '2px solid #f0f0f0',
    marginBottom: 20,
  },
  tab: {
    flex: 1,
    padding: '10px 0',
    background: 'none',
    border: 'none',
    fontSize: 15,
    fontWeight: 600,
    color: '#9ca3af',
    cursor: 'pointer',
  },
  tabActive: {
    color: '#4f46e5',
    borderBottom: '2px solid #4f46e5',
    marginBottom: -2,
  },
  error: {
    background: '#fef2f2',
    color: '#dc2626',
    border: '1px solid #fecaca',
    borderRadius: 8,
    padding: '10px 12px',
    fontSize: 13,
    marginBottom: 14,
  },
  successMsg: {
    background: '#f0fdf4',
    color: '#16a34a',
    border: '1px solid #bbf7d0',
    borderRadius: 8,
    padding: '10px 12px',
    fontSize: 13,
    marginBottom: 14,
  },
  field: { marginBottom: 14 },
  label: {
    display: 'block',
    fontSize: 13,
    fontWeight: 600,
    color: '#444',
    marginBottom: 5,
  },
  input: {
    width: '100%',
    padding: '10px 12px',
    border: '1px solid #ddd',
    borderRadius: 8,
    fontSize: 14,
    outline: 'none',
    boxSizing: 'border-box',
  },
  btn: {
    width: '100%',
    padding: 11,
    background: '#4f46e5',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    fontSize: 15,
    fontWeight: 600,
    cursor: 'pointer',
    marginTop: 6,
  },
  backLink: {
    display: 'block',
    marginTop: 18,
    background: 'none',
    border: 'none',
    color: '#6b7280',
    fontSize: 13,
    cursor: 'pointer',
    textAlign: 'center',
    width: '100%',
  },
  devNotice: {
    textAlign: 'center',
    padding: '28px 16px',
    color: '#6b7280',
  },
  devIcon: { fontSize: 36, marginBottom: 12 },
  devTitle: { fontSize: 16, fontWeight: 700, color: '#374151', marginBottom: 8 },
  devDesc: { fontSize: 13, lineHeight: 1.6 },
}
