import { useState, useEffect, useCallback } from 'react'
import { listItems, createItem, updateItemTrack, getItemPrices } from './api'

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

export default function Items() {
  const [items, setItems] = useState([])
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

  const load = useCallback(async (sort) => {
    setLoading(true)
    try {
      const res = await listItems(sort)
      setItems(res || [])
    } catch (err) {
      alert(err.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load(sortBy) }, [load, sortBy])

  function handleSortPrice() {
    setSortBy(s => s === 'price_desc' ? 'price_asc' : 'price_desc')
  }

  async function handleOpenHistory(item) {
    setHistoryItem(item)
    setHistoryRecords([])
    setHistoryLoading(true)
    try {
      const records = await getItemPrices(item.id)
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

  const filtered = items.filter(item => {
    if (filterType !== 0 && item.item_type !== filterType) return false
    if (filterPriority !== -1 && item.track_priority !== filterPriority) return false
    if (search.trim()) {
      const kw = search.trim().toLowerCase()
      if (!item.name.toLowerCase().includes(kw) && !item.category.toLowerCase().includes(kw)) return false
    }
    return true
  })

  return (
    <>
      <div className="page-header">
        <h1>道具列表</h1>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 13, color: '#6b7280' }}>共 {filtered.length} / {items.length} 筆</span>
          <button className="btn-add" onClick={() => { setForm(EMPTY_FORM); setShowCreate(true) }}>+ 新增商品</button>
        </div>
      </div>

      <div className="card">
        <div className="card-toolbar" style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
          <input
            className="search-input"
            placeholder="搜尋名稱 / 分類"
            value={search}
            onChange={e => setSearch(e.target.value)}
            style={{ flex: '1 1 200px' }}
          />

          <select
            className="search-input"
            value={filterType}
            onChange={e => setFilterType(Number(e.target.value))}
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
            onChange={e => setFilterPriority(Number(e.target.value))}
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
              <th>ID</th>
              <th>名稱</th>
              <th>分類</th>
              <th>類型</th>
              <th className="sortable-th" onClick={handleSortPrice} style={{ cursor: 'pointer' }}>
                最新價格
                <span className="sort-icon">
                  {sortBy === 'price_desc' ? ' ▼' : sortBy === 'price_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th>歷史價格</th>
              <th>查詢優先度</th>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr className="empty-row"><td colSpan={7}>載入中...</td></tr>
            )}
            {!loading && filtered.length === 0 && (
              <tr className="empty-row"><td colSpan={7}>無符合資料</td></tr>
            )}
            {filtered.map(item => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td className="text-bold">{item.name}</td>
                <td>{item.category}</td>
                <td>{ITEM_TYPE_LABEL[item.item_type] ?? item.item_type}</td>
                <td className={item.latest_price != null ? 'text-price' : 'text-muted'}>
                  {item.latest_price != null ? item.latest_price.toLocaleString() : '—'}
                </td>
                <td>
                  <button className="btn-action btn-edit" onClick={() => handleOpenHistory(item)}>查看</button>
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
      </div>
      {historyItem && (
        <div className="modal-overlay" onClick={() => setHistoryItem(null)}>
          <div className="modal" style={{ width: 520 }} onClick={e => e.stopPropagation()}>
            <h2>{historyItem.name} — 歷史價格</h2>
            {historyLoading ? (
              <p style={{ color: '#6b7280', fontSize: 14 }}>載入中...</p>
            ) : historyRecords.length === 0 ? (
              <p style={{ color: '#9ca3af', fontSize: 14 }}>尚無價格記錄</p>
            ) : (
              <div style={{ maxHeight: 400, overflowY: 'auto', marginTop: 8, border: '1px solid #f0f0f0', borderRadius: 8 }}>
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                  <thead>
                    <tr>
                      <th style={{ padding: '10px 14px', fontSize: 12, fontWeight: 600, color: '#6b7280', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>日期</th>
                      <th style={{ padding: '10px 14px', fontSize: 12, fontWeight: 600, color: '#6b7280', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>價格</th>
                      <th style={{ padding: '10px 14px', fontSize: 12, fontWeight: 600, color: '#6b7280', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>更新時間</th>
                    </tr>
                  </thead>
                  <tbody>
                    {historyRecords.map(r => (
                      <tr key={r.id} style={{ borderTop: '1px solid #f3f4f6' }}>
                        <td style={{ padding: '10px 14px', fontSize: 14, color: '#374151' }}>{new Date(r.recorded_date).toLocaleDateString('zh-TW')}</td>
                        <td style={{ padding: '10px 14px', fontSize: 14, color: '#16a34a', fontWeight: 600 }}>{r.price.toLocaleString()}</td>
                        <td style={{ padding: '10px 14px', fontSize: 12, color: '#6b7280' }}>{new Date(r.updated_at).toLocaleString('zh-TW')}</td>
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
