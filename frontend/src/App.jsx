import { useState, useEffect, useCallback, useRef } from 'react'
import PotionTable from './PotionTable'

const SUMMARY_API  = '/api/prices/summary'
const FREQUENT_API = '/api/me/frequent-items'

function getUserID() {
  let id = localStorage.getItem('artale_uid')
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem('artale_uid', id)
  }
  return id
}

const USER_ID = getUserID()

export default function App() {
  const [activeTab, setActiveTab] = useState('market')
  const [summary, setSummary] = useState([])
  const [modal, setModal] = useState(null)
  const [priceInput, setPriceInput] = useState('')
  const [selectedItem, setSelectedItem] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const localToday = () => {
    const d = new Date()
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  }
  const [filterDate, setFilterDate] = useState(localToday)
  const [searchText, setSearchText] = useState('')
  const [filterPct, setFilterPct] = useState([])
  const [filterCategories, setFilterCategories] = useState([])
  const [sortBy, setSortBy] = useState('price_desc')

  const [frequentItems, setFrequentItems] = useState([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const searchRef = useRef(null)

  const [allItems, setAllItems] = useState([])
  const [pinnedItems, setPinnedItems] = useState([])

  const fetchSummary = useCallback(async (date, pcts, categories) => {
    try {
      const url = new URL(SUMMARY_API, window.location.origin)
      url.searchParams.set('date', date)
      if (pcts && pcts.length > 0) url.searchParams.set('percentage', pcts.join(','))
      if (categories && categories.length > 0) url.searchParams.set('category', categories.join(','))
      const res = await fetch(url.toString())
      setSummary(await res.json() || [])
    } catch {
      setSummary([])
    }
  }, [])

  const fetchFrequent = useCallback(async () => {
    try {
      const res = await fetch(FREQUENT_API, { headers: { 'X-User-ID': USER_ID } })
      setFrequentItems(await res.json() || [])
    } catch {
      setFrequentItems([])
    }
  }, [])

  const fetchAllItems = useCallback(async (date) => {
    try {
      const url = new URL(SUMMARY_API, window.location.origin)
      url.searchParams.set('date', date)
      const res = await fetch(url.toString())
      setAllItems(await res.json() || [])
    } catch {
      setAllItems([])
    }
  }, [])

  useEffect(() => {
    fetchSummary(filterDate, filterPct, filterCategories)
    fetchAllItems(filterDate)
    fetchFrequent()
  }, [fetchSummary, fetchAllItems, fetchFrequent, filterDate, filterPct, filterCategories])

  useEffect(() => {
    const handleClick = (e) => {
      if (searchRef.current && !searchRef.current.contains(e.target)) {
        setShowSuggestions(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  const pinItems = (items) => {
    setPinnedItems(prev => {
      const existingIds = new Set(prev.map(p => p.item_id))
      const added = items.filter(i => !existingIds.has(i.item_id))
      return added.length ? [...prev, ...added] : prev
    })
  }

  const closeModal = () => {
    setModal(null)
    setSelectedItem(null)
    setError('')
  }

  const openRecordPrice = (item) => {
    setSelectedItem(item)
    setPriceInput(item.today_price != null ? String(item.today_price) : '')
    setError('')
    setModal('recordPrice')
  }

  const handlePriceSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const res = await fetch(`/api/items/${selectedItem.item_id}/prices`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'X-User-ID': USER_ID },
        body: JSON.stringify({ price: parseFloat(priceInput) }),
      })
      if (!res.ok) {
        setError((await res.json()).error || '記錄失敗')
        return
      }
      await fetchSummary(filterDate, filterPct, filterCategories)
      fetchAllItems(filterDate)
      fetchFrequent()
      closeModal()
    } catch {
      setError('無法連接到伺服器')
    } finally {
      setLoading(false)
    }
  }

  const suggestions = searchText.trim().length > 0
    ? [...new Set(
        allItems
          .filter(item => item.item_name.toLowerCase().includes(searchText.trim().toLowerCase()))
          .map(item => item.item_name)
      )].slice(0, 8)
    : []

  const sortItems = (items) => {
    if (sortBy === 'price_desc') {
      return [...items].sort((a, b) => {
        if (a.today_price == null && b.today_price == null) return 0
        if (a.today_price == null) return 1
        if (b.today_price == null) return -1
        return b.today_price - a.today_price
      })
    }
    if (sortBy === 'price_asc') {
      return [...items].sort((a, b) => {
        if (a.today_price == null && b.today_price == null) return 0
        if (a.today_price == null) return 1
        if (b.today_price == null) return -1
        return a.today_price - b.today_price
      })
    }
    if (sortBy === 'change_desc') {
      return [...items].sort((a, b) => {
        if (a.change_percent == null && b.change_percent == null) return 0
        if (a.change_percent == null) return 1
        if (b.change_percent == null) return -1
        return b.change_percent - a.change_percent
      })
    }
    if (sortBy === 'change_asc') {
      return [...items].sort((a, b) => {
        if (a.change_percent == null && b.change_percent == null) return 0
        if (a.change_percent == null) return 1
        if (b.change_percent == null) return -1
        return a.change_percent - b.change_percent
      })
    }
    return items
  }

  const filteredSummary = sortItems(
    pinnedItems.length > 0
      ? pinnedItems.map(p => allItems.find(i => i.item_id === p.item_id) ?? p)
      : summary
  )

  const PCT_OPTIONS = [10, 30, 60, 100]

  const CATEGORY_GROUPS = [
    {
      label: '防具',
      cols: 6,
      items: [
        { label: '帽',   value: '頭盔' },
        { label: '上衣', value: '上衣' },
        { label: '下衣', value: '下衣' },
        { label: '套服', value: '套服' },
        { label: '鞋子', value: '鞋子' },
        { label: '手套', value: '手套' },
        { label: '披風', value: '披風' },
        { label: '盾牌', value: '盾牌' },
        { label: '臉飾', value: '臉部裝飾' },
        { label: '眼飾', value: '眼部裝飾' },
        { label: '耳環', value: '耳環' },
        { label: '戒指', value: '戒指' },
        { label: '墜飾', value: '墜飾' },
        { label: '腰帶', value: '腰帶' },
        { label: '肩章', value: '肩章' },
        { label: '勳章', value: '勳章' },
      ],
    },
    {
      label: '武器',
      cols: 3,
      items: [
        { label: '單手劍', value: '單手劍' },
        { label: '雙手劍', value: '雙手劍' },
        { label: '單手斧', value: '單手斧' },
        { label: '雙手斧', value: '雙手斧' },
        { label: '單手棍', value: '單手棍' },
        { label: '雙手棍', value: '雙手棍' },
        { label: '槍',     value: '槍' },
        { label: '矛',     value: '矛' },
        { label: '短杖',   value: '短杖' },
        { label: '長杖',   value: '長杖' },
        { label: '弓',     value: '弓' },
        { label: '弩',     value: '弩' },
        { label: '短劍',   value: '短劍' },
        { label: '拳套',   value: '拳套' },
        { label: '指虎',   value: '指虎' },
        { label: '火槍',   value: '火槍' },
      ],
    },
  ]

  const daysAgo = (n) => {
    const d = new Date()
    d.setDate(d.getDate() - n)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  }

  const fmt = (price) =>
    price != null ? price.toLocaleString() : '—'

  const ChangeCell = ({ pct }) => {
    if (pct == null) return <span className="text-muted">—</span>
    const up = pct >= 0
    return (
      <span className={up ? 'change-up' : 'change-down'}>
        {up ? '▲' : '▼'} {Math.abs(pct).toFixed(2)}%
      </span>
    )
  }

  const today = new Date().toLocaleDateString('zh-TW', {
    year: 'numeric', month: 'long', day: 'numeric',
  })

  return (
    <div className="container">
      <header className="header">
        <div className="header-left">
          <h1>🏪 Artale Market</h1>
          <span className="date-label">{today}</span>
        </div>
        <nav className="tab-nav">
          <button
            className={`tab-btn ${activeTab === 'market' ? 'active' : ''}`}
            onClick={() => setActiveTab('market')}
          >
            市場行情
          </button>
          <button
            className={`tab-btn ${activeTab === 'potion' ? 'active' : ''}`}
            onClick={() => setActiveTab('potion')}
          >
            藥水參考
          </button>
        </nav>
      </header>

      {activeTab === 'potion' && <PotionTable />}

      {activeTab === 'market' && frequentItems.length > 0 && (
        <div className="frequent-bar">
          <span className="frequent-label">常用</span>
          {frequentItems.map((fi) => {
            const matched = summary.find((s) => s.item_id === fi.item_id)
            return (
              <button
                key={fi.item_id}
                className="frequent-chip"
                onClick={() => matched && openRecordPrice(matched)}
                title={`已查詢 ${fi.count} 次`}
              >
                {fi.name}
                <span className="frequent-pct">{fi.percentage}%</span>
                <span className="frequent-count">{fi.count}次</span>
              </button>
            )
          })}
        </div>
      )}

      {activeTab === 'market' && <div className="main-layout">
        <aside className="sidebar">
          <div className="sidebar-title">成功率</div>
          <div className="pct-grid">
            <button
              className={`pct-filter-btn ${filterPct.length === 0 ? 'active' : ''}`}
              onClick={() => setFilterPct([])}
            >
              全部
            </button>
            {PCT_OPTIONS.map((pct) => (
              <button
                key={pct}
                className={`pct-filter-btn ${filterPct.includes(pct) ? 'active' : ''}`}
                onClick={() => setFilterPct(prev =>
                  prev.includes(pct) ? prev.filter(p => p !== pct) : [...prev, pct]
                )}
              >
                {pct}%
              </button>
            ))}
          </div>

          <div className="sidebar-divider" />

          {filterCategories.length > 0 && (
            <button
              className="cat-clear-btn"
              onClick={() => setFilterCategories([])}
            >
              清除分類 ×
            </button>
          )}

          {CATEGORY_GROUPS.map((group) => (
            <div key={group.label}>
              <div className="sidebar-group-label">{group.label}</div>
              <div className="cat-grid" style={{ gridTemplateColumns: `repeat(${group.cols}, 1fr)` }}>
                {group.items.map(({ label, value }) => (
                  <button
                    key={value}
                    className={`cat-filter-btn ${filterCategories.includes(value) ? 'active' : ''}`}
                    onClick={() => setFilterCategories(prev =>
                      prev.includes(value) ? prev.filter(c => c !== value) : [...prev, value]
                    )}
                  >
                    {label}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </aside>

        <div className="main-content">

      <div className="filter-bar">
        <div className="search-wrapper" ref={searchRef}>
          <input
            className="search-input"
            placeholder="搜尋商品名稱或類型，可用空格分隔多個關鍵字"
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
              setShowSuggestions(true)
            }}
            onFocus={() => setShowSuggestions(true)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                const kw = searchText.trim().toLowerCase()
                if (kw) {
                  const matched = allItems.filter(item => {
                    const keywords = kw.split(/\s+/)
                    return keywords.every(k => `${item.item_name} ${item.category}`.toLowerCase().includes(k))
                  })
                  if (matched.length > 0) pinItems(matched)
                  setSearchText('')
                }
                setShowSuggestions(false)
              }
            }}
          />
          {showSuggestions && suggestions.length > 0 && (
            <ul className="search-suggestions">
              {suggestions.map((name) => (
                <li
                  key={name}
                  className="suggestion-item"
                  onMouseDown={(e) => {
                    e.preventDefault()
                    const item = allItems.find(i => i.item_name === name)
                    if (item) pinItems([item])
                    setSearchText('')
                    setShowSuggestions(false)
                  }}
                >
                  {name}
                </li>
              ))}
            </ul>
          )}
        </div>
        <div className="date-filter">
          <div className="quick-dates">
            {[7, 14, 30].map((n) => {
              const d = daysAgo(n)
              return (
                <button
                  key={n}
                  className={`quick-date-btn ${filterDate === d ? 'active' : ''}`}
                  onClick={() => setFilterDate(d)}
                >
                  {n} 天前
                </button>
              )
            })}
            <button
              className={`quick-date-btn ${filterDate === localToday() ? 'active' : ''}`}
              onClick={() => setFilterDate(localToday())}
            >
              今天
            </button>
          </div>
          <input
            type="date"
            className="date-input"
            value={filterDate}
            max={new Date().toISOString().slice(0, 10)}
            onChange={(e) => setFilterDate(e.target.value)}
          />
        </div>
      </div>

      {pinnedItems.length > 0 && (
        <div className="pinned-bar">
          {pinnedItems.map(pinned => {
            const fresh = summary.find(i => i.item_id === pinned.item_id)
                          ?? allItems.find(i => i.item_id === pinned.item_id)
                          ?? pinned
            return (
              <div key={pinned.item_id} className="pinned-chip">
                <button
                  className="pinned-chip-name"
                  onClick={() => openRecordPrice(fresh)}
                >
                  {pinned.item_name}
                  {fresh.today_price != null && (
                    <span className="pinned-price">{fresh.today_price.toLocaleString()}</span>
                  )}
                </button>
                <button
                  className="pinned-chip-remove"
                  onClick={() => setPinnedItems(prev => prev.filter(p => p.item_id !== pinned.item_id))}
                >×</button>
              </div>
            )
          })}
        </div>
      )}

      <div className="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>商品名稱</th>
              <th>類型</th>
              <th
                className="sortable-th"
                onClick={() => setSortBy(s => s === 'price_desc' ? 'price_asc' : 'price_desc')}
              >
                今日價格
                <span className="sort-icon">
                  {sortBy === 'price_desc' ? ' ▼' : sortBy === 'price_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th>昨日</th>
              <th>三天前</th>
              <th
                className="sortable-th"
                onClick={() => setSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
              >
                漲跌
                <span className="sort-icon">
                  {sortBy === 'change_desc' ? ' ▼' : sortBy === 'change_asc' ? ' ▲' : ' ⇅'}
                </span>
              </th>
              <th style={{ width: 140 }}>操作</th>
            </tr>
          </thead>
          <tbody>
            {filteredSummary.length === 0 ? (
              <tr>
                <td colSpan={7} className="empty">
                  {summary.length === 0 ? '尚無商品' : '找不到符合的商品'}
                </td>
              </tr>
            ) : (
              filteredSummary.map((item) => (
                <tr key={item.item_id}>
                  <td className="text-bold">{item.item_name}</td>
                  <td>
                    <span className="category-tag">{item.category}</span>
                  </td>
                  <td className={item.today_price != null ? 'text-price' : 'text-muted'}>
                    {fmt(item.today_price)}
                  </td>
                  <td className="text-muted">{fmt(item.yesterday_price)}</td>
                  <td className="text-muted">{fmt(item.three_days_ago_price)}</td>
                  <td>
                    <ChangeCell pct={item.change_percent} />
                  </td>
                  <td>
                    <button
                      className="btn btn-record"
                      onClick={() => openRecordPrice(item)}
                    >
                      {item.today_price != null ? '更新今日價' : '記錄今日價'}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

        </div>{/* main-content */}
      </div>}{/* activeTab === 'market' */}

      {modal && (
        <div className="overlay" onClick={closeModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>{`📝 ${selectedItem?.item_name}`}</h2>
              <button className="close-btn" onClick={closeModal}>✕</button>
            </div>

            {error && <p className="error-msg">{error}</p>}

            <form onSubmit={handlePriceSubmit}>
              {selectedItem?.yesterday_price != null && (
                <div className="price-hint">
                  昨日：<strong>{selectedItem.yesterday_price.toLocaleString()} 楓幣</strong>
                  {selectedItem.change_percent != null && (
                    <span className={selectedItem.change_percent >= 0 ? 'change-up' : 'change-down'}>
                      {' '}
                      {selectedItem.change_percent >= 0 ? '▲' : '▼'}{' '}
                      {Math.abs(selectedItem.change_percent).toFixed(2)}%
                    </span>
                  )}
                </div>
              )}
              <div className="form-group">
                <label>今日市場價格（楓幣）*</label>
                <input
                  type="number"
                  value={priceInput}
                  onChange={(e) => setPriceInput(e.target.value)}
                  placeholder="輸入今日觀察到的市場價格"
                  min="1"
                  step="1"
                  required
                  autoFocus
                />
              </div>
              <div className="form-actions">
                <button type="button" className="btn" onClick={closeModal} disabled={loading}>
                  取消
                </button>
                <button type="submit" className="btn btn-primary" disabled={loading}>
                  {loading ? '記錄中...' : '確認記錄'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
