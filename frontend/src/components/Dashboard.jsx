import { useState, useEffect } from 'react'

function Dashboard({ onLogout }) {
  const [links, setLinks] = useState([])
  const [newName, setNewName] = useState('')
  const [newLink, setNewLink] = useState('')
  const [editingId, setEditingId] = useState(null)
  const [editingName, setEditingName] = useState('')
  const [editingLink, setEditingLink] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(true)
  const [showChangePassword, setShowChangePassword] = useState(false)
  const [copiedId, setCopiedId] = useState(null)

  useEffect(() => {
    fetchLinks()
  }, [])

  const fetchLinks = async () => {
    try {
      const response = await fetch('/api/list/vpn')
      if (!response.ok) {
        throw new Error('Ошибка загрузки ссылок')
      }
      const data = await response.json()
      setLinks(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    if (!newName.trim() || !newLink.trim()) {
      setError('Введите название и ссылку')
      return
    }

    try {
      const response = await fetch('/api/insert/vpn', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newName, link: newLink }),
      })

      if (!response.ok) {
        throw new Error('Ошибка добавления ссылки')
      }

      setSuccess('Ссылка успешно добавлена')
      setNewName('')
      setNewLink('')
      fetchLinks()
    } catch (err) {
      setError(err.message)
    }
  }

  const handleUpdate = async (id) => {
    setError('')
    setSuccess('')

    if (!editingName.trim() || !editingLink.trim()) {
      setError('Введите название и ссылку')
      return
    }

    try {
      const response = await fetch('/api/update/vpn', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id, name: editingName, link: editingLink }),
      })

      if (!response.ok) {
        throw new Error('Ошибка обновления ссылки')
      }

      setSuccess('Ссылка успешно обновлена')
      setEditingId(null)
      setEditingName('')
      setEditingLink('')
      fetchLinks()
    } catch (err) {
      setError(err.message)
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Удалить эту ссылку?')) {
      return
    }

    setError('')
    setSuccess('')

    try {
      const response = await fetch(`/api/delete/vpn?id=${id}`, {
        method: 'DELETE',
      })

      if (!response.ok) {
        throw new Error('Ошибка удаления ссылки')
      }

      setSuccess('Ссылка успешно удалена')
      fetchLinks()
    } catch (err) {
      setError(err.message)
    }
  }

  const startEdit = (link) => {
    setEditingId(link.id)
    setEditingName(link.name)
    setEditingLink(link.link)
  }

  const cancelEdit = () => {
    setEditingId(null)
    setEditingName('')
    setEditingLink('')
  }

  const handleCopyLink = (link) => {
    navigator.clipboard.writeText(link)
    setCopiedId(link)
    setTimeout(() => setCopiedId(null), 2000)
  }

  const handleChangePassword = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    const formData = new FormData(e.target)
    const oldPassword = formData.get('old_password')
    const newPassword = formData.get('new_password')
    const confirmPassword = formData.get('confirm_password')

    if (!oldPassword || !newPassword || !confirmPassword) {
      setError('Заполните все поля')
      return
    }

    if (newPassword !== confirmPassword) {
      setError('Новые пароли не совпадают')
      return
    }

    try {
      const response = await fetch('/api/change-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          old_password: oldPassword,
          new_password: newPassword,
          password_confirm: confirmPassword,
        }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.error || 'Ошибка смены пароля')
      }

      setSuccess('Пароль успешно изменен')
      setShowChangePassword(false)
      e.target.reset()
    } catch (err) {
      setError(err.message)
    }
  }

  if (loading) {
    return (
      <div className="container">
        <div className="loading">Загрузка...</div>
      </div>
    )
  }

  return (
    <div className="container">
      <div className="header">
        <h1>VPN Управление</h1>
        <div className="header-buttons">
          <button className="btn-change-password" onClick={() => setShowChangePassword(true)}>
            Сменить пароль
          </button>
          <button className="btn-logout" onClick={onLogout}>
            Выйти
          </button>
        </div>
      </div>

      {error && <div className="error">{error}</div>}
      {success && <div className="success">{success}</div>}

      {showChangePassword && (
        <div className="modal-overlay" onClick={() => setShowChangePassword(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>Смена пароля</h2>
              <button className="modal-close" onClick={() => setShowChangePassword(false)}>✕</button>
            </div>
            <form onSubmit={handleChangePassword} className="change-password-form">
              <div className="form-group">
                <label>Текущий пароль</label>
                <input
                  type="password"
                  name="old_password"
                  placeholder="Введите текущий пароль"
                  required
                />
              </div>
              <div className="form-group">
                <label>Новый пароль</label>
                <input
                  type="password"
                  name="new_password"
                  placeholder="Введите новый пароль"
                  required
                />
              </div>
              <div className="form-group">
                <label>Подтверждение пароля</label>
                <input
                  type="password"
                  name="confirm_password"
                  placeholder="Повторите новый пароль"
                  required
                />
              </div>
              <button type="submit" className="btn-primary">Сменить пароль</button>
            </form>
          </div>
        </div>
      )}

      <div className="add-form">
        <h3>Добавить новую ссылку</h3>
        <form onSubmit={handleAdd}>
          <div className="form-group">
            <input
              type="text"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              placeholder="Название ссылки (например: VPN Server 1)..."
            />
          </div>
          <div className="form-group">
            <textarea
              value={newLink}
              onChange={(e) => setNewLink(e.target.value)}
              placeholder="Вставьте VPN ссылку..."
            />
          </div>
          <button type="submit">Добавить ссылку</button>
        </form>
      </div>

      <div className="vpn-list">
        <h3>Список ссылок ({links.length})</h3>

        {links.length === 0 ? (
          <p style={{ textAlign: 'center', color: '#999', padding: '20px' }}>
            Ссылки отсутствуют
          </p>
        ) : (
          links.map((link) => (
            <div key={link.id} className="vpn-item">
              {editingId === link.id ? (
                <>
                  <div style={{ flex: 1 }}>
                    <input
                      type="text"
                      value={editingName}
                      onChange={(e) => setEditingName(e.target.value)}
                      placeholder="Название ссылки"
                      style={{ width: '100%', marginBottom: '8px', padding: '8px', border: '1px solid #ddd', borderRadius: '4px' }}
                    />
                    <textarea
                      value={editingLink}
                      onChange={(e) => setEditingLink(e.target.value)}
                      style={{ width: '100%', padding: '8px', border: '1px solid #ddd', borderRadius: '4px', minHeight: '60px' }}
                    />
                  </div>
                  <div className="vpn-actions">
                    <button
                      className="btn-small btn-edit"
                      onClick={() => handleUpdate(link.id)}
                    >
                      Сохранить
                    </button>
                    <button
                      className="btn-small"
                      onClick={cancelEdit}
                      style={{ background: '#999' }}
                    >
                      Отмена
                    </button>
                  </div>
                </>
              ) : (
                <>
                  <div style={{ flex: 1 }}>
                    <div style={{ fontWeight: 'bold', marginBottom: '8px', color: '#333' }}>{link.name}</div>
                    <div className="vpn-link">{link.link}</div>
                  </div>
                  <div className="vpn-actions">
                    <button
                      className={`btn-small ${copiedId === link.link ? 'btn-copied' : 'btn-copy'}`}
                      onClick={() => handleCopyLink(link.link)}
                    >
                      {copiedId === link.link ? '✓ Скопировано' : 'Копировать'}
                    </button>
                    <button
                      className="btn-small btn-edit"
                      onClick={() => startEdit(link)}
                    >
                      Изменить
                    </button>
                    <button
                      className="btn-small btn-delete"
                      onClick={() => handleDelete(link.id)}
                    >
                      Удалить
                    </button>
                  </div>
                </>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default Dashboard
