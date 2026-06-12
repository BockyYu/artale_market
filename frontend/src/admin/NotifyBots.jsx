import { useState, useEffect, useCallback } from 'react'
import { listBots, createBot, updateBot, deleteBot, toggleBotActive, sendBotMessage, testDiscordWebhook } from './api'

const PLATFORM_LABEL = { tg: 'Telegram', line: 'LINE Notify' }
const PLATFORM_COLOR = { tg: '#2ca5e0', line: '#06c755' }

const EMPTY_FORM = { name: '', platform: 'tg', token: '', chat_id: '' }

function BotForm({ f, setF, isEdit }) {
  return (
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
        </select>
      </div>
      <div className="form-group">
        <label>
          {f.platform === 'tg' && 'Bot Token *'}
          {f.platform === 'line' && 'LINE Notify Token *'}
        </label>
        <input
          required={!isEdit}
          value={f.token}
          onChange={e => setF(p => ({ ...p, token: e.target.value }))}
          placeholder={
            isEdit ? '留空則不修改' :
            f.platform === 'tg' ? '例：123456789:ABCdef...' :
            'LINE Notify Access Token'
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
}

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
  const [sendTarget, setSendTarget] = useState(null)
  const [sendMsg, setSendMsg] = useState('')
  const [sending, setSending] = useState(false)
  const [testingDc, setTestingDc] = useState(false)

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

  async function handleSend(e) {
    e.preventDefault()
    if (!sendMsg.trim()) return
    setSending(true)
    try {
      await sendBotMessage(sendTarget.id, sendMsg)
      setSendTarget(null)
      setSendMsg('')
      alert('訊息已送出')
    } catch (err) {
      alert(err.message)
    } finally {
      setSending(false)
    }
  }

  function openEdit(bot) {
    setEditingBot(bot)
    setEditForm({ name: bot.name, platform: bot.platform, token: '', chat_id: bot.chat_id })
  }

  return (
    <>
      <div className="page-header">
        <h1>通知機器人</h1>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <span style={{ fontSize: 13, color: '#374151', fontWeight: 600 }}>共 {bots.length} 筆</span>
          <button
            className="btn-action"
            style={{ background: '#5865f2', color: '#fff', border: 'none', opacity: testingDc ? 0.6 : 1 }}
            disabled={testingDc}
            onClick={async () => {
              setTestingDc(true)
              try {
                await testDiscordWebhook()
                alert('Discord 測試訊息已送出')
              } catch (err) {
                alert('發送失敗：' + err.message)
              } finally {
                setTestingDc(false)
              }
            }}
          >
            {testingDc ? '發送中...' : '測試 Discord'}
          </button>
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
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {loading && <tr className="empty-row"><td colSpan={6}>載入中...</td></tr>}
            {!loading && bots.length === 0 && <tr className="empty-row"><td colSpan={6}>尚無機器人</td></tr>}
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
                <td style={{ whiteSpace: 'nowrap' }}>
                  <button className="btn-action btn-edit" style={{ marginRight: 8 }} onClick={() => openEdit(bot)}>修改</button>
                  <button
                    className="btn-action"
                    style={{ marginRight: 8, color: '#7c3aed' }}
                    onClick={() => { setSendTarget(bot); setSendMsg('') }}
                  >
                    發送訊息
                  </button>
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

      {sendTarget && (
        <div className="modal-overlay" onClick={() => setSendTarget(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>發送訊息 — {sendTarget.name}</h2>
            <p style={{ color: '#6b7280', fontSize: 14, marginBottom: 16 }}>
              平台：<span style={{ color: PLATFORM_COLOR[sendTarget.platform], fontWeight: 700 }}>
                {PLATFORM_LABEL[sendTarget.platform] ?? sendTarget.platform}
              </span>
            </p>
            <form onSubmit={handleSend}>
              <div className="form-group">
                <label>訊息內容 *</label>
                <textarea
                  required
                  rows={4}
                  style={{ width: '100%', resize: 'vertical', padding: '8px 12px', borderRadius: 6, border: '1px solid #d1d5db', fontSize: 14, fontFamily: 'inherit', boxSizing: 'border-box' }}
                  value={sendMsg}
                  onChange={e => setSendMsg(e.target.value)}
                  placeholder="輸入要發送的訊息..."
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setSendTarget(null)}>取消</button>
                <button type="submit" className="btn-save" disabled={sending}>{sending ? '發送中...' : '發送'}</button>
              </div>
            </form>
          </div>
        </div>
      )}

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
