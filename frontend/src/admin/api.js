const API_BASE = import.meta.env.VITE_API_BASE || ''
const BASE = `${API_BASE}/api/v1/admin`

function getToken() {
  return localStorage.getItem('admin_token')
}

function authHeaders() {
  return {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${getToken()}`,
  }
}

async function handleResponse(res) {
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`)
  return data
}

export async function login(username, password) {
  const res = await fetch(`${BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  const data = await handleResponse(res)
  localStorage.setItem('admin_token', data.token)
  localStorage.setItem('admin_user', JSON.stringify(data.admin))
  return data
}

export function logout() {
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_user')
}

export function currentUser() {
  try {
    return JSON.parse(localStorage.getItem('admin_user'))
  } catch {
    return null
  }
}

// Admins
export async function listAdmins() {
  const res = await fetch(`${BASE}/admins`, { headers: authHeaders() })
  return handleResponse(res)
}

export async function createAdmin(data) {
  const res = await fetch(`${BASE}/admins`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify(data),
  })
  return handleResponse(res)
}

export async function updateAdmin(id, data) {
  const res = await fetch(`${BASE}/admins/${id}`, {
    method: 'PUT',
    headers: authHeaders(),
    body: JSON.stringify(data),
  })
  return handleResponse(res)
}

export async function deleteAdmin(id) {
  const res = await fetch(`${BASE}/admins/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return handleResponse(res)
}

// Permissions
export async function getPermissions(adminId) {
  const res = await fetch(`${BASE}/admins/${adminId}/permissions`, { headers: authHeaders() })
  return handleResponse(res)
}

export async function updatePermissions(adminId, perms) {
  const res = await fetch(`${BASE}/admins/${adminId}/permissions`, {
    method: 'PUT',
    headers: authHeaders(),
    body: JSON.stringify(perms),
  })
  return handleResponse(res)
}

// Items
export async function createItem(data) {
  const res = await fetch(`${BASE}/items`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify(data),
  })
  return handleResponse(res)
}

export async function listItems({ sortBy = '', search = '', filterType = 0, filterPriority = -1, page = 1, pageSize = 20 } = {}) {
  const params = new URLSearchParams({ page, page_size: pageSize })
  if (sortBy) params.set('sort_by', sortBy)
  if (search) params.set('search', search)
  if (filterType > 0) params.set('filter_type', filterType)
  if (filterPriority >= 0) params.set('filter_priority', filterPriority)
  const res = await fetch(`${BASE}/items?${params}`, { headers: authHeaders() })
  return handleResponse(res)
}

export async function getItemPrices(id) {
  const res = await fetch(`${BASE}/items/${id}/prices`, { headers: authHeaders() })
  return handleResponse(res)
}

export async function getItemHistories(id) {
  const res = await fetch(`${BASE}/items/${id}/histories`, { headers: authHeaders() })
  return handleResponse(res)
}

export async function updateItem(id, data) {
  const res = await fetch(`${BASE}/items/${id}`, {
    method: 'PUT',
    headers: authHeaders(),
    body: JSON.stringify(data),
  })
  return handleResponse(res)
}

export async function recordItemPrice(id, price) {
  const res = await fetch(`${BASE}/items/${id}/prices`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ price }),
  })
  return handleResponse(res)
}

export async function updateItemTrack(id, trackPriority) {
  const res = await fetch(`${BASE}/items/${id}/track`, {
    method: 'PATCH',
    headers: authHeaders(),
    body: JSON.stringify({ track_priority: trackPriority }),
  })
  return handleResponse(res)
}

// Members
export async function listMembers(page = 1, pageSize = 20, search = '') {
  const params = new URLSearchParams({ page, page_size: pageSize, search })
  const res = await fetch(`${BASE}/members?${params}`, { headers: authHeaders() })
  return handleResponse(res)
}

export async function updateMemberStatus(id, status) {
  const res = await fetch(`${BASE}/members/${id}/status`, {
    method: 'PUT',
    headers: authHeaders(),
    body: JSON.stringify({ status }),
  })
  return handleResponse(res)
}

export async function deleteMember(id) {
  const res = await fetch(`${BASE}/members/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return handleResponse(res)
}
