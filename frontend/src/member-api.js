const API_BASE = import.meta.env.VITE_API_BASE || ''
const BASE = `${API_BASE}/api/v1/member`

export async function fetchAppConfig() {
  try {
    const res = await fetch(`${API_BASE}/api/v1/system`)
    if (!res.ok) return { maintenance: true, mode: 'prod', message: '' }
    const data = await res.json()
    return { maintenance: data.enabled, mode: data.mode || 'test', message: data.message || '' }
  } catch {
    return { maintenance: true, mode: 'prod', message: '無法連線到伺服器' }
  }
}

export function memberFetch(url, options = {}) {
  const fullUrl = url.startsWith('http') ? url : `${API_BASE}${url}`
  const token = localStorage.getItem('member_token')
  const headers = { ...(options.headers || {}) }
  if (token) headers['Authorization'] = `Bearer ${token}`
  return fetch(fullUrl, { ...options, headers })
}

export async function memberLogin(username, password) {
  const res = await fetch(`${BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || '登入失敗')
  localStorage.setItem('member_token', data.token)
  const info = { id: data.id, nickname: data.nickname, username: data.username, email: data.email, status: data.status }
  localStorage.setItem('member_info', JSON.stringify(info))
  return info
}

export async function memberRegister({ nickname, username, password, email, invite_code }) {
  const res = await fetch(`${BASE}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname, username, password, email, invite_code }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || '註冊失敗')
  return data
}

export async function memberLogout() {
  await memberFetch(`${BASE}/logout`, { method: 'POST' })
  localStorage.removeItem('member_token')
  localStorage.removeItem('member_info')
}

export async function fetchMe() {
  const info = getMemberInfo()
  if (!info) return null
  const res = await memberFetch(`${BASE}/me`)
  if (!res.ok) {
    localStorage.removeItem('member_info')
    return null
  }
  return res.json()
}

export function getMemberInfo() {
  try {
    return JSON.parse(localStorage.getItem('member_info'))
  } catch {
    return null
  }
}

export async function fetchPriceHistory(itemId, days = 7) {
  const res = await memberFetch(`${BASE}/items/${itemId}/price-history?days=${days}`)
  const data = await res.json()
  return data?.data || []
}
