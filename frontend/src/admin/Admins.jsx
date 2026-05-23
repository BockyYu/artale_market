import { useState, useEffect, useCallback } from 'react'
import { listAdmins, createAdmin, updateAdmin, deleteAdmin, getPermissions, updatePermissions } from './api'

const EMPTY_FORM = { username: '', password: '', role: 'admin' }

const PERM_LABELS = {
  price_write: '新增每日價格',
  admin_manage: '管理員操作',
}

export default function Admins() {
  const [admins, setAdmins] = useState([])
  const [loading, setLoading] = useState(false)

  // 編輯 modal
  const [modal, setModal] = useState(null) // null | 'create' | { admin }
  const [form, setForm] = useState(EMPTY_FORM)
  const [formError, setFormError] = useState('')
  const [saving, setSaving] = useState(false)

  // 權限 modal
  const [permModal, setPermModal] = useState(null) // null | { admin }
  const [perms, setPerms] = useState({ price_write: false, admin_manage: false })
  const [permSaving, setPermSaving] = useState(false)
  const [permError, setPermError] = useState('')

  const load = useCallback(async () => {
    setLoading(true)
    try { setAdmins(await listAdmins()) }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { load() }, [load])

  // ── 編輯 ──
  function openCreate() {
    setForm(EMPTY_FORM)
    setFormError('')
    setModal('create')
  }

  function openEdit(admin) {
    setForm({ username: admin.username, password: '', role: admin.role })
    setFormError('')
    setModal(admin)
  }

  async function handleSave() {
    setSaving(true)
    setFormError('')
    try {
      if (modal === 'create') {
        await createAdmin(form)
      } else {
        const payload = { role: form.role }
        if (form.username) payload.username = form.username
        if (form.password) payload.password = form.password
        await updateAdmin(modal.id, payload)
      }
      setModal(null)
      load()
    } catch (err) {
      setFormError(err.message)
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id) {
    if (!confirm('確定刪除此管理員？')) return
    try { await deleteAdmin(id); load() }
    catch (err) { alert(err.message) }
  }

  // ── 權限 ──
  async function openPerms(admin) {
    setPermError('')
    setPermModal(admin)
    try {
      const data = await getPermissions(admin.id)
      setPerms(data)
    } catch (err) {
      setPermError(err.message)
    }
  }

  async function handlePermSave() {
    setPermSaving(true)
    setPermError('')
    try {
      const updated = await updatePermissions(permModal.id, perms)
      setPerms(updated)
      setPermModal(null)
    } catch (err) {
      setPermError(err.message)
    } finally {
      setPermSaving(false)
    }
  }

  const isSuperadmin = (admin) => admin.role === 'superadmin'

  return (
    <>
      <div className="page-header">
        <h1>管理員帳號</h1>
        <button className="btn-add" onClick={openCreate}>+ 新增管理員</button>
      </div>

      <div className="card">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>帳號</th>
              <th>角色</th>
              <th>建立時間</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {loading && <tr className="empty-row"><td colSpan={5}>載入中...</td></tr>}
            {!loading && admins.length === 0 && (
              <tr className="empty-row"><td colSpan={5}>目前無管理員資料</td></tr>
            )}
            {admins.map(a => (
              <tr key={a.id}>
                <td>{a.id}</td>
                <td>{a.username}</td>
                <td><span className={`badge badge-${a.role}`}>{a.role}</span></td>
                <td>{new Date(a.created_at).toLocaleDateString('zh-TW')}</td>
                <td>
                  <button className="btn-action btn-edit" onClick={() => openEdit(a)}>編輯</button>
                  <button
                    className="btn-action"
                    style={{ background: '#ede9fe', color: '#7c3aed', marginRight: 6 }}
                    onClick={() => openPerms(a)}
                  >
                    權限
                  </button>
                  <button className="btn-action btn-delete" onClick={() => handleDelete(a.id)}>刪除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* 編輯 Modal */}
      {modal && (
        <div className="modal-overlay" onClick={() => setModal(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>{modal === 'create' ? '新增管理員' : '編輯管理員'}</h2>
            {formError && <div className="error-msg">{formError}</div>}
            <div className="form-group">
              <label>帳號</label>
              <input
                value={form.username}
                placeholder={modal !== 'create' ? '不填則保留原帳號' : '請輸入帳號'}
                onChange={e => setForm(f => ({ ...f, username: e.target.value }))}
              />
              <p className="hint">只允許英文字母與數字</p>
            </div>
            <div className="form-group">
              <label>密碼</label>
              <input
                type="password"
                value={form.password}
                placeholder={modal !== 'create' ? '不填則保留原密碼' : '請輸入密碼（至少6碼）'}
                onChange={e => setForm(f => ({ ...f, password: e.target.value }))}
              />
              <p className="hint">只允許英文字母與數字，至少6碼</p>
            </div>
            <div className="form-group">
              <label>角色</label>
              <select
                value={form.role}
                onChange={e => setForm(f => ({ ...f, role: e.target.value }))}
                style={{ width: '100%', padding: '10px 12px', border: '1px solid #ddd', borderRadius: 8, fontSize: 14 }}
              >
                <option value="admin">admin</option>
                <option value="superadmin">superadmin</option>
              </select>
            </div>
            <div className="modal-actions">
              <button className="btn-cancel" onClick={() => setModal(null)}>取消</button>
              <button className="btn-save" onClick={handleSave} disabled={saving}>
                {saving ? '儲存中...' : '儲存'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 權限 Modal */}
      {permModal && (
        <div className="modal-overlay" onClick={() => setPermModal(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>權限設定 — {permModal.username}</h2>
            {permError && <div className="error-msg">{permError}</div>}

            {isSuperadmin(permModal) ? (
              <p style={{ color: '#6b7280', fontSize: 14, margin: '12px 0 20px' }}>
                superadmin 擁有所有權限，無法個別調整。
              </p>
            ) : (
              <div style={{ margin: '12px 0 20px' }}>
                {Object.entries(PERM_LABELS).map(([key, label]) => (
                  <label
                    key={key}
                    style={{
                      display: 'flex', alignItems: 'center', gap: 10,
                      padding: '10px 0', borderBottom: '1px solid #f3f4f6',
                      cursor: 'pointer', fontSize: 14,
                    }}
                  >
                    <input
                      type="checkbox"
                      checked={perms[key] ?? false}
                      onChange={e => setPerms(p => ({ ...p, [key]: e.target.checked }))}
                      style={{ width: 16, height: 16, cursor: 'pointer' }}
                    />
                    {label}
                  </label>
                ))}
              </div>
            )}

            <div className="modal-actions">
              <button className="btn-cancel" onClick={() => setPermModal(null)}>取消</button>
              {!isSuperadmin(permModal) && (
                <button className="btn-save" onClick={handlePermSave} disabled={permSaving}>
                  {permSaving ? '儲存中...' : '儲存'}
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  )
}
