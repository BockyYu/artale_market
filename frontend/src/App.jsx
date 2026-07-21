import { useState, useEffect, useCallback, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import PotionTable from './PotionTable'
import Portfolio from './Portfolio'
import { getMemberInfo, memberLogout, memberLogin, memberFetch, fetchAppConfig, fetchPriceHistory } from './member-api'

const SCROLL_API    = '/api/v1/member/scrolls/search'
const SKILLBOOK_API = '/api/v1/member/skillbooks/search'
const EQUIP_API     = '/api/v1/member/equips/search'
const OTHER_API     = '/api/v1/member/others/search'

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
  { label: '劍士', cols: 2, items: [
    { label: '劍士',   value: '劍士' },
    { label: '英雄',   value: '英雄' },
    { label: '聖騎士', value: '聖騎士' },
    { label: '黑騎士', value: '黑騎士' },
  ]},
  { label: '弓手', cols: 3, items: [
    { label: '弓手',   value: '弓手' },
    { label: '箭神',   value: '箭神' },
    { label: '神射手', value: '神射手' },
  ]},
  { label: '法師', cols: 2, items: [
    { label: '法師',   value: '法師' },
    { label: '火毒',   value: '火毒' },
    { label: '冰雷',   value: '冰雷' },
    { label: '主教',   value: '主教' },
  ]},
  { label: '盜賊', cols: 3, items: [
    { label: '盜賊',   value: '盜賊' },
    { label: '神偷',   value: '神偷' },
    { label: '夜使者', value: '夜使者' },
  ]},
  { label: '海賊', cols: 2, items: [
    { label: '槍神',   value: '槍神' },
    { label: '拳霸',   value: '拳霸' },
  ]},
]

function MobileBlock() {
  return (
    <div style={{
      position: 'fixed', inset: 0, background: '#f0f2f5',
      display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
      padding: 32, textAlign: 'center', zIndex: 9999,
    }}>
      <div style={{ fontSize: 56, marginBottom: 20 }}>🖥️</div>
      <h1 style={{ fontSize: 22, fontWeight: 700, color: '#1a1a2e', marginBottom: 12 }}>
        請使用電腦瀏覽器開啟
      </h1>
      <p style={{ fontSize: 15, color: '#6b7280', lineHeight: 1.6 }}>
        本網站目前不支援手機或平板裝置，<br />請改用桌機或筆電瀏覽。
      </p>
    </div>
  )
}

