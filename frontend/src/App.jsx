import { useState, useEffect, useCallback, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import PotionTable from './PotionTable'
import { getMemberInfo, memberLogout, memberLogin, memberFetch, fetchAppConfig } from './member-api'

const SCROLL_API    = '/api/v1/member/scrolls/search'
const SKILLBOOK_API = '/api/v1/member/skillbooks/search'

function getUserID() {
  let id = localStorage.getItem('artale_uid')
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem('artale_uid', id)
  }
  return id
}

const USER_ID = getUserID()

const ALL_SKILLBOOK_JOB = '全部'

const JOB_GROUPS = [
  { label: '全職業', cols: 1, items: [
    { label: '全職業', value: '全職業共通' },
  ]},
  { label: '劍士系', cols: 2, items: [
    { label: '劍士',   value: '劍士' },
    { label: '英雄',   value: '英雄' },
    { label: '聖騎士', value: '聖騎士' },
    { label: '黑騎士', value: '黑騎士' },
  ]},
  { label: '弓手系', cols: 3, items: [
    { label: '弓手',   value: '弓手' },
    { label: '箭神',   value: '箭神' },
    { label: '神射手', value: '神射手' },
  ]},
  { label: '法師系', cols: 2, items: [
    { label: '法師',   value: '法師' },
    { label: '火毒',   value: '火毒' },
    { label: '冰雷',   value: '冰雷' },
    { label: '主教',   value: '主教' },
  ]},
  { label: '盜賊系', cols: 3, items: [
    { label: '盜賊',   value: '盜賊' },
    { label: '神偷',   value: '神偷' },
    { label: '夜使者', value: '夜使者' },
  ]},
  { label: '海賊系', cols: 2, items: [
    { label: '槍神',   value: '槍神' },
    { label: '拳霸',   value: '拳霸' },
  ]},
]

