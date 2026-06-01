import { useState, useEffect, useCallback, useRef } from 'react'
import { listAlerts, createAlert, updateAlert, deleteAlert, toggleAlertActive, listItems, listBots } from './api'

const PLATFORM_LABEL = { tg: 'Telegram', line: 'LINE Notify', dc: 'Discord' }
const PLATFORM_COLOR = { tg: '#2ca5e0', line: '#06c755', dc: '#5865f2' }

const EMPTY_FORM = { itemID: 0, itemName: '', thresholdPrice: '', botID: '', note: '' }

export default function PriceAlerts() {
  const [alerts, setAlerts] = useState([])
  const [bots, setBots] = useState([])
  const [loading, setLoading] = useState(false)
  const [showCreate, setShowCreate] = useState(false)
  const [form, setForm] = useState(EMPTY_FORM)
  const [creating, setCreating] = useState(false)
  const [toggling, setToggling] = useState(null)
  const [deleting, setDeleting] = useState(null)
  const [editingAlert, setEditingAlert] = useState(null)
  const [editForm, setEditForm] = useState({ thresholdPrice: '', botID: '', note: '' })
  const [saving, setSaving] = useState(false)

  // 道具搜尋 autocomplete
  const [itemSearch, setItemSearch] = useState('')
  const [itemSuggestions, setItemSuggestions] = useState([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const searchRef = useRef(null)
  const searchTimer = useRef(null)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const [alertData, botData] = await Promise.all([listAlerts(), listBots()])
      setAlerts(alertData?.data || [])
      setBots((botData?.data || []).filter(b => b.is_active))
    } catch (err) {
      alert(err.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  useEffect(() => {
    const handler = (e) => {
      if (searchRef.current && !searchRef.current.contains(e.target)) {
        setShowSuggestions(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  function handleItemSearchChange(val) {
    setItemSearch(val)
    setForm(f => ({ ...f, itemID: 0, itemName: '' }))
    clearTimeout(searchTimer.current)
    if (!val.trim()) {
      setItemSuggestions([])
      setShowSuggestions(false)
      return
    }
    searchTimer.current = setTimeout(async () => {
      try {
        const res = await listItems({ search: val, page: 1, pageSize: 10 })
        setItemSuggestions(res.data || [])
        setShowSuggestions(true)
      } catch {
        setItemSuggestions([])
      }
    }, 300)
  }

  function handleSelectItem(item) {
    setForm(f => ({ ...f, itemID: item.id, itemName: item.name }))
    setItemSearch(item.name)
    setItemSuggestions([])
    setShowSuggestions(false)
  }

  async function handleCreate(e) {
    e.preventDefault()
    if (!form.itemID) {
      alert('請從下拉選單選擇道具')
      return
    }
    const price = parseFloat(String(form.thresholdPrice).replace(/,/g, ''))
    if (isNaN(price) || price <= 0) {
      alert('請輸入有效的觸發價格')
      return
    }
    setCreating(true)
    try {
      await createAlert({
        item_id: form.itemID,
        threshold_price: price,
        bot_id: form.botID ? Number(form.botID) : undefined,
        note: form.note,
      })
      await load()
      setShowCreate(false)
      setForm(EMPTY_FORM)
      setItemSearch('')
    } catch (err) {
      alert(err.message)
    } finally {
      setCreating(false)
    }
  }

  async function handleToggle(alert) {
    setToggling(alert.id)
    try {
      await toggleAlertActive(alert.id, !alert.is_active)
      setAlerts(prev => prev.map(a => a.id === alert.id ? { ...a, is_active: !a.is_active } : a))
    } catch (err) {
      alert(err.message)
    } finally {
      setToggling(null)
    }
  }

  function handleOpenEdit(a) {
    setEditingAlert(a)
    setEditForm({
      thresholdPrice: Number(a.threshold_price).toLocaleString(),
      botID: a.bot_id ?? '',
      note: a.note ?? '',
    })
  }

  async function handleUpdate(e) {
    e.preventDefault()
    const price = parseFloat(String(editForm.thresholdPrice).replace(/,/g, ''))
    if (isNaN(price) || price <= 0) {
      alert('請輸入有效的觸發價格')
      return
    }
    setSaving(true)
    try {
      await updateAlert(editingAlert.id, {
        threshold_price: price,
        bot_id: editForm.botID ? Number(editForm.botID) : undefined,
        note: editForm.note,
      })
      setAlerts(prev => prev.map(a => a.id === editingAlert.id
        ? { ...a, threshold_price: price, bot_id: editForm.botID ? Number(editForm.botID) : null, note: editForm.note, bot: editForm.botID ? bots.find(b => b.id === Number(editForm.botID)) ?? a.bot : null }
        : a
      ))
      setEditingAlert(null)
    } catch (err) {
      alert(err.message)
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id) {
    if (!confirm('確定要刪除這筆提醒？')) return
    setDeleting(id)
    try {
      await deleteAlert(id)
      setAlerts(prev => prev.filter(a => a.id !== id))
    } catch (err) {
      alert(err.message)
    } finally {
      setDeleting(null)
    }
  }

  function fmtPrice(p) {
    if (p == null) return '—'
    return Number(p).toLocaleString()
  }

  function fmtTime(t) {
    if (!t) return '—'
    return new Date(t).toLocaleString('zh-TW')
  }

  return (
    <>
      <div className="page-header">
        <h1>價格提醒</h1>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 13, color: '#374151', fontWeight: 600 }}>共 {alerts.length} 筆</span>
          <button className="btn-add" onClick={() => { setForm(EMPTY_FORM); setItemSearch(''); setShowCreate(true) }}>
            + 新增提醒
          </button>
        </div>
      </div>

      <div className="card">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>道具名稱</th>
              <th>觸發價格（低於）</th>
              <th>通知機器人</th>
              <th>備註</th>
              <th>上次觸發</th>
              <th>狀態</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr className="empty-row"><td colSpan={8}>載入中...</td></tr>
            )}
            {!loading && alerts.length === 0 && (
              <tr className="empty-row"><td colSpan={8}>尚無提醒設定</td></tr>
            )}
            {alerts.map(a => (
              <tr key={a.id}>
                <td>{a.id}</td>
                <td className="text-bold">{a.item?.name ?? `#${a.item_id}`}</td>
                <td style={{ color: '#dc2626', fontWeight: 700 }}>{fmtPrice(a.threshold_price)}</td>
                <td>
                  {a.bot ? (
                    <span style={{
                      display: 'inline-block', padding: '2px 10px', borderRadius: 12, fontSize: 12, fontWeight: 700,
                      background: (PLATFORM_COLOR[a.bot.platform] ?? '#888') + '22',
                      color: PLATFORM_COLOR[a.bot.platform] ?? '#888',
                    }}>
                      {a.bot.name}
                    </span>
                  ) : (
                    <span style={{ color: '#9ca3af', fontSize: 12 }}>環境變數 TG</span>
                  )}
                </td>
                <td style={{ color: '#6b7280', fontSize: 13 }}>{a.note || '—'}</td>
                <td style={{ fontSize: 12, color: '#9ca3af' }}>{fmtTime(a.last_triggered_at)}</td>
                <td>
                  <button
                    disabled={toggling === a.id}
                    className={`btn-action ${a.is_active ? 'btn-active-priority' : ''}`}
                    style={{
                      background: a.is_active ? '#16a34a' : undefined,
                      color: a.is_active ? '#fff' : undefined,
                      border: a.is_active ? 'none' : undefined,
                      opacity: toggling === a.id ? 0.5 : 1,
                    }}
                    onClick={() => handleToggle(a)}
                  >
                    {a.is_active ? '啟用中' : '已停用'}
                  </button>
                </td>
                <td style={{ whiteSpace: 'nowrap' }}>
                  <button
                    className="btn-action"
                    style={{ marginRight: 8 }}
                    onClick={() => handleOpenEdit(a)}
                  >
                    編輯
                  </button>
                  <button
                    className="btn-action"
                    style={{ color: '#dc2626', opacity: deleting === a.id ? 0.5 : 1 }}
                    disabled={deleting === a.id}
                    onClick={() => handleDelete(a.id)}
                  >
                    刪除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {editingAlert && (
        <div className="modal-overlay" onClick={() => setEditingAlert(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>編輯價格提醒</h2>
            <p style={{ color: '#6b7280', marginBottom: 16, fontSize: 14 }}>
              道具：<strong style={{ color: '#374151' }}>{editingAlert.item?.name ?? `#${editingAlert.item_id}`}</strong>
            </p>
            <form onSubmit={handleUpdate}>
              <div className="form-group">
                <label>觸發價格（低於此價格時通知）*</label>
                <input
                  type="text"
                  inputMode="numeric"
                  placeholder="例：1,000,000"
                  value={editForm.thresholdPrice}
                  onChange={e => {
                    const raw = e.target.value.replace(/,/g, '').replace(/[^0-9]/g, '')
                    setEditForm(f => ({ ...f, thresholdPrice: raw === '' ? '' : Number(raw).toLocaleString() }))
                  }}
                />
              </div>

              <div className="form-group">
                <label>通知機器人</label>
                <select
                  className="search-input"
                  style={{ width: '100%', maxWidth: '100%' }}
                  value={editForm.botID}
                  onChange={e => setEditForm(f => ({ ...f, botID: e.target.value }))}
                >
                  <option value="">環境變數 TG（預設）</option>
                  {bots.map(b => (
                    <option key={b.id} value={b.id}>
                      [{PLATFORM_LABEL[b.platform] ?? b.platform}] {b.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>備註（選填）</label>
                <input
                  value={editForm.note}
                  onChange={e => setEditForm(f => ({ ...f, note: e.target.value }))}
                  placeholder="例：想買的價位"
                />
              </div>

              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setEditingAlert(null)}>取消</button>
                <button type="submit" className="btn-save" disabled={saving}>{saving ? '儲存中...' : '儲存'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showCreate && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>新增價格提醒</h2>
            <form onSubmit={handleCreate}>
              <div className="form-group" ref={searchRef} style={{ position: 'relative' }}>
                <label>道具名稱 *</label>
                <input
                  placeholder="輸入名稱搜尋..."
                  value={itemSearch}
                  onChange={e => handleItemSearchChange(e.target.value)}
                  onFocus={() => itemSuggestions.length > 0 && setShowSuggestions(true)}
                  autoComplete="off"
                />
                {showSuggestions && itemSuggestions.length > 0 && (
                  <ul style={{
                    position: 'absolute', top: '100%', left: 0, right: 0,
                    background: '#fff', border: '1px solid #e5e7eb', borderRadius: 6,
                    boxShadow: '0 4px 12px rgba(0,0,0,0.1)', listStyle: 'none',
                    margin: 0, padding: 0, zIndex: 100, maxHeight: 220, overflowY: 'auto',
                  }}>
                    {itemSuggestions.map(item => (
                      <li
                        key={item.id}
                        onMouseDown={e => { e.preventDefault(); handleSelectItem(item) }}
                        style={{
                          padding: '8px 12px', cursor: 'pointer', fontSize: 13,
                          borderBottom: '1px solid #f3f4f6',
                        }}
                        onMouseEnter={e => e.currentTarget.style.background = '#f9fafb'}
                        onMouseLeave={e => e.currentTarget.style.background = ''}
                      >
                        <span style={{ fontWeight: 600 }}>{item.name}</span>
                        <span style={{ color: '#9ca3af', marginLeft: 8, fontSize: 12 }}>{item.category}</span>
                      </li>
                    ))}
                  </ul>
                )}
                {form.itemID > 0 && (
                  <div style={{ marginTop: 4, fontSize: 12, color: '#16a34a' }}>
                    已選取：{form.itemName} (ID: {form.itemID})
                  </div>
                )}
              </div>

              <div className="form-group">
                <label>觸發價格（低於此價格時通知）*</label>
                <input
                  type="text"
                  inputMode="numeric"
                  placeholder="例：1,000,000"
                  value={form.thresholdPrice}
                  onChange={e => {
                    const raw = e.target.value.replace(/,/g, '').replace(/[^0-9]/g, '')
                    setForm(f => ({ ...f, thresholdPrice: raw === '' ? '' : Number(raw).toLocaleString() }))
                  }}
                />
              </div>

              <div className="form-group">
                <label>通知機器人</label>
                <select
                  className="search-input"
                  style={{ width: '100%', maxWidth: '100%' }}
                  value={form.botID}
                  onChange={e => setForm(f => ({ ...f, botID: e.target.value }))}
                >
                  <option value="">環境變數 TG（預設）</option>
                  {bots.map(b => (
                    <option key={b.id} value={b.id}>
                      [{PLATFORM_LABEL[b.platform] ?? b.platform}] {b.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>備註（選填）</label>
                <input
                  value={form.note}
                  onChange={e => setForm(f => ({ ...f, note: e.target.value }))}
                  placeholder="例：想買的價位"
                />
              </div>

              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setShowCreate(false)}>取消</button>
                <button type="submit" className="btn-save" disabled={creating}>{creating ? '建立中...' : '建立'}</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  )
}
