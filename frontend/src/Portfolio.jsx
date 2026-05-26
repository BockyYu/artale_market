import { useState, useEffect, useRef } from 'react'
import { memberFetch } from './member-api'

const STORAGE_KEY = 'artale_portfolio'

function loadRecords() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY)) || [] }
  catch { return [] }
}
function saveRecords(records) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(records))
}

function formatPrice(val) {
  const raw = String(val).replace(/,/g, '').replace(/[^0-9]/g, '')
  return raw === '' ? '' : Number(raw).toLocaleString()
}

const fmt = n => n != null ? Math.round(n).toLocaleString() : '—'

function DiffCell({ value, pct }) {
  if (value == null) return <span style={{ color: '#9ca3af' }}>—</span>
  const up = value >= 0
  return (
    <span className={up ? 'change-up' : 'change-down'}>
      {up ? '+' : ''}{Math.round(value).toLocaleString()}
      {pct != null && (
        <span style={{ fontSize: '0.78rem', marginLeft: 4 }}>
          ({up ? '+' : ''}{pct.toFixed(1)}%)
        </span>
      )}
    </span>
  )
}

export default function Portfolio() {
  const [records, setRecords] = useState(loadRecords)
  const [prices, setPrices] = useState({})
  const [loadingPrices, setLoadingPrices] = useState(false)
  const [allScrolls, setAllScrolls] = useState([])
  const [showAdd, setShowAdd] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [showSugg, setShowSugg] = useState(false)
  const [selectedItem, setSelectedItem] = useState(null)
  const [qty, setQty] = useState('')
  const [buyPrice, setBuyPrice] = useState('')
  const searchRef = useRef(null)

  useEffect(() => {
    memberFetch('/api/v1/member/items')
      .then(r => r.json())
      .then(data => setAllScrolls((data || []).filter(i => i.item_type === 1)))
      .catch(() => {})
  }, [])

  useEffect(() => {
    const ids = [...new Set(records.map(r => r.item_id))]
    if (!ids.length) { setPrices({}); return }
    setLoadingPrices(true)
    Promise.all(
      ids.map(id =>
        memberFetch(`/api/v1/member/items/${id}/prices`)
          .then(r => r.json())
          .then(d => [id, d.today_price ?? null])
          .catch(() => [id, null])
      )
    ).then(entries => {
      setPrices(Object.fromEntries(entries))
      setLoadingPrices(false)
    })
  }, [records])

  useEffect(() => {
    const handler = e => {
      if (searchRef.current && !searchRef.current.contains(e.target)) setShowSugg(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  function handleAdd(e) {
    e.preventDefault()
    if (!selectedItem) { alert('請選擇卷軸'); return }
    const quantity = parseInt(qty)
    const price = parseFloat(String(buyPrice).replace(/,/g, ''))
    if (!quantity || quantity <= 0) { alert('請輸入張數'); return }
    if (!price || price <= 0) { alert('請輸入買入單價'); return }
    const next = [...records, {
      id: Date.now(),
      item_id: selectedItem.id,
      item_name: selectedItem.name,
      percentage: selectedItem.percentage,
      category: selectedItem.category,
      quantity,
      bought_price: price,
      bought_at: new Date().toISOString(),
    }]
    setRecords(next); saveRecords(next)
    setShowAdd(false)
    setSelectedItem(null); setSearchText(''); setQty(''); setBuyPrice('')
  }

  function handleDelete(id) {
    if (!window.confirm('確定刪除此筆記錄？')) return
    const next = records.filter(r => r.id !== id)
    setRecords(next); saveRecords(next)
  }

  const suggestions = searchText.trim()
    ? allScrolls.filter(i => i.name.toLowerCase().includes(searchText.trim().toLowerCase())).slice(0, 8)
    : []

  // Totals
  let totalCost = 0, totalValue = 0, allHavePrice = records.length > 0
  records.forEach(r => {
    totalCost += r.bought_price * r.quantity
    const cur = prices[r.item_id]
    if (cur != null) totalValue += cur * r.quantity
    else allHavePrice = false
  })
  const totalProfit = allHavePrice ? totalValue - totalCost : null
  const totalPct = totalProfit != null && totalCost > 0 ? (totalProfit / totalCost) * 100 : null

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <div>
          <h2 style={{ fontSize: '1.2rem', fontWeight: 700 }}>卷軸持倉紀錄</h2>
          <p style={{ fontSize: '0.82rem', color: '#9ca3af', marginTop: 2 }}>資料儲存於本機瀏覽器，清除快取會消失</p>
        </div>
        <button
          className="btn btn-primary"
          onClick={() => setShowAdd(true)}
          style={{ padding: '8px 18px', fontSize: '0.88rem', fontWeight: 600 }}
        >
          + 新增購買
        </button>
      </div>

      {records.length === 0 ? (
        <div className="table-wrapper" style={{ textAlign: 'center', padding: '48px 16px', color: '#9ca3af', fontSize: '0.92rem' }}>
          尚無記錄，點「+ 新增購買」開始追蹤
        </div>
      ) : (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>卷軸名稱</th>
                <th>分類</th>
                <th>成功率</th>
                <th>張數</th>
                <th>買入單價</th>
                <th>買入總額</th>
                <th>{loadingPrices ? '現在單價…' : '現在單價'}</th>
                <th>現在總值</th>
                <th>損益</th>
                <th>買入時間</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {records.map(r => {
                const cur = prices[r.item_id]
                const cost = r.bought_price * r.quantity
                const value = cur != null ? cur * r.quantity : null
                const profit = value != null ? value - cost : null
                const pct = profit != null && cost > 0 ? (profit / cost) * 100 : null
                return (
                  <tr key={r.id}>
                    <td className="text-bold">{r.item_name}</td>
                    <td><span className="category-tag">{r.category}</span></td>
                    <td style={{ color: '#f57f17', fontWeight: 700, fontSize: '0.85rem' }}>
                      {r.percentage > 0 ? `${r.percentage}%` : '—'}
                    </td>
                    <td style={{ fontWeight: 600 }}>{r.quantity.toLocaleString()} 張</td>
                    <td style={{ color: '#374151' }}>{r.bought_price.toLocaleString()}</td>
                    <td style={{ color: '#374151', fontWeight: 600 }}>{cost.toLocaleString()}</td>
                    <td style={{ color: cur != null ? '#e65100' : '#9ca3af', fontWeight: cur != null ? 600 : 400 }}>
                      {loadingPrices ? '…' : fmt(cur)}
                    </td>
                    <td style={{ color: value != null ? '#374151' : '#9ca3af', fontWeight: 600 }}>
                      {loadingPrices ? '…' : fmt(value)}
                    </td>
                    <td><DiffCell value={profit} pct={pct} /></td>
                    <td style={{ fontSize: '0.78rem', color: '#9ca3af', whiteSpace: 'nowrap' }}>
                      {new Date(r.bought_at).toLocaleDateString('zh-TW')}
                    </td>
                    <td>
                      <button
                        onClick={() => handleDelete(r.id)}
                        style={{ padding: '3px 10px', fontSize: '0.78rem', border: '1px solid #fecaca', borderRadius: 6, background: '#fef2f2', color: '#dc2626', cursor: 'pointer' }}
                      >
                        刪除
                      </button>
                    </td>
                  </tr>
                )
              })}
            </tbody>
            {records.length > 1 && (
              <tfoot>
                <tr style={{ background: '#f8f9fb', fontWeight: 700 }}>
                  <td colSpan={5} style={{ padding: '10px 16px', fontSize: '0.85rem', color: '#374151' }}>合計</td>
                  <td style={{ padding: '10px 16px', color: '#374151' }}>{totalCost.toLocaleString()}</td>
                  <td></td>
                  <td style={{ padding: '10px 16px', color: '#374151' }}>{allHavePrice ? totalValue.toLocaleString() : '—'}</td>
                  <td style={{ padding: '10px 16px' }}>
                    <DiffCell value={totalProfit} pct={totalPct} />
                  </td>
                  <td colSpan={2}></td>
                </tr>
              </tfoot>
            )}
          </table>
        </div>
      )}

      {showAdd && (
        <div className="overlay" onClick={() => setShowAdd(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h2>新增購買記錄</h2>
            </div>
            <form onSubmit={handleAdd}>
              <div className="form-group">
                <label>卷軸名稱 *</label>
                <div style={{ position: 'relative' }} ref={searchRef}>
                  <input
                    placeholder="輸入名稱搜尋"
                    value={searchText}
                    onChange={e => { setSearchText(e.target.value); setSelectedItem(null); setShowSugg(true) }}
                    onFocus={() => setShowSugg(true)}
                    autoFocus
                  />
                  {selectedItem && (
                    <div style={{ marginTop: 4, fontSize: '0.82rem', color: '#16a34a', fontWeight: 600 }}>
                      ✓ {selectedItem.name}
                      {selectedItem.percentage > 0 && ` (${selectedItem.percentage}%)`}
                      {' · '}{selectedItem.category}
                    </div>
                  )}
                  {showSugg && suggestions.length > 0 && (
                    <ul className="search-suggestions" style={{ top: 'calc(100% + 2px)' }}>
                      {suggestions.map(item => (
                        <li
                          key={item.id}
                          className="suggestion-item"
                          onMouseDown={e => {
                            e.preventDefault()
                            setSelectedItem(item)
                            setSearchText(item.name)
                            setShowSugg(false)
                          }}
                        >
                          {item.name}
                          {item.percentage > 0 && <span style={{ color: '#f57f17', marginLeft: 6, fontSize: '0.78rem' }}>{item.percentage}%</span>}
                          <span style={{ color: '#9ca3af', marginLeft: 6, fontSize: '0.78rem' }}>{item.category}</span>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              </div>
              <div className="form-group">
                <label>購買張數 *</label>
                <input
                  type="number"
                  min={1}
                  value={qty}
                  onChange={e => setQty(e.target.value)}
                  placeholder="例：10"
                />
              </div>
              <div className="form-group">
                <label>買入單價 *</label>
                <input
                  type="text"
                  inputMode="numeric"
                  value={buyPrice}
                  onChange={e => setBuyPrice(formatPrice(e.target.value))}
                  placeholder="例：1,200,000"
                />
              </div>
              <div className="form-actions">
                <button type="button" className="btn" onClick={() => setShowAdd(false)}>取消</button>
                <button type="submit" className="btn btn-primary">新增</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  )
}
