import { useState, useEffect, useCallback } from 'react'
import { listMembers, updateMemberStatus, deleteMember } from './api'

export default function Members() {
  const [data, setData] = useState({ data: [], total: 0 })
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [searchInput, setSearchInput] = useState('')
  const [loading, setLoading] = useState(false)
  const PAGE_SIZE = 20

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await listMembers(page, PAGE_SIZE, search)
      setData(res)
    } catch (err) {
      alert(err.message)
    } finally {
      setLoading(false)
    }
  }, [page, search])

  useEffect(() => { load() }, [load])

  function handleSearch(e) {
    e.preventDefault()
    setPage(1)
    setSearch(searchInput)
  }

  async function toggleStatus(member) {
    const newStatus = member.status === 1 ? 0 : 1
    const action = newStatus === 0 ? '停用' : '啟用'
    if (!confirm(`確定${action}此會員？`)) return
    try {
      await updateMemberStatus(member.id, newStatus)
      load()
    } catch (err) {
      alert(err.message)
    }
  }

  async function handleDelete(id) {
    if (!confirm('確定刪除此會員？此操作無法復原。')) return
    try {
      await deleteMember(id)
      if (data.data.length === 1 && page > 1) setPage(p => p - 1)
      else load()
    } catch (err) {
      alert(err.message)
    }
  }

  const totalPages = Math.ceil(data.total / PAGE_SIZE) || 1

  return (
    <>
      <div className="page-header">
        <h1>會員列表</h1>
        <span style={{ fontSize: 13, color: '#6b7280' }}>共 {data.total} 位會員</span>
      </div>

      <div className="card">
        <div className="card-toolbar">
          <form onSubmit={handleSearch} style={{ display: 'flex', gap: 8 }}>
            <input
              className="search-input"
              placeholder="搜尋帳號 / Email"
              value={searchInput}
              onChange={e => setSearchInput(e.target.value)}
            />
            <button type="submit" className="btn-add" style={{ padding: '8px 16px' }}>搜尋</button>
          </form>
        </div>

        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>帳號</th>
              <th>Email</th>
              <th>狀態</th>
              <th>加入時間</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr className="empty-row"><td colSpan={6}>載入中...</td></tr>
            )}
            {!loading && data.data.length === 0 && (
              <tr className="empty-row"><td colSpan={6}>目前無會員資料</td></tr>
            )}
            {data.data.map(m => (
              <tr key={m.id}>
                <td>{m.id}</td>
                <td>{m.username}</td>
                <td>{m.email || '—'}</td>
                <td>
                  <span className={`badge badge-${m.status === 1 ? 'active' : 'banned'}`}>
                    {m.status === 1 ? '正常' : '停用'}
                  </span>
                </td>
                <td>{new Date(m.created_at).toLocaleDateString('zh-TW')}</td>
                <td>
                  <button
                    className={`btn-action ${m.status === 1 ? 'btn-ban' : 'btn-unban'}`}
                    onClick={() => toggleStatus(m)}
                  >
                    {m.status === 1 ? '停用' : '啟用'}
                  </button>
                  <button className="btn-action btn-delete" onClick={() => handleDelete(m.id)}>刪除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {data.total > PAGE_SIZE && (
          <div className="pagination">
            <span>第 {page} / {totalPages} 頁</span>
            <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>‹ 上一頁</button>
            <button disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>下一頁 ›</button>
          </div>
        )}
      </div>
    </>
  )
}