export default function App() {
  const navigate = useNavigate()
  const [member, setMember] = useState(getMemberInfo)
  const [appConfig, setAppConfig] = useState(null)
  const [activeTab, setActiveTab] = useState('market')
  const [viewMode, setViewMode] = useState('scroll') // 'scroll' | 'skillbook'
  const [summary, setSummary] = useState([])
  const localToday = () => {
    const d = new Date()
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  }
  const [searchText, setSearchText] = useState('')
  const [filterPct, setFilterPct] = useState([])
  const [filterCategories, setFilterCategories] = useState([])
  const [sortBy, setSortBy] = useState('price_desc')

  const [showSuggestions, setShowSuggestions] = useState(false)
  const searchRef = useRef(null)

  const [allItems, setAllItems] = useState([])
  const [pinnedItems, setPinnedItems] = useState([])
  const [pinnedPrices, setPinnedPrices] = useState({})

  const [selectedJob, setSelectedJob] = useState(ALL_SKILLBOOK_JOB)
  const [skillBookItems, setSkillBookItems] = useState([])
  const [skillBookSortBy, setSkillBookSortBy] = useState('percentage_asc')

  const [scrollPage, setScrollPage] = useState(1)
  const [scrollPageSize, setScrollPageSize] = useState(10)
  const [scrollTotal, setScrollTotal] = useState(0)
  const [skillBookPage, setSkillBookPage] = useState(1)
  const [skillBookPageSize, setSkillBookPageSize] = useState(10)
  const [skillBookTotal, setSkillBookTotal] = useState(0)

  const fetchSummary = useCallback(async (pcts, categories, sortBy, page, pageSize) => {
    try {
      const res = await memberFetch(SCROLL_API, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ date: localToday(), percentage: pcts, category: categories.length === 0 ? ['scroll_all'] : categories, sort_by: sortBy, page, page_size: pageSize }),
      })
      const result = await res.json()
      setSummary(result?.data || [])
      setScrollTotal(result?.total || 0)
    } catch {
      setSummary([])
      setScrollTotal(0)
    }
  }, [])

  const fetchAllItems = useCallback(async () => {
    try {
      const res = await memberFetch('/api/v1/member/items')
      const result = await res.json()
      setAllItems(result || [])
    } catch {
      setAllItems([])
    }
  }, [])

  const fetchSkillBooks = useCallback(async (job, sortBy, page, pageSize) => {
    try {
      const categories = job === ALL_SKILLBOOK_JOB ? [] : [job]
      const res = await memberFetch(SKILLBOOK_API, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ date: localToday(), category: categories, sort_by: sortBy, page, page_size: pageSize }),
      })
      const result = await res.json()
      setSkillBookItems(result?.data || [])
      setSkillBookTotal(result?.total || 0)
    } catch {
      setSkillBookItems([])
      setSkillBookTotal(0)
    }
  }, [])

  useEffect(() => {
    fetchAppConfig().then(setAppConfig)
  }, [])

  useEffect(() => {
    if (!appConfig || appConfig.maintenance) return
    fetchAllItems()
  }, [fetchAllItems, appConfig])

  useEffect(() => {
    if (!appConfig || appConfig.maintenance) return
    if (viewMode === 'scroll') {
      fetchSummary(filterPct, filterCategories, sortBy, scrollPage, scrollPageSize)
    } else {
      fetchSkillBooks(selectedJob, skillBookSortBy, skillBookPage, skillBookPageSize)
    }
  }, [fetchSummary, fetchSkillBooks,
      filterPct, filterCategories, sortBy, viewMode, selectedJob, skillBookSortBy,
      scrollPage, scrollPageSize, skillBookPage, skillBookPageSize, appConfig])

  useEffect(() => {
    const handleClick = (e) => {
      if (searchRef.current && !searchRef.current.contains(e.target)) {
        setShowSuggestions(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  useEffect(() => { setScrollPage(1) }, [filterPct, filterCategories, sortBy, pinnedItems.length, scrollPageSize])
  useEffect(() => { setSkillBookPage(1) }, [selectedJob, skillBookSortBy, skillBookPageSize])

  const fetchPinnedItemPrices = useCallback(async (items) => {
    if (items.length === 0) return
    const results = await Promise.all(items.map(async (item) => {
      try {
        const res = await memberFetch(`/api/v1/member/items/${item.id}/prices`)
        return await res.json()
      } catch {
        return { item_id: item.id, item_name: item.name, category: item.category }
      }
    }))
    setPinnedPrices(prev => {
      const next = { ...prev }
      for (const r of results) next[r.item_id] = r
      return next
    })
  }, [])

  const pinItems = useCallback((items) => {
    const existingIds = new Set(pinnedItems.map(p => p.id))
    const added = items.filter(i => !existingIds.has(i.id))
    if (!added.length) return
    setPinnedItems(prev => [...prev, ...added])
    fetchPinnedItemPrices(added)
  }, [pinnedItems, fetchPinnedItemPrices])

  const suggestions = searchText.trim().length > 0
    ? [...new Set(
        allItems
          .filter(item => item.name.toLowerCase().includes(searchText.trim().toLowerCase()))
          .map(item => item.name)
      )].slice(0, 8)
    : []

  const sortItems = (items, by) => {
    if (by === 'price_desc') {
      return [...items].sort((a, b) => {
        if (a.today_price == null && b.today_price == null) return 0
        if (a.today_price == null) return 1
        if (b.today_price == null) return -1
        return b.today_price - a.today_price
      })
    }
    if (by === 'price_asc') {
      return [...items].sort((a, b) => {
        if (a.today_price == null && b.today_price == null) return 0
        if (a.today_price == null) return 1
        if (b.today_price == null) return -1
        return a.today_price - b.today_price
      })
    }
    if (by === 'change_desc') {
      return [...items].sort((a, b) => {
        if (a.change_percent == null && b.change_percent == null) return 0
        if (a.change_percent == null) return 1
        if (b.change_percent == null) return -1
        return b.change_percent - a.change_percent
      })
    }
    if (by === 'change_asc') {
      return [...items].sort((a, b) => {
        if (a.change_percent == null && b.change_percent == null) return 0
        if (a.change_percent == null) return 1
        if (b.change_percent == null) return -1
        return a.change_percent - b.change_percent
      })
    }
    return items
  }

  const filteredSummary = pinnedItems.length > 0
    ? sortItems(
        pinnedItems.map(p => pinnedPrices[p.id] ?? { item_id: p.id, item_name: p.name, category: p.category }),
        sortBy
      )
    : summary

  const sortedSkillBooks = skillBookItems

  const getPageNumbers = (current, total) => {
    if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
    const pages = [1]
    if (current > 3) pages.push('...')
    for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) pages.push(i)
    if (current < total - 2) pages.push('...')
    pages.push(total)
    return pages
  }

  const PaginationBar = ({ page, pageSize, total, onPageChange, onPageSizeChange }) => {
    const totalPages = Math.ceil(total / pageSize)
    if (total === 0) return null
    const pageNums = getPageNumbers(page, totalPages)
    return (
      <div className="pagination-bar">
        <div className="page-size-selector">
          <span className="pagination-label">每頁</span>
          {[10, 20, 40, 60, 80, 100].map(size => (
            <button
              key={size}
              className={`page-size-btn ${pageSize === size ? 'active' : ''}`}
              onClick={() => onPageSizeChange(size)}
            >{size}</button>
          ))}
        </div>
        <div className="page-nav">
          <button className="page-btn" disabled={page === 1} onClick={() => onPageChange(page - 1)}>←</button>
          {pageNums.map((p, i) =>
            p === '...'
              ? <span key={`e${i}`} className="page-ellipsis">…</span>
              : <button key={p} className={`page-btn ${page === p ? 'active' : ''}`} onClick={() => onPageChange(p)}>{p}</button>
          )}
          <button className="page-btn" disabled={page === totalPages} onClick={() => onPageChange(page + 1)}>→</button>
        </div>
        <span className="pagination-info">共 {total} 筆</span>
      </div>
    )
  }

  const PCT_OPTIONS = [10, 30, 60, 100]

  const CATEGORY_GROUPS = [
    {
      label: '防具',
      cols: 5,
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

  if (!appConfig) return null

  if (appConfig.maintenance) return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100vh', gap: 12 }}>
      <h2>系統維護中</h2>
      <p style={{ color: '#888' }}>{appConfig.message || 'We\'ll be back soon.'}</p>
    </div>
  )

  return (
    <div className="container">
      {/* {!member && <LoginModal onLogin={setMember} />} */}
      <header className="header">
        <div className="header-left">
          <h1>🏪 Artale Market</h1>
          <span className="date-label">{today}</span>
        </div>
        <div className="header-right">
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
          {member ? (
            <div className="member-bar">
              <span className="member-nickname">{member.nickname}</span>
              <button className="member-logout-btn" onClick={async () => {
                await memberLogout()
                setMember(null)
              }}>登出</button>
            </div>
          ) : null}
        </div>
      </header>

      {activeTab === 'potion' && <PotionTable />}

{activeTab === 'market' && <div className="main-layout">
        <aside className="sidebar">

          {viewMode === 'scroll' ? (
            <>
              <div className="sidebar-title">卷軸成功率</div>
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
            </>
          ) : (
            <>
              <div className="sidebar-title">職業</div>
              <button
                className={`cat-filter-btn ${selectedJob === ALL_SKILLBOOK_JOB ? 'active' : ''}`}
                style={{ marginBottom: 4 }}
                onClick={() => setSelectedJob(ALL_SKILLBOOK_JOB)}
              >
                全部
              </button>
              {JOB_GROUPS.map((group) => (
                <div key={group.label}>
                  <div className="sidebar-group-label">{group.label}</div>
                  <div className="cat-grid" style={{ gridTemplateColumns: `repeat(${group.cols}, 1fr)` }}>
                    {group.items.map(({ label, value }) => (
                      <button
                        key={value}
                        className={`cat-filter-btn ${selectedJob === value ? 'active' : ''}`}
                        onClick={() => setSelectedJob(value)}
                      >
                        {label}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </>
          )}

          <div className="sidebar-divider" />
          <button
            className={`skillbook-toggle-btn ${viewMode === 'skillbook' ? 'active' : ''}`}
            onClick={() => setViewMode(v => v === 'scroll' ? 'skillbook' : 'scroll')}
          >
            {viewMode === 'scroll' ? '技能書' : '← 回卷軸'}
          </button>

        </aside>

        <div className="main-content">

          <div className="filter-bar">
            {viewMode === 'scroll' && (
              <div className="search-wrapper" ref={searchRef}>
                <input
                  className="search-input"
                  placeholder="搜尋商品名稱"
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
                          return keywords.every(k => `${item.name} ${item.category}`.toLowerCase().includes(k))
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
                          const item = allItems.find(i => i.name === name)
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
            )}
          </div>

          {viewMode === 'scroll' && pinnedItems.length > 0 && (
            <div className="pinned-bar">
              <button
                className="pinned-clear-all"
                onClick={() => setPinnedItems([])}
              >
                清除全部
              </button>
              {pinnedItems.map(pinned => {
                const fresh = summary.find(i => i.item_id === pinned.id) ?? pinned
                return (
                  <div key={pinned.id} className="pinned-chip">
                    <span className="pinned-chip-name">
                      {pinned.name}
                      {fresh.today_price != null && (
                        <span className="pinned-price">{fresh.today_price.toLocaleString()}</span>
                      )}
                    </span>
                    <button
                      className="pinned-chip-remove"
                      onClick={() => setPinnedItems(prev => prev.filter(p => p.id !== pinned.id))}
                    >×</button>
                  </div>
                )
              })}
            </div>
          )}

          {viewMode === 'scroll' ? (
            <>
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
                      <th
                        className="sortable-th"
                        onClick={() => setSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
                      >
                        漲跌
                        <span className="sort-icon">
                          {sortBy === 'change_desc' ? ' ▼' : sortBy === 'change_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {filteredSummary.length === 0 ? (
                      <tr>
                        <td colSpan={5} className="empty">
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
                              {(item.today_updated_at || item.today_created_at) && (
                                <div className="price-updated-at">
                                  {new Date(item.today_updated_at ?? item.today_created_at).toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' })}
                                </div>
                              )}
                            </td>
                            <td className="text-muted">
                              {fmt(item.yesterday_price)}
                              {(item.yesterday_updated_at || item.yesterday_created_at) && (
                                <div className="price-updated-at">
                                  {new Date(item.yesterday_updated_at ?? item.yesterday_created_at).toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' })}
                                </div>
                              )}
                            </td>
                            <td><ChangeCell pct={item.change_percent} /></td>
                          </tr>
                        ))
                    )}
                  </tbody>
                </table>
              </div>
              <PaginationBar
                page={scrollPage}
                pageSize={scrollPageSize}
                total={pinnedItems.length > 0 ? filteredSummary.length : scrollTotal}
                onPageChange={setScrollPage}
                onPageSizeChange={setScrollPageSize}
              />
            </>
          ) : (
            <>
              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th>技能書名稱</th>
                      <th>職業</th>
                      <th
                        className="sortable-th"
                        onClick={() => setSkillBookSortBy(s => s === 'price_desc' ? 'price_asc' : 'price_desc')}
                      >
                        今日價格
                        <span className="sort-icon">
                          {skillBookSortBy === 'price_desc' ? ' ▼' : skillBookSortBy === 'price_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th>昨日</th>
                      <th
                        className="sortable-th"
                        onClick={() => setSkillBookSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
                      >
                        漲跌
                        <span className="sort-icon">
                          {skillBookSortBy === 'change_desc' ? ' ▼' : skillBookSortBy === 'change_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {sortedSkillBooks.length === 0 ? (
                      <tr>
                        <td colSpan={5} className="empty">尚無資料</td>
                      </tr>
                    ) : (
                      sortedSkillBooks.map((item) => (
                          <tr key={item.item_id}>
                            <td className="text-bold">{item.item_name}</td>
                            <td><span className="category-tag">{item.category}</span></td>
                            <td className={item.today_price != null ? 'text-price' : 'text-muted'}>
                              {fmt(item.today_price)}
                              {(item.today_updated_at || item.today_created_at) && (
                                <div className="price-updated-at">
                                  {new Date(item.today_updated_at ?? item.today_created_at).toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' })}
                                </div>
                              )}
                            </td>
                            <td className="text-muted">
                              {fmt(item.yesterday_price)}
                              {(item.yesterday_updated_at || item.yesterday_created_at) && (
                                <div className="price-updated-at">
                                  {new Date(item.yesterday_updated_at ?? item.yesterday_created_at).toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' })}
                                </div>
                              )}
                            </td>
                            <td><ChangeCell pct={item.change_percent} /></td>
                          </tr>
                        ))
                    )}
                  </tbody>
                </table>
              </div>
              <PaginationBar
                page={skillBookPage}
                pageSize={skillBookPageSize}
                total={skillBookTotal}
                onPageChange={setSkillBookPage}
                onPageSizeChange={setSkillBookPageSize}
              />
            </>
          )}

        </div>{/* main-content */}
      </div>}{/* activeTab === 'market' */}

    </div>
  )
}

function LoginModal({ onLogin }) {
  const [form, setForm] = useState({ username: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const member = await memberLogin(form.username, form.password)
      onLogin(member)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-modal-overlay">
      <div className="login-modal-card">
        <h2 className="login-modal-title">🏪 Artale Market</h2>
        <p className="login-modal-sub">請登入以繼續使用</p>
        {error && <div className="login-modal-error">{error}</div>}
        <form onSubmit={handleSubmit}>
          <div className="login-modal-field">
            <label>帳號</label>
            <input
              type="text"
              value={form.username}
              onChange={e => setForm(f => ({ ...f, username: e.target.value }))}
              placeholder="請輸入帳號"
              autoFocus
              required
            />
          </div>
          <div className="login-modal-field">
            <label>密碼</label>
            <input
              type="password"
              value={form.password}
              onChange={e => setForm(f => ({ ...f, password: e.target.value }))}
              placeholder="請輸入密碼"
              required
            />
          </div>
          <button className="login-modal-btn" type="submit" disabled={loading}>
            {loading ? '登入中...' : '登入'}
          </button>
        </form>
      </div>
    </div>
  )
}
