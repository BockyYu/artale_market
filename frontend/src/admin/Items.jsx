import { useState, useEffect, useCallback } from 'react'
import { listItems, createItem, updateItem, updateItemTrack, getItemHistories, recordItemPrice } from './api'

const EMPTY_FORM = { name: '', item_type: 1, category: '', percentage: 0, description: '', track_priority: 0 }

const ITEM_TYPE_LABEL = {
  1: '卷軸',
  2: '素材',
  3: '消耗品',
  4: '技能書',
  5: '商城',
  6: '裝備',
  7: '活動道具',
}

const TRACK_PRIORITY_LABEL = {
  0: '不追蹤',
  1: '優先',
  2: '次要',
  3: '未出現',
}

const TRACK_PRIORITY_CLASS = {
  0: '',
  1: 'badge-active',
  2: 'badge-pending',
  3: 'badge-banned',
}

const PAGE_SIZE = 20

export default function Items() {
  const [items, setItems] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [filterType, setFilterType] = useState(0)
  const [filterPriority, setFilterPriority] = useState(-1)
  const [updating, setUpdating] = useState(null)
  const [sortBy, setSortBy] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [form, setForm] = useState(EMPTY_FORM)
  const [creating, setCreating] = useState(false)
  const [historyItem, setHistoryItem] = useState(null)
  const [historyRecords, setHistoryRecords] = useState([])
  const [historyLoading, setHistoryLoading] = useState(false)
  const [editingItem, setEditingItem] = useState(null)
  const [editForm, setEditForm] = useState({})
  const [savingItem, setSavingItem] = useState(false)

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await listItems({ sortBy, search, filterType, filterPriority, page, pageSize: PAGE_SIZE })
      setItems(res.data || [])
      setTotal(res.total || 0)
    } catch (err) {
      alert(err.message)
    } finally {
      setLoading(false)
    }
  }, [sortBy, search, filterType, filterPriority, page])

  useEffect(() => { load() }, [load])

  function handleSearchChange(val) { setSearch(val); setPage(1) }
  function handleFilterTypeChange(val) { setFilterType(val); setPage(1) }
  function handleFilterPriorityChange(val) { setFilterPriority(val); setPage(1) }

  function handleSortId() {
    setSortBy(s => s === 'id_desc' ? '' : 'id_desc')
    setPage(1)
  }

  function handleSortPrice() {
    setSortBy(s => s === 'price_desc' ? 'price_asc' : 'price_desc')
    setPage(1)
  }

  function handleSortChanges() {
    setSortBy(s => s === 'changes_desc' ? 'changes_asc' : 'changes_desc')
    setPage(1)
  }

  function handleSortViews() {
    setSortBy(s => s === 'views_desc' ? 'views_asc' : 'views_desc')
    setPage(1)
  }

  async function handleOpenHistory(item) {
    setHistoryItem(item)
    setHistoryRecords([])
    setHistoryLoading(true)
    try {
      const records = await getItemHistories(item.id)
      setHistoryRecords(records || [])
    } catch (err) {
      alert(err.message)
    } finally {
      setHistoryLoading(false)
    }
  }

  async function handleCreate(e) {
    e.preventDefault()
    setCreating(true)
    try {
      const payload = {
        ...form,
        item_type: Number(form.item_type),
        percentage: Number(form.item_type) === 1 ? Number(form.percentage) : 0,
        track_priority: Number(form.track_priority),
      }
      const created = await createItem(payload)
      setItems(prev => [...prev, { ...created, latest_price: null }])
      setShowCreate(false)
      setForm(EMPTY_FORM)
    } catch (err) {
      alert(err.message)
    } finally {
      setCreating(false)
    }
  }

  function handleOpenEdit(item) {
    setEditingItem(item)
    setEditForm({
      name: item.name,
      item_type: item.item_type,
      category: item.category,
      percentage: item.percentage,
      description: item.description,
      price: item.latest_price != null ? Number(item.latest_price).toLocaleString() : '',
    })
  }

  async function handleSaveItem(e) {
    e.preventDefault()
    setSavingItem(true)
    try {
      const payload = {
        name: editForm.name,
        item_type: Number(editForm.item_type),
        percentage: Number(editForm.item_type) === 1 ? Number(editForm.percentage) : 0,
        category: editForm.category,
        description: editForm.description,
      }
      const updated = await updateItem(editingItem.id, payload)
      let priceRecord = null
      if (editForm.price !== '') {
        const price = parseFloat(String(editForm.price).replace(/,/g, ''))
        if (isNaN(price) || price <= 0) {
          alert('請輸入有效價格')
          setSavingItem(false)
          return
        }
        priceRecord = await recordItemPrice(editingItem.id, price)
      }
      setItems(prev => prev.map(i => {
        if (i.id !== editingItem.id) return i
        const next = { ...i, ...updated }
        if (priceRecord) {
          next.latest_price = priceRecord.price
          next.latest_price_at = priceRecord.updated_at || priceRecord.created_at
          next.today_changes = (i.today_changes || 0) + 1
        }
        return next
      }))
      setEditingItem(null)
    } catch (err) {
      alert(err.message)
    } finally {
      setSavingItem(false)
    }
  }

  async function handleTrackChange(item, priority) {
    setUpdating(item.id)
    try {
      const updated = await updateItemTrack(item.id, priority)
      setItems(prev => prev.map(i => i.id === item.id ? { ...i, ...updated } : i))
    } catch (err) {
      alert(err.message)
    } finally {
      setUpdating(null)
    }
  }

  return (
    <>
      <div className="page-header">
        <h1>道具列表</h1>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 13, color: '#374151', fontWeight: 600 }}>共 {total} 筆</span>
          <button className="btn-add" onClick={() => { setForm(EMPTY_FORM); setShowCreate(true) }}>+ 新增商品</button>
        </div>
      </div>

      <div className="card">
        <div className="card-toolbar" style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
          <input
            className="search-input"
            placeholder="搜尋名稱 / 分類"
            value={search}
            onChange={e => handleSearchChange(e.target.value)}
            style={{ flex: '1 1 200px' }}
          />

          <select
            className="search-input"
            value={filterType}
            onChange={e => handleFilterTypeChange(Number(e.target.value))}
            style={{ flex: '0 0 auto' }}
          >
            <option value={0}>全部類型</option>
            {Object.entries(ITEM_TYPE_LABEL).map(([k, v]) => (
              <option key={k} value={Number(k)}>{v}</option>
            ))}
          </select>

          <select
            className="search-input"
            value={filterPriority}
            onChange={e => handleFilterPriorityChange(Number(e.target.value))}
            style={{ flex: '0 0 auto' }}
          >
            <option value={-1}>全部優先度</option>
            {Object.entries(TRACK_PRIORITY_LABEL).map(([k, v]) => (
              <option key={k} value={Number(k)}>{v}</option>
            ))}
          </select>
        </div>

        <table>
          <thead>
            <tr>
              <th className="sortable-th" onClick={handleSortId} style={{ cursor: 'pointer' }}>
                ID
                <span className="sort-icon">
                  {sortBy === '' ? ' ▲' : sortBy === 'id_desc' ? ' ▼' : ' ⇅'}
                </span>
              </th>
              <th>名稱</th>
              <th>分類</th>
              <th>類型</th>
              <th className="sortable-th" onClick={handleSortPrice} style={{ cursor: 'pointer' }}>
                最新價格
                <span className="sort-icon">
                  {sortBy === 'price_desc' ? ' ▼' : sortBy === 'price_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th className="sortable-th" onClick={handleSortChanges} style={{ cursor: 'pointer' }}>
                今日修改
                <span className="sort-icon">
                  {sortBy === 'changes_desc' ? ' ▼' : sortBy === 'changes_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th className="sortable-th" onClick={handleSortViews} style={{ cursor: 'pointer' }}>
                今日查詢
                <span className="sort-icon">
                  {sortBy === 'views_desc' ? ' ▼' : sortBy === 'views_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th>歷史價格</th>
              <th>修改</th>
              <th>查詢優先度</th>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr className="empty-row"><td colSpan={10}>載入中...</td></tr>
            )}
            {!loading && items.length === 0 && (
              <tr className="empty-row"><td colSpan={10}>無符合資料</td></tr>
            )}
            {items.map(item => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td className="text-bold">{item.name}</td>
                <td>{item.category}</td>
                <td>{ITEM_TYPE_LABEL[item.item_type] ?? item.item_type}</td>
                <td>
                  <div style={{ color: item.latest_price != null ? '#16a34a' : '#9ca3af', fontWeight: item.latest_price != null ? 700 : 400 }}>
                    {item.latest_price != null ? item.latest_price.toLocaleString() : '—'}
                  </div>
                  {item.latest_price_at && (() => { const d = new Date(item.latest_price_at); return d.getFullYear() > 2000 ? <div style={{ fontSize: 11, color: '#9ca3af', fontWeight: 400, marginTop: 2 }}>{d.toLocaleString('zh-TW')}</div> : null })()}
                </td>
                <td style={{ color: item.today_changes > 0 ? '#374151' : '#9ca3af' }}>
                  {item.today_changes} 次
                </td>
                <td style={{ color: item.today_views > 0 ? '#374151' : '#9ca3af' }}>
                  {item.today_views} 次
                </td>
                <td>
                  <button className="btn-action btn-edit" onClick={() => handleOpenHistory(item)}>查看</button>
                </td>
                <td>
                  <button className="btn-action btn-edit" onClick={() => handleOpenEdit(item)}>修改</button>
                </td>
                <td>
                  {item.track_priority === 3 ? (
                    <span className={`badge ${TRACK_PRIORITY_CLASS[3]}`}>未出現</span>
                  ) : (
                    <div style={{ display: 'flex', gap: 4 }}>
                      {[0, 1, 2].map(p => (
                        <button
                          key={p}
                          disabled={updating === item.id}
                          className={`btn-action ${item.track_priority === p ? 'btn-active-priority' : ''}`}
                          style={{
                            opacity: updating === item.id ? 0.5 : 1,
                            fontWeight: item.track_priority === p ? 700 : 400,
                            background: item.track_priority === p
                              ? (p === 1 ? '#16a34a' : p === 2 ? '#ca8a04' : '#6b7280')
                              : undefined,
                            color: item.track_priority === p ? '#fff' : undefined,
                            border: item.track_priority === p ? 'none' : undefined,
                          }}
                          onClick={() => handleTrackChange(item, p)}
                        >
                          {TRACK_PRIORITY_LABEL[p]}
                        </button>
                      ))}
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 16, padding: '0 4px' }}>
          <span style={{ fontSize: 13, color: '#6b7280' }}>
            第 {Math.min((page - 1) * PAGE_SIZE + 1, total)} – {Math.min(page * PAGE_SIZE, total)} 筆，共 {total} 筆
          </span>
          <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
            <button className="btn-action" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>上一頁</button>
            <span style={{ fontSize: 13, color: '#374151', minWidth: 80, textAlign: 'center' }}>
              第 {page} / {totalPages} 頁
            </span>
            <button className="btn-action" disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>下一頁</button>
          </div>
        </div>
      </div>
      {editingItem && (
        <div className="modal-overlay" onClick={() => setEditingItem(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>修改道具 — {editingItem.name}</h2>
            <form onSubmit={handleSaveItem}>
              <div className="form-group">
                <label>名稱 *</label>
                <input required value={editForm.name} onChange={e => setEditForm(f => ({ ...f, name: e.target.value }))} />
              </div>
              <div className="form-group">
                <label>類型 *</label>
                <select className="search-input" style={{ width: '100%', maxWidth: '100%' }}
                  value={editForm.item_type}
                  onChange={e => setEditForm(f => ({ ...f, item_type: Number(e.target.value) }))}>
                  {Object.entries(ITEM_TYPE_LABEL).map(([k, v]) => (
                    <option key={k} value={Number(k)}>{v}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>分類 *</label>
                <input required value={editForm.category} onChange={e => setEditForm(f => ({ ...f, category: e.target.value }))} />
              </div>
              {Number(editForm.item_type) === 1 && (
                <div className="form-group">
                  <label>成功率 (%)</label>
                  <input type="number" min={1} max={100} value={editForm.percentage}
                    onChange={e => setEditForm(f => ({ ...f, percentage: e.target.value }))} />
                </div>
              )}
              <div className="form-group">
                <label>備註</label>
                <input value={editForm.description} onChange={e => setEditForm(f => ({ ...f, description: e.target.value }))} placeholder="選填" />
              </div>
              <div className="form-group">
                <label>今日價格</label>
                <input
                  type="text"
                  inputMode="numeric"
                  value={editForm.price}
                  onChange={e => {
                    const raw = e.target.value.replace(/,/g, '').replace(/[^0-9]/g, '')
                    setEditForm(f => ({ ...f, price: raw === '' ? '' : Number(raw).toLocaleString() }))
                  }}
                  placeholder="留空則不修改"
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setEditingItem(null)}>取消</button>
                <button type="submit" className="btn-save" disabled={savingItem}>{savingItem ? '儲存中...' : '儲存'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {historyItem && (
        <div className="modal-overlay" onClick={() => setHistoryItem(null)}>
          <div className="modal" style={{ width: '50vw', maxWidth: 675 }} onClick={e => e.stopPropagation()}>
            <h2>{historyItem.name} — 歷史價格</h2>
            {historyLoading ? (
              <p style={{ color: '#6b7280', fontSize: 14 }}>載入中...</p>
            ) : historyRecords.length === 0 ? (
              <p style={{ color: '#9ca3af', fontSize: 14 }}>尚無價格記錄</p>
            ) : (
              <div style={{ maxHeight: 400, overflowY: 'auto', overflowX: 'hidden', marginTop: 8, border: '1px solid #f0f0f0', borderRadius: 8 }}>
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                  <thead>
                    <tr>
                      <th style={{ padding: '10px 14px', fontSize: 12, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>時間</th>
                      <th style={{ padding: '10px 14px', fontSize: 12, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>價格</th>
                      <th style={{ padding: '10px 14px', fontSize: 12, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>來源</th>
                    </tr>
                  </thead>
                  <tbody>
                    {historyRecords.map(r => (
                      <tr key={r.id} style={{ borderTop: '1px solid #f3f4f6' }}>
                        <td style={{ padding: '10px 14px', fontSize: 13, color: '#374151' }}>{new Date(r.recorded_at).toLocaleString('zh-TW')}</td>
                        <td style={{ padding: '10px 14px', fontSize: 14, color: '#16a34a', fontWeight: 600 }}>{r.price.toLocaleString()}</td>
                        <td style={{ padding: '10px 14px', fontSize: 12, color: '#6b7280' }}>{r.source === 'admin' ? '手動' : '自動'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            <div className="modal-actions">
              <button className="btn-cancel" onClick={() => setHistoryItem(null)}>關閉</button>
            </div>
          </div>
        </div>
      )}

      {showCreate && (
        <div className="modal-overlay" onClick={() => setShowCreate(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>新增商品</h2>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>名稱 *</label>
                <input required value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))} placeholder="商品名稱" />
              </div>
              <div className="form-group">
                <label>類型 *</label>
                <select
                  className="search-input"
                  style={{ width: '100%', maxWidth: '100%' }}
                  value={form.item_type}
                  onChange={e => setForm(f => ({ ...f, item_type: Number(e.target.value) }))}
                >
                  {Object.entries(ITEM_TYPE_LABEL).map(([k, v]) => (
                    <option key={k} value={Number(k)}>{v}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>分類 *</label>
                <input required value={form.category} onChange={e => setForm(f => ({ ...f, category: e.target.value }))} placeholder="例：頭盔、劍士、消耗品..." />
              </div>
              {Number(form.item_type) === 1 && (
                <div className="form-group">
                  <label>成功率 (%)</label>
                  <input type="number" min={1} max={100} value={form.percentage} onChange={e => setForm(f => ({ ...f, percentage: e.target.value }))} placeholder="10 / 30 / 60 / 100" />
                </div>
              )}
              <div className="form-group">
                <label>查詢優先度</label>
                <select
                  className="search-input"
                  style={{ width: '100%', maxWidth: '100%' }}
                  value={form.track_priority}
                  onChange={e => setForm(f => ({ ...f, track_priority: Number(e.target.value) }))}
                >
                  {[0, 1, 2].map(p => (
                    <option key={p} value={p}>{TRACK_PRIORITY_LABEL[p]}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>備註</label>
                <input value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))} placeholder="選填" />
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
