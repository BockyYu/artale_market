import { useState, useEffect, useCallback } from 'react'
import { listBots, createBot, updateBot, deleteBot, toggleBotActive } from './api'

const PLATFORM_LABEL = { tg: 'Telegram', line: 'LINE Notify', dc: 'Discord' }
const PLATFORM_COLOR = { tg: '#2ca5e0', line: '#06c755', dc: '#5865f2' }

const EMPTY_FORM = { name: '', platform: 'tg', token: '', chat_id: '' }

export default function NotifyBots() {
  const [bots, setBots] = useState([])
  const [loading, setLoading] = useState(false)
  const [showCreate, setShowCreate] = useState(false)
  const [form, setForm] = useState(EMPTY_FORM)
  const [creating, setCreating] = useState(false)
  const [editingBot, setEditingBot] = useState(null)
  const [editForm, setEditForm] = useState(EMPTY_FORM)
  const [saving, setSaving] = useState(false)
  const [toggling, setToggling] = useState(null)
  const [deleting, setDeleting] = useState(null)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const data = await listBots()
      setBots(data?.data || [])
    } catch (err) {
      alert(err.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  async function handleCreate(e) {
    e.preventDefault()
    setCreating(true)
    try {
      const created = await createBot({
        name: form.name,
        platform: form.platform,
        token: form.token,
        chat_id: form.chat_id,
      })
      setBots(prev => [created, ...prev])
      setShowCreate(false)
      setForm(EMPTY_FORM)
    } catch (err) {
      alert(err.message)
    } finally {
      setCreating(false)
    }
  }

  async function handleSave(e) {
    e.preventDefault()
    setSaving(true)
    try {
      await updateBot(editingBot.id, {
        name: editForm.name,
        platform: editForm.platform,
        token: editForm.token,
        chat_id: editForm.chat_id,
      })
      await load()
      setEditingBot(null)
    } catch (err) {
      alert(err.message)
    } finally {
      setSaving(false)
    }
  }

  async function handleToggle(bot) {
    setToggling(bot.id)
    try {
      await toggleBotActive(bot.id, !bot.is_active)
      setBots(prev => prev.map(b => b.id === bot.id ? { ...b, is_active: !b.is_active } : b))
    } catch (err) {
      alert(err.message)
    } finally {
      setToggling(null)
    }
  }

  async function handleDelete(id) {
    if (!confirm('確定要刪除這個機器人？')) return
    setDeleting(id)
    try {
      await deleteBot(id)
      setBots(prev => prev.filter(b => b.id !== id))
    } catch (err) {
      alert(err.message)
    } finally {
      setDeleting(null)
    }
  }

  function openEdit(bot) {
    setEditingBot(bot)
    setEditForm({ name: bot.name, platform: bot.platform, token: '', chat_id: bot.chat_id })
  }

  const BotForm = ({ f, setF, isEdit }) => (
    <>
      <div className="form-group">
        <label>機器人名稱 *</label>
        <input required value={f.name} onChange={e => setF(p => ({ ...p, name: e.target.value }))} placeholder="例：主要通知" />
      </div>
      <div className="form-group">
        <label>平台 *</label>
        <select
          className="search-input"
          style={{ width: '100%', maxWidth: '100%' }}
          value={f.platform}
          onChange={e => setF(p => ({ ...p, platform: e.target.value, chat_id: '' }))}
        >
          <option value="tg">Telegram</option>
          <option value="line">LINE Notify</option>
          <option value="dc">Discord</option>
        </select>
      </div>
      <div className="form-group">
        <label>
          {f.platform === 'tg' && 'Bot Token *'}
          {f.platform === 'line' && 'LINE Notify Token *'}
          {f.platform === 'dc' && 'Webhook URL *'}
        </label>
        <input
          required={!isEdit}
          value={f.token}
          onChange={e => setF(p => ({ ...p, token: e.target.value }))}
          placeholder={
            isEdit ? '留空則不修改' :
            f.platform === 'tg' ? '例：123456789:ABCdef...' :
            f.platform === 'line' ? 'LINE Notify Access Token' :
            'https://discord.com/api/webhooks/...'
          }
        />
      </div>
      {f.platform === 'tg' && (
        <div className="form-group">
          <label>Chat ID *</label>
          <input
            required={!isEdit}
            value={f.chat_id}
            onChange={e => setF(p => ({ ...p, chat_id: e.target.value }))}
            placeholder="例：-100123456789"
          />
        </div>
      )}
    </>
  )

  return (
    <>
      <div className="page-header">
        <h1>通知機器人</h1>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 13, color: '#374151', fontWeight: 600 }}>共 {bots.length} 筆</span>
          <button className="btn-add" onClick={() => { setForm(EMPTY_FORM); setShowCreate(true) }}>+ 新增機器人</button>
        </div>
      </div>

      <div className="card">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>名稱</th>
              <th>平台</th>
              <th>Chat ID</th>
              <th>狀態</th>
              <th>修改</th>
              <th>刪除</th>
            </tr>
          </thead>
          <tbody>
            {loading && <tr className="empty-row"><td colSpan={7}>載入中...</td></tr>}
            {!loading && bots.length === 0 && <tr className="empty-row"><td colSpan={7}>尚無機器人</td></tr>}
            {bots.map(bot => (
              <tr key={bot.id}>
                <td>{bot.id}</td>
                <td className="text-bold">{bot.name}</td>
                <td>
                  <span style={{
                    display: 'inline-block', padding: '2px 10px', borderRadius: 12, fontSize: 12, fontWeight: 700,
                    background: PLATFORM_COLOR[bot.platform] + '22',
                    color: PLATFORM_COLOR[bot.platform],
                  }}>
                    {PLATFORM_LABEL[bot.platform] ?? bot.platform}
                  </span>
                </td>
                <td style={{ fontSize: 12, color: '#6b7280' }}>{bot.chat_id || '—'}</td>
                <td>
                  <button
                    disabled={toggling === bot.id}
                    className={`btn-action ${bot.is_active ? 'btn-active-priority' : ''}`}
                    style={{
                      background: bot.is_active ? '#16a34a' : undefined,
                      color: bot.is_active ? '#fff' : undefined,
                      border: bot.is_active ? 'none' : undefined,
                      opacity: toggling === bot.id ? 0.5 : 1,
                    }}
                    onClick={() => handleToggle(bot)}
                  >
                    {bot.is_active ? '啟用中' : '已停用'}
                  </button>
                </td>
                <td>
                  <button className="btn-action btn-edit" onClick={() => openEdit(bot)}>修改</button>
                </td>
                <td>
                  <button
                    className="btn-action"
                    style={{ color: '#dc2626', opacity: deleting === bot.id ? 0.5 : 1 }}
                    disabled={deleting === bot.id}
                    onClick={() => handleDelete(bot.id)}
                  >
                    刪除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showCreate && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>新增通知機器人</h2>
            <form onSubmit={handleCreate}>
              <BotForm f={form} setF={setForm} isEdit={false} />
              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setShowCreate(false)}>取消</button>
                <button type="submit" className="btn-save" disabled={creating}>{creating ? '建立中...' : '建立'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {editingBot && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>修改機器人 — {editingBot.name}</h2>
            <form onSubmit={handleSave}>
              <BotForm f={editForm} setF={setEditForm} isEdit={true} />
              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setEditingBot(null)}>取消</button>
                <button type="submit" className="btn-save" disabled={saving}>{saving ? '儲存中...' : '儲存'}</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  )
}