export default function App() {
  const navigate = useNavigate()

  if (window.innerWidth < 768) return <MobileBlock />
  const [member, setMember] = useState(getMemberInfo)
  const [appConfig, setAppConfig] = useState(null)
  const [activeTab, setActiveTab] = useState('market')
  const [viewMode, setViewMode] = useState('scroll') // 'scroll' | 'skillbook'
  const [summary, setSummary] = useState([])
  const localToday = () => {
    const d = new Date()
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  }
  const prevDay = (dateStr) => {
    const d = new Date(dateStr + 'T00:00:00')
    d.setDate(d.getDate() - 1)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  }

  const [searchText, setSearchText] = useState('')
  const [filterPct, setFilterPct] = useState([])
  const [filterCategories, setFilterCategories] = useState([])
  const [sortBy, setSortBy] = useState('price_desc')

  const [showSuggestions, setShowSuggestions] = useState(false)
  const searchRef = useRef(null)
  const tableTopRef = useRef(null)

  const [allItems, setAllItems] = useState([])
  const [pinnedItems, setPinnedItems] = useState([])
  const [pinnedPrices, setPinnedPrices] = useState({})

  const [selectedJob, setSelectedJob] = useState(ALL_SKILLBOOK_JOB)
  const [skillBookItems, setSkillBookItems] = useState([])
  const [skillBookSortBy, setSkillBookSortBy] = useState('price_desc')

  const [scrollPage, setScrollPage] = useState(1)
  const [scrollPageSize, setScrollPageSize] = useState(10)
  const [scrollTotal, setScrollTotal] = useState(0)
  const [scrollDataDate, setScrollDataDate] = useState(null)
  const [skillBookPage, setSkillBookPage] = useState(1)
  const [skillBookPageSize, setSkillBookPageSize] = useState(10)
  const [skillBookTotal, setSkillBookTotal] = useState(0)
  const [skillBookDataDate, setSkillBookDataDate] = useState(null)

  const [equipItems, setEquipItems] = useState([])
  const [equipSortBy, setEquipSortBy] = useState('price_desc')
  const [equipFilterCategories, setEquipFilterCategories] = useState([])
  const [equipPage, setEquipPage] = useState(1)
  const [equipPageSize, setEquipPageSize] = useState(10)
  const [equipTotal, setEquipTotal] = useState(0)
  const [equipDataDate, setEquipDataDate] = useState(null)

  const [otherItems, setOtherItems] = useState([])
  const [otherSortBy, setOtherSortBy] = useState('price_desc')
  const [otherPage, setOtherPage] = useState(1)
  const [otherPageSize, setOtherPageSize] = useState(10)
  const [otherTotal, setOtherTotal] = useState(0)
  const [otherDataDate, setOtherDataDate] = useState(null)

  const [skillBookSearch, setSkillBookSearch] = useState('')
  const [equipSearch, setEquipSearch] = useState('')
  const [otherSearch, setOtherSearch] = useState('')
  const [otherFilterTypes, setOtherFilterTypes] = useState([])

  const [historyModal, setHistoryModal] = useState(null) // { itemId, itemName } | null

  const fetchSummary = useCallback(async (pcts, categories, sortBy, page, pageSize) => {
    try {
      let date = localToday()
      let result = null
      for (let i = 0; i < 2; i++) {
        const res = await memberFetch(SCROLL_API, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ date, percentage: pcts, category: categories.length === 0 ? ['scroll_all'] : categories, sort_by: sortBy, page, page_size: pageSize }),
        })
        result = await res.json()
        if ((result?.total || 0) > 0) break
        date = prevDay(date)
      }
      setSummary(result?.data || [])
      setScrollTotal(result?.total || 0)
      setScrollDataDate(date)
    } catch {
      setSummary([])
      setScrollTotal(0)
      setScrollDataDate(null)
    }
  }, [])

  const fetchAllItems = useCallback(async () => {
    try {
      const res = await memberFetch('/api/v1/member/items')
      const result = await res.json()
      setAllItems(result?.data || [])
    } catch {
      setAllItems([])
    }
  }, [])

  const fetchEquips = useCallback(async (categories, name, sortBy, page, pageSize) => {
    try {
      let date = localToday()
      let result = null
      for (let i = 0; i < 2; i++) {
        const res = await memberFetch(EQUIP_API, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ date, category: categories, name: name || undefined, sort_by: sortBy, page, page_size: pageSize }),
        })
        result = await res.json()
        if ((result?.total || 0) > 0) break
        date = prevDay(date)
      }
      setEquipItems(result?.data || [])
      setEquipTotal(result?.total || 0)
      setEquipDataDate(date)
    } catch {
      setEquipItems([])
      setEquipTotal(0)
      setEquipDataDate(null)
    }
  }, [])

  const fetchOthers = useCallback(async (types, name, sortBy, page, pageSize) => {
    try {
      let date = localToday()
      let result = null
      for (let i = 0; i < 2; i++) {
        const res = await memberFetch(OTHER_API, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ date, types: types.length > 0 ? types : undefined, name: name || undefined, sort_by: sortBy, page, page_size: pageSize }),
        })
        result = await res.json()
        if ((result?.total || 0) > 0) break
        date = prevDay(date)
      }
      setOtherItems(result?.data || [])
      setOtherTotal(result?.total || 0)
      setOtherDataDate(date)
    } catch {
      setOtherItems([])
      setOtherTotal(0)
      setOtherDataDate(null)
    }
  }, [])

  const fetchSkillBooks = useCallback(async (job, name, sortBy, page, pageSize) => {
    try {
      const categories = job === ALL_SKILLBOOK_JOB ? [] : [job]
      let date = localToday()
      let result = null
      for (let i = 0; i < 2; i++) {
        const res = await memberFetch(SKILLBOOK_API, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ date, category: categories, name: name || undefined, sort_by: sortBy, page, page_size: pageSize }),
        })
        result = await res.json()
        if ((result?.total || 0) > 0) break
        date = prevDay(date)
      }
      setSkillBookItems(result?.data || [])
      setSkillBookTotal(result?.total || 0)
      setSkillBookDataDate(date)
    } catch {
      setSkillBookItems([])
      setSkillBookTotal(0)
      setSkillBookDataDate(null)
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
    } else if (viewMode === 'skillbook') {
      fetchSkillBooks(selectedJob, skillBookSearch, skillBookSortBy, skillBookPage, skillBookPageSize)
    } else if (viewMode === 'equip') {
      fetchEquips(equipFilterCategories, equipSearch, equipSortBy, equipPage, equipPageSize)
    } else if (viewMode === 'other') {
      fetchOthers(otherFilterTypes, otherSearch, otherSortBy, otherPage, otherPageSize)
    }
  }, [fetchSummary, fetchSkillBooks, fetchEquips, fetchOthers,
      filterPct, filterCategories, sortBy, viewMode, selectedJob, skillBookSortBy,
      scrollPage, scrollPageSize, skillBookPage, skillBookPageSize,
      equipFilterCategories, equipSortBy, equipPage, equipPageSize,
      otherSortBy, otherPage, otherPageSize,
      skillBookSearch, equipSearch, otherSearch, otherFilterTypes, appConfig])

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
  useEffect(() => { setSkillBookPage(1) }, [selectedJob, skillBookSearch, skillBookSortBy, skillBookPageSize])
  useEffect(() => { setEquipPage(1) }, [equipFilterCategories, equipSearch, equipSortBy, equipPageSize])
  useEffect(() => { setOtherPage(1) }, [otherSearch, otherFilterTypes, otherSortBy, otherPageSize])

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

  function buildSearchRegex(keyword) {
    const escaped = keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    const pattern = escaped.replace(/(\d+)/g, '$1(?!\\d)')
    return new RegExp(pattern, 'i')
  }

  function hasEnglish(text) {
    return /[a-zA-Z]/.test(text)
  }

  function itemMatchesKeyword(item, keyword) {
    const re = buildSearchRegex(keyword)
    if (hasEnglish(keyword) && item.english_name) return re.test(item.english_name)
    return re.test(item.name)
  }

  const scrollItems = allItems.filter(i => i.item_type === 1)

  const suggestions = searchText.trim().length > 0
    ? [...new Set(
        scrollItems
          .filter(item => itemMatchesKeyword(item, searchText.trim()))
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
        { label: '頭盔', value: '頭盔' },
        { label: '上衣', value: '上衣' },
        { label: '下衣', value: '下衣' },
        { label: '套服', value: '套服' },
        { label: '鞋子', value: '鞋子' },
        { label: '手套', value: '手套' },
        { label: '披風', value: '披風' },
        { label: '盾牌', value: '盾牌' },
        { label: '臉飾', value: '臉部' },
        { label: '眼飾', value: '眼部' },
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

  const EQUIP_CATEGORY_GROUPS = [
    {
      label: '防具',
      cols: 5,
      items: [
        { label: '頭盔', value: '頭盔' },
        { label: '上衣', value: '上衣' },
        { label: '下衣', value: '下衣' },
        { label: '套服', value: '套服' },
        { label: '鞋子', value: '鞋子' },
        { label: '手套', value: '手套' },
        { label: '披風', value: '披風' },
        { label: '盾牌', value: '盾牌' },
        { label: '臉飾', value: '臉部' },
        { label: '眼飾', value: '眼部' },
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
        { label: '武器',   value: '武器' },
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
        { label: '飛鏢',   value: '飛鏢' },
      ],
    },
  ]

  const DataDateBanner = ({ date }) => {
    if (!date || date === localToday()) return null
    return (
      <div style={{
        display: 'flex', alignItems: 'center', gap: 6,
        padding: '6px 12px', marginBottom: 8, borderRadius: 6,
        background: '#fffbeb', border: '1px solid #fcd34d',
        fontSize: 13, color: '#92400e',
      }}>
        ⚠️ 今日尚無資料，顯示 <strong>{date}</strong> 的價格
      </div>
    )
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

  if (!appConfig) return null

  if (appConfig.maintenance) return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100vh', gap: 12 }}>
      <h2>系統維護中</h2>
      <p style={{ color: '#888' }}>{appConfig.message || 'We\'ll be back soon.'}</p>
    </div>
  )

  return (
    <>
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
          <div className="fs-wrap">
            <div className="fs-tabs">
              <button
                className={`fs-tab fs-tab--skillbook ${viewMode === 'skillbook' ? 'active' : ''}`}
                onClick={() => setViewMode('skillbook')}
              >職業技能書</button>
              <button
                className={`fs-tab fs-tab--scroll ${viewMode === 'scroll' ? 'active' : ''}`}
                onClick={() => setViewMode('scroll')}
              >卷軸</button>
              <button
                className={`fs-tab fs-tab--equip ${viewMode === 'equip' ? 'active' : ''}`}
                onClick={() => setViewMode('equip')}
              >裝備</button>
              <button
                className={`fs-tab fs-tab--other ${viewMode === 'other' ? 'active' : ''}`}
                onClick={() => setViewMode('other')}
              >其他</button>
            </div>

            {/* 職業 panel */}
            <div className={`fs-panel ${viewMode === 'skillbook' ? 'active' : ''}`}>
              <div className="fs-row" style={{ gridTemplateColumns: 'repeat(2, 1fr)' }}>
                <button
                  className={`fs-btn ${selectedJob === ALL_SKILLBOOK_JOB ? 'active' : ''}`}
                  onClick={() => setSelectedJob(ALL_SKILLBOOK_JOB)}
                >全部</button>
                <button
                  className={`fs-btn ${selectedJob === '全職業共通' ? 'active' : ''}`}
                  onClick={() => setSelectedJob('全職業共通')}
                >全職業通用</button>
              </div>
              {JOB_GROUPS.map((group) => (
                <div key={group.label}>
                  <div className="fs-sub-label">{group.label}</div>
                  <div className="fs-row" style={{ gridTemplateColumns: `repeat(${group.cols}, 1fr)` }}>
                    {group.items.map(({ label, value }) => (
                      <button
                        key={value}
                        className={`fs-btn ${selectedJob === value ? 'active' : ''}`}
                        onClick={() => setSelectedJob(value)}
                      >{label}</button>
                    ))}
                  </div>
                </div>
              ))}
            </div>

            {/* 卷軸 panel */}
            <div className={`fs-panel ${viewMode === 'scroll' ? 'active' : ''}`}>
              <div className="fs-sub-label">成功率</div>
              <div className="fs-row" style={{ gridTemplateColumns: 'repeat(5, 1fr)' }}>
                <button
                  className={`fs-btn ${filterPct.length === 0 ? 'active' : ''}`}
                  onClick={() => setFilterPct([])}
                >全部</button>
                {PCT_OPTIONS.map((pct) => (
                  <button
                    key={pct}
                    className={`fs-btn ${filterPct.includes(pct) ? 'active' : ''}`}
                    onClick={() => setFilterPct(prev =>
                      prev.includes(pct) ? prev.filter(p => p !== pct) : [...prev, pct]
                    )}
                  >{pct}%</button>
                ))}
              </div>

              <div className="fs-divider" />

              {filterCategories.length > 0 && (
                <button
                  className="fs-btn fs-btn-clear"
                  style={{ width: '100%', marginBottom: 4 }}
                  onClick={() => setFilterCategories([])}
                >清除分類 ×</button>
              )}

              {CATEGORY_GROUPS.map((group) => {
                const groupValues = group.items.map(i => i.value)
                const allSelected = groupValues.every(v => filterCategories.includes(v))
                return (
                  <div key={group.label}>
                    <div className="fs-sub-label">{group.label}</div>
                    <button
                      className={`fs-btn ${allSelected ? 'active' : ''}`}
                      style={{ width: '100%', marginBottom: 4 }}
                      onClick={() => setFilterCategories(prev =>
                        allSelected
                          ? prev.filter(c => !groupValues.includes(c))
                          : [...new Set([...prev, ...groupValues])]
                      )}
                    >全部</button>
                    <div className="fs-row" style={{ gridTemplateColumns: `repeat(${group.cols}, 1fr)` }}>
                      {group.items.map(({ label, value }) => (
                        <button
                          key={value}
                          className={`fs-btn ${filterCategories.includes(value) ? 'active' : ''}`}
                          onClick={() => setFilterCategories(prev =>
                            prev.includes(value) ? prev.filter(c => c !== value) : [...prev, value]
                          )}
                        >{label}</button>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
            {/* 其他 panel */}
            <div className={`fs-panel ${viewMode === 'other' ? 'active' : ''}`}>
              {otherFilterTypes.length > 0 && (
                <button
                  className="fs-btn fs-btn-clear"
                  style={{ width: '100%', marginBottom: 4 }}
                  onClick={() => setOtherFilterTypes([])}
                >清除分類 ×</button>
              )}
              <button
                className={`fs-btn ${otherFilterTypes.length === 0 ? 'active' : ''}`}
                style={{ width: '100%', marginBottom: 4 }}
                onClick={() => setOtherFilterTypes([])}
              >全部</button>
              <div className="fs-row" style={{ gridTemplateColumns: 'repeat(3, 1fr)' }}>
                {[
                  { label: '消耗', value: 3 },
                  { label: '素材', value: 2 },
                  { label: '商城', value: 5 },
                ].map(({ label, value }) => (
                  <button
                    key={value}
                    className={`fs-btn ${otherFilterTypes.includes(value) ? 'active' : ''}`}
                    onClick={() => setOtherFilterTypes(prev =>
                      prev.includes(value) ? prev.filter(t => t !== value) : [...prev, value]
                    )}
                  >{label}</button>
                ))}
              </div>
            </div>

            {/* 裝備 panel */}
            <div className={`fs-panel ${viewMode === 'equip' ? 'active' : ''}`}>
              {equipFilterCategories.length > 0 && (
                <button
                  className="fs-btn fs-btn-clear"
                  style={{ width: '100%', marginBottom: 4 }}
                  onClick={() => setEquipFilterCategories([])}
                >清除分類 ×</button>
              )}

              {EQUIP_CATEGORY_GROUPS.map((group) => {
                const groupValues = group.items.map(i => i.value)
                const allSelected = groupValues.every(v => equipFilterCategories.includes(v))
                return (
                  <div key={group.label}>
                    <div className="fs-sub-label">{group.label}</div>
                    <button
                      className={`fs-btn ${allSelected ? 'active' : ''}`}
                      style={{ width: '100%', marginBottom: 4 }}
                      onClick={() => setEquipFilterCategories(prev =>
                        allSelected
                          ? prev.filter(c => !groupValues.includes(c))
                          : [...new Set([...prev, ...groupValues])]
                      )}
                    >全部</button>
                    <div className="fs-row" style={{ gridTemplateColumns: `repeat(${group.cols}, 1fr)` }}>
                      {group.items.map(({ label, value }) => (
                        <button
                          key={value}
                          className={`fs-btn ${equipFilterCategories.includes(value) ? 'active' : ''}`}
                          onClick={() => setEquipFilterCategories(prev =>
                            prev.includes(value) ? prev.filter(c => c !== value) : [...prev, value]
                          )}
                        >{label}</button>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </aside>

        <div className="main-content">

          <div className="filter-bar">
            {viewMode === 'skillbook' && (
              <input
                className="search-input"
                placeholder="搜尋技能書名稱"
                value={skillBookSearch}
                onChange={e => setSkillBookSearch(e.target.value)}
              />
            )}
            {viewMode === 'equip' && (
              <input
                className="search-input"
                placeholder="搜尋裝備名稱"
                value={equipSearch}
                onChange={e => setEquipSearch(e.target.value)}
              />
            )}
            {viewMode === 'other' && (
              <input
                className="search-input"
                placeholder="搜尋道具名稱"
                value={otherSearch}
                onChange={e => setOtherSearch(e.target.value)}
              />
            )}
            {viewMode === 'scroll' && (
              <div className="search-wrapper" ref={searchRef}>
                <input
                  className="search-input"
                  placeholder="搜尋卷軸名稱"
                  value={searchText}
                  onChange={(e) => { setSearchText(e.target.value); setShowSuggestions(true) }}
                  onFocus={() => setShowSuggestions(true)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      const kw = searchText.trim().toLowerCase()
                      if (kw) {
                        const matched = scrollItems.filter(item => {
                          const keywords = kw.split(/\s+/)
                          return keywords.every(k => {
                            const re = buildSearchRegex(k)
                            if (hasEnglish(k) && item.english_name) return re.test(item.english_name)
                            return re.test(`${item.name} ${item.category}`)
                          })
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
                          const item = scrollItems.find(i => i.name === name)
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

          <div ref={tableTopRef} />
          {viewMode === 'equip' ? (
            <>
              <DataDateBanner date={equipDataDate} />
              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th style={{ width: 36, textAlign: 'center', color: '#9ca3af' }}>#</th>
                      <th>裝備名稱</th>
                      <th>類型</th>
                      <th
                        className="sortable-th"
                        onClick={() => setEquipSortBy(s => s === 'price_desc' ? 'price_asc' : 'price_desc')}
                      >
                        今日價格
                        <span className="sort-icon">
                          {equipSortBy === 'price_desc' ? ' ▼' : equipSortBy === 'price_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th
                        className="sortable-th"
                        onClick={() => setEquipSortBy(s => s === 'yesterday_price_desc' ? 'yesterday_price_asc' : 'yesterday_price_desc')}
                      >
                        昨日
                        <span className="sort-icon">
                          {equipSortBy === 'yesterday_price_desc' ? ' ▼' : equipSortBy === 'yesterday_price_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th
                        className="sortable-th"
                        onClick={() => setEquipSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
                      >
                        漲跌
                        <span className="sort-icon">
                          {equipSortBy === 'change_desc' ? ' ▼' : equipSortBy === 'change_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th style={{ width: 60 }}></th>
                    </tr>
                  </thead>
                  <tbody>
                    {equipItems.length === 0 ? (
                      <tr>
                        <td colSpan={7} className="empty">尚無資料</td>
                      </tr>
                    ) : (
                      equipItems.map((item, idx) => (
                        <tr key={item.item_id}>
                          <td style={{ textAlign: 'center', color: '#9ca3af', fontSize: '0.88rem', fontWeight: 600 }}>
                            {(equipPage - 1) * equipPageSize + idx + 1}
                          </td>
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
                          <td style={{ textAlign: 'center' }}>
                            <button className="history-btn" onClick={() => setHistoryModal({ itemId: item.item_id, itemName: item.item_name })}>歷史資料</button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
              <PaginationBar
                page={equipPage}
                pageSize={equipPageSize}
                total={equipTotal}
                onPageChange={p => { tableTopRef.current?.scrollIntoView({ behavior: 'instant' }); setEquipPage(p) }}
                onPageSizeChange={setEquipPageSize}
              />
            </>
          ) : viewMode === 'scroll' ? (
            <>
              <DataDateBanner date={scrollDataDate} />
              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th style={{ width: 36, textAlign: 'center', color: '#9ca3af' }}>#</th>
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
                      <th
                        className="sortable-th"
                        onClick={() => setSortBy(s => s === 'yesterday_price_desc' ? 'yesterday_price_asc' : 'yesterday_price_desc')}
                      >
                        昨日
                        <span className="sort-icon">
                          {sortBy === 'yesterday_price_desc' ? ' ▼' : sortBy === 'yesterday_price_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th
                        className="sortable-th"
                        onClick={() => setSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
                      >
                        漲跌
                        <span className="sort-icon">
                          {sortBy === 'change_desc' ? ' ▼' : sortBy === 'change_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th style={{ width: 60 }}></th>
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
                      filteredSummary.map((item, idx) => (
                          <tr key={item.item_id}>
                            <td style={{ textAlign: 'center', color: '#9ca3af', fontSize: '0.88rem', fontWeight: 600 }}>
                              {pinnedItems.length > 0 ? idx + 1 : (scrollPage - 1) * scrollPageSize + idx + 1}
                            </td>
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
                            <td style={{ textAlign: 'center' }}>
                              <button className="history-btn" onClick={() => setHistoryModal({ itemId: item.item_id, itemName: item.item_name })}>歷史資料</button>
                            </td>
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
                onPageChange={p => { tableTopRef.current?.scrollIntoView({ behavior: 'instant' }); setScrollPage(p) }}
                onPageSizeChange={setScrollPageSize}
              />
            </>
          ) : viewMode === 'other' ? (
            <>
              <DataDateBanner date={otherDataDate} />
              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th style={{ width: 36, textAlign: 'center', color: '#9ca3af' }}>#</th>
                      <th>道具名稱</th>
                      <th>分類</th>
                      <th
                        className="sortable-th"
                        onClick={() => setOtherSortBy(s => s === 'price_desc' ? 'price_asc' : 'price_desc')}
                      >
                        今日價格
                        <span className="sort-icon">{otherSortBy === 'price_desc' ? ' ▼' : otherSortBy === 'price_asc' ? ' ▲' : ' ⇅'}</span>
                      </th>
                      <th
                        className="sortable-th"
                        onClick={() => setOtherSortBy(s => s === 'yesterday_price_desc' ? 'yesterday_price_asc' : 'yesterday_price_desc')}
                      >
                        昨日價格
                        <span className="sort-icon">{otherSortBy === 'yesterday_price_desc' ? ' ▼' : otherSortBy === 'yesterday_price_asc' ? ' ▲' : ' ⇅'}</span>
                      </th>
                      <th
                        className="sortable-th"
                        onClick={() => setOtherSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
                      >
                        漲跌
                        <span className="sort-icon">{otherSortBy === 'change_desc' ? ' ▼' : otherSortBy === 'change_asc' ? ' ▲' : ' ⇅'}</span>
                      </th>
                      <th style={{ width: 60 }}></th>
                    </tr>
                  </thead>
                  <tbody>
                    {otherItems.length === 0 ? (
                      <tr><td colSpan={7} className="empty">尚無資料</td></tr>
                    ) : (
                      otherItems.map((item, idx) => (
                        <tr key={item.item_id}>
                          <td style={{ textAlign: 'center', color: '#9ca3af', fontSize: '0.88rem', fontWeight: 600 }}>
                            {(otherPage - 1) * otherPageSize + idx + 1}
                          </td>
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
                          <td style={{ textAlign: 'center' }}>
                            <button className="history-btn" onClick={() => setHistoryModal({ itemId: item.item_id, itemName: item.item_name })}>歷史資料</button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
              <PaginationBar
                page={otherPage}
                pageSize={otherPageSize}
                total={otherTotal}
                onPageChange={p => { tableTopRef.current?.scrollIntoView({ behavior: 'instant' }); setOtherPage(p) }}
                onPageSizeChange={setOtherPageSize}
              />
            </>
          ) : (
            <>
              <DataDateBanner date={skillBookDataDate} />
              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th style={{ width: 36, textAlign: 'center', color: '#9ca3af' }}>#</th>
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
                      <th
                        className="sortable-th"
                        onClick={() => setSkillBookSortBy(s => s === 'yesterday_price_desc' ? 'yesterday_price_asc' : 'yesterday_price_desc')}
                      >
                        昨日
                        <span className="sort-icon">
                          {skillBookSortBy === 'yesterday_price_desc' ? ' ▼' : skillBookSortBy === 'yesterday_price_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th
                        className="sortable-th"
                        onClick={() => setSkillBookSortBy(s => s === 'change_desc' ? 'change_asc' : 'change_desc')}
                      >
                        漲跌
                        <span className="sort-icon">
                          {skillBookSortBy === 'change_desc' ? ' ▼' : skillBookSortBy === 'change_asc' ? ' ▲' : ' ⇅'}
                        </span>
                      </th>
                      <th style={{ width: 60 }}></th>
                    </tr>
                  </thead>
                  <tbody>
                    {sortedSkillBooks.length === 0 ? (
                      <tr>
                        <td colSpan={7} className="empty">尚無資料</td>
                      </tr>
                    ) : (
                      sortedSkillBooks.map((item, idx) => (
                          <tr key={item.item_id}>
                            <td style={{ textAlign: 'center', color: '#9ca3af', fontSize: '0.88rem', fontWeight: 600 }}>
                              {(skillBookPage - 1) * skillBookPageSize + idx + 1}
                            </td>
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
                            <td style={{ textAlign: 'center' }}>
                              <button className="history-btn" onClick={() => setHistoryModal({ itemId: item.item_id, itemName: item.item_name })}>歷史資料</button>
                            </td>
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
                onPageChange={p => { tableTopRef.current?.scrollIntoView({ behavior: 'instant' }); setSkillBookPage(p) }}
                onPageSizeChange={setSkillBookPageSize}
              />
            </>
          )}

        </div>{/* main-content */}
      </div>}{/* activeTab === 'market' */}

    </div>

    {historyModal && (
      <PriceHistoryModal
        itemId={historyModal.itemId}
        itemName={historyModal.itemName}
        onClose={() => setHistoryModal(null)}
      />
    )}
    </>
  )
}

function PriceHistoryModal({ itemId, itemName, onClose }) {
  const [days, setDays] = useState(7)
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)   // 初次載入
  const [fetching, setFetching] = useState(false) // 切換天數時的背景更新

  useEffect(() => {
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = '' }
  }, [])

  useEffect(() => {
    // 已有資料時只做背景更新（不清空畫面），初次載入才顯示 loading
    if (data.length > 0) {
      setFetching(true)
      fetchPriceHistory(itemId, days)
        .then(setData)
        .catch(() => {})
        .finally(() => setFetching(false))
    } else {
      setLoading(true)
      fetchPriceHistory(itemId, days)
        .then(setData)
        .catch(() => setData([]))
        .finally(() => setLoading(false))
    }
  }, [itemId, days])

  const rows = data.map((r, i) => {
    const prev = data[i + 1]
    const changePct = prev && prev.price > 0
      ? (r.price - prev.price) / prev.price * 100
      : null
    return { ...r, changePct }
  })

  return (
    <div
      style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.45)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }}
    >
      <div
        style={{ background: '#fff', borderRadius: 12, width: 1200, height: '72vh', display: 'flex', flexDirection: 'column', boxShadow: '0 20px 60px rgba(0,0,0,0.25)' }}
        onClick={e => e.stopPropagation()}
      >
        <div style={{ padding: '32px 48px 28px', borderBottom: '1px solid #e5e7eb', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <div>
            <div style={{ fontWeight: 700, fontSize: 20, color: '#1a1a2e' }}>{itemName}</div>
            <div style={{ fontSize: 13, color: '#9ca3af', marginTop: 4 }}>歷史價格紀錄</div>
          </div>
          <button onClick={onClose} style={{ border: 'none', background: 'none', fontSize: 22, cursor: 'pointer', color: '#9ca3af', lineHeight: 1, padding: '0 2px' }}>×</button>
        </div>

        <div style={{ padding: '18px 48px', borderBottom: '1px solid #e5e7eb', display: 'flex', gap: 14, alignItems: 'center' }}>
          {[7, 14, 30].map(d => (
            <button
              key={d}
              onClick={() => setDays(d)}
              style={{
                padding: '12px 20px', borderRadius: 28, border: '1px solid',
                borderColor: days === d ? '#4f46e5' : '#e5e7eb',
                background: days === d ? '#4f46e5' : '#fff',
                color: days === d ? '#fff' : '#6b7280',
                fontSize: 13, cursor: 'pointer', fontWeight: days === d ? 600 : 400,
                textAlign: 'center',
              }}
            >{d} 天</button>
          ))}
          <span style={{ marginLeft: 'auto', fontSize: 20, color: '#9ca3af' }}>每日最低價</span>
        </div>

        <div style={{ overflowY: 'auto', flex: 1 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: 40, color: '#9ca3af', fontSize: 14 }}>載入中…</div>
          ) : rows.length === 0 ? (
            <div style={{ textAlign: 'center', padding: 40, color: '#9ca3af', fontSize: 14 }}>尚無歷史資料</div>
          ) : (
            <table style={{ width: '100%', borderCollapse: 'collapse', opacity: fetching ? 0.5 : 1, transition: 'opacity 0.15s' }}>
              <thead>
                <tr style={{ background: '#f9fafb', position: 'sticky', top: 0 }}>
                  <th style={histThStyle}>日期</th>
                  <th style={{ ...histThStyle, textAlign: 'right' }}>最低價</th>
                  <th style={{ ...histThStyle, textAlign: 'right' }}>更新時間</th>
                  <th style={{ ...histThStyle, textAlign: 'right' }}>漲跌</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((r, i) => (
                  <tr key={r.id} style={{ borderTop: '1px solid #f3f4f6', background: i === 0 ? '#f5f3ff' : 'transparent' }}>
                    <td style={histTdStyle}>
                      {new Date(r.recorded_date).toLocaleDateString('zh-TW', { month: '2-digit', day: '2-digit' })}
                      {i === 0 && (
                        <span style={{ marginLeft: 6, fontSize: 11, color: '#7c3aed', background: '#ede9fe', padding: '2px 8px', borderRadius: 4, fontWeight: 600 }}>最新</span>
                      )}
                    </td>
                    <td style={{ ...histTdStyle, textAlign: 'right', fontWeight: 700, color: '#111827' }}>
                      {r.price.toLocaleString()}
                    </td>
                    <td style={{ ...histTdStyle, textAlign: 'right', color: '#9ca3af' }}>
                      {new Date(r.updated_at || r.created_at).toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' })}
                    </td>
                    <td style={{ ...histTdStyle, textAlign: 'right' }}>
                      {r.changePct == null ? (
                        <span style={{ color: '#d1d5db' }}>—</span>
                      ) : r.changePct > 0 ? (
                        <span style={{ color: '#ef4444' }}>▲ {Math.abs(r.changePct).toFixed(1)}%</span>
                      ) : r.changePct < 0 ? (
                        <span style={{ color: '#22c55e' }}>▼ {Math.abs(r.changePct).toFixed(1)}%</span>
                      ) : (
                        <span style={{ color: '#9ca3af' }}>—</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>

        <div style={{ padding: '20px 48px', borderTop: '1px solid #e5e7eb', display: 'flex', justifyContent: 'flex-end' }}>
          <button onClick={onClose} style={{ border: 'none', background: '#ef4444', fontSize: 13, cursor: 'pointer', color: '#fff', fontWeight: 700, padding: '6px 18px', borderRadius: 6, letterSpacing: 1 }}>關閉</button>
        </div>
      </div>
    </div>
  )
}

const histThStyle = { padding: '8px 24px', textAlign: 'left', fontSize: 13, color: '#6b7280', fontWeight: 600 }
const histTdStyle = { padding: '10px 24px', fontSize: 14, color: '#374151' }

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
