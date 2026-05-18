import { useState, useEffect, useCallback } from 'react'

const API = '/api/items'

const emptyForm = { name: '', description: '', price: '', quantity: '' }

export default function App() {
  const [items, setItems] = useState([])
  const [modal, setModal] = useState(null) // null | 'create' | 'edit'
  const [form, setForm] = useState(emptyForm)
  const [editId, setEditId] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const fetchItems = useCallback(async () => {
    try {
      const res = await fetch(API)
      const data = await res.json()
      setItems(Array.isArray(data) ? data : [])
    } catch {
      setItems([])
    }
  }, [])

  useEffect(() => {
    fetchItems()
  }, [fetchItems])

  const openCreate = () => {
    setForm(emptyForm)
    setError('')
    setModal('create')
  }

  const openEdit = (item) => {
    setForm({
      name: item.name,
      description: item.description,
      price: String(item.price),
      quantity: String(item.quantity),
    })
    setEditId(item.id)
    setError('')
    setModal('edit')
  }

  const closeModal = () => {
    setModal(null)
    setEditId(null)
    setError('')
  }

  const handleChange = (e) =>
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }))

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    const body = {
      name: form.name,
      description: form.description,
      price: parseFloat(form.price),
      quantity: parseInt(form.quantity) || 0,
    }

    try {
      const res = await fetch(modal === 'create' ? API : `${API}/${editId}`, {
        method: modal === 'create' ? 'POST' : 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })

      if (!res.ok) {
        const err = await res.json()
        setError(err.error || '操作失敗，請重試')
        return
      }

      await fetchItems()
      closeModal()
    } catch {
      setError('無法連接到伺服器')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id, name) => {
    if (!window.confirm(`確定要刪除「${name}」嗎？`)) return
    try {
      await fetch(`${API}/${id}`, { method: 'DELETE' })
      await fetchItems()
    } catch {
      alert('刪除失敗')
    }
  }

  return (
    <div className="container">
      <header className="header">
        <div className="header-title">
          <h1>🏪 Artale Market</h1>
          <span className="badge">{items.length} 件商品</span>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          + 新增商品
        </button>
      </header>

      <div className="table-wrapper">
        <table>
          <thead>
            <tr>
              <th style={{ width: 60 }}>ID</th>
              <th>名稱</th>
              <th>描述</th>
              <th style={{ width: 140 }}>價格</th>
              <th style={{ width: 80 }}>數量</th>
              <th style={{ width: 140 }}>操作</th>
            </tr>
          </thead>
          <tbody>
            {items.length === 0 ? (
              <tr>
                <td colSpan={6} className="empty">
                  尚無商品，點擊「新增商品」開始吧！
                </td>
              </tr>
            ) : (
              items.map((item) => (
                <tr key={item.id}>
                  <td className="text-center text-muted">{item.id}</td>
                  <td className="text-bold">{item.name}</td>
                  <td className="text-muted">{item.description || '—'}</td>
                  <td className="text-price">
                    {item.price.toLocaleString()} 楓幣
                  </td>
                  <td className="text-center">{item.quantity}</td>
                  <td>
                    <div className="action-btns">
                      <button
                        className="btn btn-edit"
                        onClick={() => openEdit(item)}
                      >
                        編輯
                      </button>
                      <button
                        className="btn btn-delete"
                        onClick={() => handleDelete(item.id, item.name)}
                      >
                        刪除
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {modal && (
        <div className="overlay" onClick={closeModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>{modal === 'create' ? '新增商品' : '編輯商品'}</h2>
              <button className="close-btn" onClick={closeModal}>
                ✕
              </button>
            </div>

            {error && <p className="error-msg">{error}</p>}

            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>名稱 *</label>
                <input
                  name="name"
                  value={form.name}
                  onChange={handleChange}
                  placeholder="商品名稱"
                  required
                />
              </div>
              <div className="form-group">
                <label>描述</label>
                <textarea
                  name="description"
                  value={form.description}
                  onChange={handleChange}
                  placeholder="商品描述（選填）"
                  rows={3}
                />
              </div>
              <div className="form-row">
                <div className="form-group">
                  <label>價格（楓幣）*</label>
                  <input
                    name="price"
                    type="number"
                    value={form.price}
                    onChange={handleChange}
                    placeholder="0"
                    min="1"
                    step="1"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>數量</label>
                  <input
                    name="quantity"
                    type="number"
                    value={form.quantity}
                    onChange={handleChange}
                    placeholder="0"
                    min="0"
                  />
                </div>
              </div>
              <div className="form-actions">
                <button
                  type="button"
                  className="btn"
                  onClick={closeModal}
                  disabled={loading}
                >
                  取消
                </button>
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={loading}
                >
                  {loading ? '處理中...' : modal === 'create' ? '新增' : '儲存'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
