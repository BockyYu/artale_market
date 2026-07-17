import { useState, useEffect, useCallback, useRef } from 'react'
import { listItems, listItemCategories, listUsedCategories, createItem, updateItem, updateItemTrack, getItemHistories, togglePriceHistoryHidden, recordItemPrice, exportExcel, sendExcelToDiscord, setItemHidden } from './api'

const EMPTY_FORM = { name: '', english_name: '', search_mode: 1, item_type: 1, category: '', percentage: 0, description: '', track_priority: 0 }


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
  const [categories, setCategories] = useState([])
  const [slotCategories, setSlotCategories] = useState([])
  const [classCategories, setClassCategories] = useState([])
  const [items, setItems] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [filterTypes, setFilterTypes] = useState([])
  const [filterCategories, setFilterCategories] = useState([])
  const [filterPriority, setFilterPriority] = useState(-1)
  const [showSlotDrop, setShowSlotDrop] = useState(false)
  const [showClassDrop, setShowClassDrop] = useState(false)
  const slotDropRef = useRef(null)
  const classDropRef = useRef(null)
  const [updating, setUpdating] = useState(null)
  const [sortBy, setSortBy] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [form, setForm] = useState(EMPTY_FORM)
  const [creating, setCreating] = useState(false)
  const [historyItem, setHistoryItem] = useState(null)
  const [historyRecords, setHistoryRecords] = useState([])
  const [historyLoading, setHistoryLoading] = useState(false)
  const [historyActing, setHistoryActing] = useState(null)
  const [editingItem, setEditingItem] = useState(null)
  const [editForm, setEditForm] = useState({})
  const [savingItem, setSavingItem] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [sendingDc, setSendingDc] = useState(false)
  const [confirmHideItem, setConfirmHideItem] = useState(null)
  const [hiding, setHiding] = useState(false)

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await listItems({ sortBy, search, filterTypes, filterCategories, filterPriority, page, pageSize: PAGE_SIZE })
      setItems(res.data || [])
      setTotal(res.total || 0)
    } catch (err) {
      alert(err.message)
    } finally {
      setLoading(false)
    }
  }, [sortBy, search, filterTypes, filterCategories, filterPriority, page])

  useEffect(() => { load() }, [load])

  useEffect(() => {
    listUsedCategories(6).then(setSlotCategories).catch(() => setSlotCategories([]))
    listUsedCategories(4).then(setClassCategories).catch(() => setClassCategories([]))
  }, [])

  useEffect(() => {
    if (!showSlotDrop) return
    function handler(e) { if (slotDropRef.current && !slotDropRef.current.contains(e.target)) setShowSlotDrop(false) }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [showSlotDrop])

  useEffect(() => {
    if (!showClassDrop) return
    function handler(e) { if (classDropRef.current && !classDropRef.current.contains(e.target)) setShowClassDrop(false) }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [showClassDrop])

  function handleSearchChange(val) { setSearch(val); setPage(1) }
  function toggleFilterType(typeId) {
    setFilterTypes(prev => prev.includes(typeId) ? prev.filter(t => t !== typeId) : [...prev, typeId])
    setPage(1)
  }
  function toggleFilterCategory(cat) {
    setFilterCategories(prev => prev.includes(cat) ? prev.filter(c => c !== cat) : [...prev, cat])
    setPage(1)
  }
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

  async function handleDeleteHistory(id) {
    if (!confirm('確定要刪除這筆記錄？')) return
    setHistoryActing(id)
    try {
      await togglePriceHistoryHidden(id, true)
      setHistoryRecords(prev => prev.filter(r => r.id !== id))
    } catch (err) {
      alert(err.message)
    } finally {
      setHistoryActing(null)
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
      english_name: item.english_name || '',
      search_mode: item.search_mode ?? 1,
      item_type: item.item_type,
      category: item.category,
      percentage: item.percentage,
      description: item.description,
      price: '',
    })
    listItemCategories(item.item_type).then(setCategories).catch(() => setCategories([]))
  }

  async function handleSaveItem(e) {
    e.preventDefault()
    setSavingItem(true)
    try {
      const payload = {
        name: editForm.name,
        english_name: editForm.english_name,
        search_mode: Number(editForm.search_mode),
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

  async function handleConfirmHide() {
    if (!confirmHideItem) return
    setHiding(true)
    try {
      await setItemHidden(confirmHideItem.id, true)
      setItems(prev => prev.filter(i => i.id !== confirmHideItem.id))
      setTotal(prev => prev - 1)
      setConfirmHideItem(null)
    } catch (err) {
      alert(err.message)
    } finally {
      setHiding(false)
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
          <button
            className="btn-add"
            style={{ background: '#16a34a' }}
            disabled={exporting}
            onClick={async () => {
              setExporting(true)
              try { await exportExcel() } catch (err) { alert(err.message) } finally { setExporting(false) }
            }}
          >{exporting ? '產生中...' : '匯出 Excel'}</button>
          <button
            className="btn-add"
            style={{ background: '#5865f2' }}
            disabled={sendingDc}
            onClick={async () => {
              setSendingDc(true)
              try { await sendExcelToDiscord(); alert('已發布到 Discord') } catch (err) { alert(err.message) } finally { setSendingDc(false) }
            }}
          >{sendingDc ? '發布中...' : '發布到 Discord'}</button>
          <button className="btn-add" onClick={() => { setForm(EMPTY_FORM); setShowCreate(true); listItemCategories(EMPTY_FORM.item_type).then(setCategories).catch(() => setCategories([])) }}>+ 新增商品</button>
        </div>
      </div>

      <div className="card">
        <div className="card-toolbar" style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
          <div style={{ position: 'relative', flex: '1 1 200px' }}>
            <input
              className="search-input"
              placeholder="搜尋道具名稱"
              value={search}
              onChange={e => handleSearchChange(e.target.value)}
              style={{ width: '100%', paddingRight: search ? 32 : undefined }}
            />
            {search && (
              <button
                onClick={() => handleSearchChange('')}
                style={{
                  position: 'absolute', right: 10, top: '50%', transform: 'translateY(-50%)',
                  background: 'none', border: 'none', cursor: 'pointer',
                  fontSize: 20, color: '#6b7280', lineHeight: 1, padding: '0 2px',
                }}
              >×</button>
            )}
          </div>

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

        <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center', padding: '8px 0 4px' }}>
          <span style={{ fontSize: 13, color: '#6b7280', whiteSpace: 'nowrap' }}>類型：</span>
          {Object.entries(ITEM_TYPE_LABEL).map(([k, v]) => {
            const typeId = Number(k)
            const active = filterTypes.includes(typeId)
            return (
              <button
                key={k}
                onClick={() => toggleFilterType(typeId)}
                style={{
                  padding: '3px 10px', fontSize: 13, borderRadius: 4, cursor: 'pointer', border: '1px solid',
                  borderColor: active ? '#2563eb' : '#d1d5db',
                  background: active ? '#2563eb' : '#fff',
                  color: active ? '#fff' : '#374151',
                  fontWeight: active ? 600 : 400,
                }}
              >{v}</button>
            )
          })}
          {filterTypes.length > 0 && (
            <button
              onClick={() => { setFilterTypes([]); setPage(1) }}
              style={{ fontSize: 12, color: '#6b7280', background: 'none', border: 'none', cursor: 'pointer', padding: '2px 4px' }}
            >清除</button>
          )}
        </div>

        <div style={{ display: 'flex', gap: 12, alignItems: 'center', padding: '4px 0 8px', flexWrap: 'wrap' }}>
          {[
            { label: '分類', ref: slotDropRef, items: slotCategories, show: showSlotDrop, setShow: setShowSlotDrop },
            { label: '職業', ref: classDropRef, items: classCategories, show: showClassDrop, setShow: setShowClassDrop },
          ].map(({ label, ref, items: opts, show, setShow }) => {
            const activeCount = opts.filter(c => filterCategories.includes(c)).length
            return (
              <div key={label} style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span style={{ fontSize: 13, color: '#6b7280', whiteSpace: 'nowrap' }}>{label}：</span>
                <div ref={ref} style={{ position: 'relative' }}>
                  <button
                    onClick={() => setShow(v => !v)}
                    style={{
                      padding: '4px 10px', fontSize: 13, borderRadius: 4, cursor: 'pointer',
                      border: '1px solid', borderColor: activeCount > 0 ? '#2563eb' : '#d1d5db',
                      background: activeCount > 0 ? '#eff6ff' : '#fff',
                      color: activeCount > 0 ? '#2563eb' : '#374151',
                      minWidth: 90,
                    }}
                  >
                    {activeCount === 0 ? `全部${label}` : `已選 ${activeCount} 項`} ▾
                  </button>
                  {show && (
                    <div style={{
                      position: 'absolute', top: 'calc(100% + 4px)', left: 0, zIndex: 200,
                      background: '#fff', border: '1px solid #e5e7eb', borderRadius: 6,
                      boxShadow: '0 4px 16px rgba(0,0,0,0.1)', padding: '4px 0',
                      maxHeight: 260, overflowY: 'auto', minWidth: 140,
                    }}>
                      {opts.length === 0 && (
                        <div style={{ padding: '8px 12px', fontSize: 13, color: '#9ca3af' }}>無資料</div>
                      )}
                      {opts.map(cat => (
                        <label key={cat} style={{ display: 'flex', alignItems: 'center', gap: 8, padding: '5px 12px', cursor: 'pointer', fontSize: 13, userSelect: 'none' }}>
                          <input type="checkbox" checked={filterCategories.includes(cat)} onChange={() => toggleFilterCategory(cat)} style={{ cursor: 'pointer' }} />
                          {cat}
                        </label>
                      ))}
                    </div>
                  )}
                </div>
                {activeCount > 0 && (
                  <button
                    onClick={() => { setFilterCategories(prev => prev.filter(c => !opts.includes(c))); setPage(1) }}
                    style={{ fontSize: 12, color: '#6b7280', background: 'none', border: 'none', cursor: 'pointer', padding: '2px 4px' }}
                  >清除</button>
                )}
              </div>
            )
          })}
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
              <th className="sortable-th" onClick={handleSortChanges} style={{ cursor: 'pointer', padding: '11px 8px' }}>
                今日修改
                <span className="sort-icon">
                  {sortBy === 'changes_desc' ? ' ▼' : sortBy === 'changes_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th className="sortable-th" onClick={handleSortViews} style={{ cursor: 'pointer', padding: '11px 8px' }}>
                今日查詢
                <span className="sort-icon">
                  {sortBy === 'views_desc' ? ' ▼' : sortBy === 'views_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th>歷史價格</th>
              <th>修改</th>
              <th>查詢優先度</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr className="empty-row"><td colSpan={10}>載入中...</td></tr>
            )}
            {!loading && items.length === 0 && (
              <tr className="empty-row"><td colSpan={10}>無符合資料</td></tr>
            )}
            {items.map((item, index) => (
              <tr key={item.id} className={((page - 1) * PAGE_SIZE + index) % 2 === 0 ? 'row-odd' : 'row-even'}>
                <td>{item.id}</td>
                <td className="text-bold">{item.name}</td>
                <td>{item.category}</td>
                <td>{ITEM_TYPE_LABEL[item.item_type] ?? item.item_type}</td>
                <td>
                  <div style={{ color: item.latest_price != null ? '#16a34a' : '#9ca3af', fontWeight: item.latest_price != null ? 700 : 400 }}>
                    {item.latest_price != null ? item.latest_price.toLocaleString() : '—'}
                  </div>
                  {item.latest_price_at && (() => { const d = new Date(item.latest_price_at); return d.getFullYear() > 2000 ? <div style={{ fontSize: 15, color: '#111827', fontWeight: 400, marginTop: 2 }}>{d.toLocaleString('zh-TW')}</div> : null })()}
                </td>
                <td style={{ color: item.today_changes > 0 ? '#374151' : '#9ca3af', padding: '12px 8px' }}>
                  {item.today_changes} 次
                </td>
                <td style={{ color: item.today_views > 0 ? '#374151' : '#9ca3af', padding: '12px 8px' }}>
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
                <td style={{ textAlign: 'center' }}>
                  <button
                    title="刪除"
                    onClick={() => setConfirmHideItem(item)}
                    style={{ background: 'none', border: 'none', cursor: 'pointer', padding: '4px 6px', color: '#dc2626' }}
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                      <path d="M10 11v6m4-6v6"/><path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/>
                    </svg>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 16, padding: '0 4px' }}>
          <span style={{ fontSize: 16, color: '#6b7280', fontWeight: 600 }}>
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
                <label>英文名稱</label>
                <input value={editForm.english_name} onChange={e => setEditForm(f => ({ ...f, english_name: e.target.value }))} placeholder="選填" />
              </div>
              <div className="form-group">
                <label>查詢方式</label>
                <select className="search-input" style={{ width: '100%', maxWidth: '100%' }}
                  value={editForm.search_mode}
                  onChange={e => setEditForm(f => ({ ...f, search_mode: Number(e.target.value) }))}>
                  <option value={1}>1 - 中文</option>
                  <option value={2}>2 - 英文</option>
                </select>
              </div>
              <div className="form-group">
                <label>類型 *</label>
                <select className="search-input" style={{ width: '100%', maxWidth: '100%' }}
                  value={editForm.item_type}
                  onChange={e => {
                    const t = Number(e.target.value)
                    setEditForm(f => ({ ...f, item_type: t, category: '' }))
                    listItemCategories(t).then(setCategories).catch(() => setCategories([]))
                  }}>
                  {Object.entries(ITEM_TYPE_LABEL).map(([k, v]) => (
                    <option key={k} value={Number(k)}>{v}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>分類 *</label>
                <select className="search-input" style={{ width: '100%', maxWidth: '100%' }}
                  value={editForm.category}
                  onChange={e => setEditForm(f => ({ ...f, category: e.target.value }))}>
                  <option value="">-- 選擇分類 --</option>
                  {categories.map(c => <option key={c} value={c}>{c}</option>)}
                </select>
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
          <div className="modal" style={{ width: '70vw', maxWidth: 950 }} onClick={e => e.stopPropagation()}>
            <h2 style={{ fontSize: 22 }}>{historyItem.name} — 歷史價格</h2>
            {historyLoading ? (
              <p style={{ color: '#6b7280', fontSize: 18 }}>載入中...</p>
            ) : historyRecords.length === 0 ? (
              <p style={{ color: '#9ca3af', fontSize: 18 }}>尚無價格記錄</p>
            ) : (
              <div style={{ maxHeight: 560, overflowY: 'auto', overflowX: 'hidden', marginTop: 8, border: '1px solid #f0f0f0', borderRadius: 8 }}>
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                  <thead>
                    <tr>
                      <th style={{ padding: '12px 16px', fontSize: 18, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>時間</th>
                      <th style={{ padding: '12px 16px', fontSize: 18, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>價格</th>
                      <th style={{ padding: '12px 16px', fontSize: 18, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>來源</th>
                      <th style={{ padding: '12px 16px', fontSize: 18, fontWeight: 700, color: '#374151', background: '#f8f9fb', textAlign: 'left', position: 'sticky', top: 0 }}>操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    {historyRecords.map(r => (
                      <tr key={r.id} style={{ borderTop: '1px solid #f3f4f6' }}>
                        <td style={{ padding: '12px 16px', fontSize: 17, color: '#374151' }}>{new Date(r.recorded_at).toLocaleString('zh-TW')}</td>
                        <td style={{ padding: '12px 16px', fontSize: 19, color: '#16a34a', fontWeight: 600 }}>{r.price.toLocaleString()}</td>
                        <td style={{ padding: '12px 16px', fontSize: 17, color: '#6b7280' }}>{r.source === 'admin' ? '手動' : '自動'}</td>
                        <td style={{ padding: '12px 16px' }}>
                          <button
                            className="btn-action"
                            style={{ color: '#dc2626', opacity: historyActing === r.id ? 0.5 : 1 }}
                            disabled={historyActing === r.id}
                            onClick={() => handleDeleteHistory(r.id)}
                          >
                            刪除
                          </button>
                        </td>
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
                <label>英文名稱</label>
                <input value={form.english_name} onChange={e => setForm(f => ({ ...f, english_name: e.target.value }))} placeholder="選填" />
              </div>
              <div className="form-group">
                <label>查詢方式</label>
                <select className="search-input" style={{ width: '100%', maxWidth: '100%' }}
                  value={form.search_mode}
                  onChange={e => setForm(f => ({ ...f, search_mode: Number(e.target.value) }))}>
                  <option value={1}>1 - 中文</option>
                  <option value={2}>2 - 英文</option>
                </select>
              </div>
              <div className="form-group">
                <label>類型 *</label>
                <select
                  className="search-input"
                  style={{ width: '100%', maxWidth: '100%' }}
                  value={form.item_type}
                  onChange={e => {
                    const t = Number(e.target.value)
                    setForm(f => ({ ...f, item_type: t, category: '' }))
                    listItemCategories(t).then(setCategories).catch(() => setCategories([]))
                  }}
                >
                  {Object.entries(ITEM_TYPE_LABEL).map(([k, v]) => (
                    <option key={k} value={Number(k)}>{v}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>分類 *</label>
                <select className="search-input" style={{ width: '100%', maxWidth: '100%' }}
                  value={form.category}
                  onChange={e => setForm(f => ({ ...f, category: e.target.value }))}>
                  <option value="">-- 選擇分類 --</option>
                  {categories.map(c => <option key={c} value={c}>{c}</option>)}
                </select>
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
      {confirmHideItem && (
        <div className="modal-overlay" onClick={() => setConfirmHideItem(null)}>
          <div className="modal" style={{ maxWidth: 400 }} onClick={e => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, marginBottom: 8 }}>確認刪除</h2>
            <p style={{ color: '#374151', marginBottom: 4 }}>
              確定要刪除 <strong>{confirmHideItem.name}</strong>？
            </p>
            <p style={{ color: '#6b7280', fontSize: 13, marginBottom: 24 }}>
              刪除後此道具將不再出現於前台與後台，且無法復原。
            </p>
            <div className="modal-actions">
              <button className="btn-cancel" onClick={() => setConfirmHideItem(null)} disabled={hiding}>取消</button>
              <button
                className="btn-save"
                style={{ background: '#dc2626' }}
                disabled={hiding}
                onClick={handleConfirmHide}
              >{hiding ? '刪除中...' : '確認刪除'}</button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}
