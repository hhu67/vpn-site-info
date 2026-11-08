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
  const [generatingHappId, setGeneratingHappId] = useState(null)

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

  const handleGenerateCryptoLink = async (id) => {
    setError('')
    setSuccess('')
    setGeneratingHappId(id)

    try {
      const response = await fetch('/api/crypto-link', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
      })

      if (!response.ok) {
        let errMessage = 'Ошибка генерации Happ ссылки'
        try {
          const errData = await response.json()
          if (errData && errData.error) errMessage = errData.error
        } catch {
          const text = await response.text()
          if (text) errMessage = text
        }
        throw new Error(errMessage)
      }

      const data = await response.json()
      const cryptoLink = data.crypto_link || data.link

      setLinks((prevLinks) =>
        prevLinks.map((item) =>
          item.id === id ? { ...item, crypto_link: cryptoLink } : item
        )
      )

      setSuccess('Happ крипто-ссылка успешно создана!')
    } catch (err) {
      setError(err.message)
    } finally {
      setGeneratingHappId(null)
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

  const handleCopyLink = (text, identifier) => {
    navigator.clipboard.writeText(text)
    setCopiedId(identifier || text)
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
                <div className="vpn-edit-mode">
                  <div className="form-group" style={{ marginBottom: '10px' }}>
                    <label>Название</label>
                    <input
                      type="text"
                      value={editingName}
                      onChange={(e) => setEditingName(e.target.value)}
                      placeholder="Название ссылки"
                    />
                  </div>
                  <div className="form-group" style={{ marginBottom: '12px' }}>
                    <label>VPN ссылка</label>
                    <textarea
                      value={editingLink}
                      onChange={(e) => setEditingLink(e.target.value)}
                      placeholder="VPN ссылка"
                      style={{ minHeight: '70px' }}
                    />
                  </div>
                  <div className="vpn-actions">
                    <button
                      className="btn-small btn-primary"
                      onClick={() => handleUpdate(link.id)}
                    >
                      Сохранить
                    </button>
                    <button
                      className="btn-small"
                      onClick={cancelEdit}
                      style={{ background: '#64748b' }}
                    >
                      Отмена
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <div className="vpn-item-header">
                    <div className="vpn-item-title">{link.name}</div>
                    <div className="vpn-actions">
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
                  </div>

                  {/* Regular VPN Link */}
                  <div className="link-block">
                    <div className="link-block-label">VPN ссылка</div>
                    <div className="link-block-content">
                      <a
                        href={link.link}
                        target={link.link.startsWith('http') ? '_blank' : '_self'}
                        rel="noopener noreferrer"
                        className="vpn-link"
                        title={`Нажмите для перехода\n${link.link}`}
                      >
                        {link.link}
                      </a>
                      <div className="link-actions">
                        <a
                          href={link.link}
                          target={link.link.startsWith('http') ? '_blank' : '_self'}
                          rel="noopener noreferrer"
                          className="btn-small btn-open"
                          title="Перейти по ссылке"
                        >
                          ↗ Перейти
                        </a>
                        <button
                          className={`btn-small ${copiedId === link.id + '-vpn' ? 'btn-copied' : 'btn-copy'}`}
                          onClick={() => handleCopyLink(link.link, link.id + '-vpn')}
                        >
                          {copiedId === link.id + '-vpn' ? '✓ Скопировано' : 'Копировать'}
                        </button>
                      </div>
                    </div>
                  </div>

                  {/* Happ Crypto Link */}
                  <div className="link-block">
                    <div className="link-block-label">
                      <span className="happ-badge">Happ</span>
                      <span>Крипто-ссылка (happ://)</span>
                    </div>
                    {link.crypto_link ? (
                      <div className="link-block-content link-block-happ">
                        <a
                          href={link.crypto_link}
                          target={link.crypto_link.startsWith('http') ? '_blank' : '_self'}
                          rel="noopener noreferrer"
                          className="vpn-link happ-link-text"
                          title={`Нажмите для перехода в Happ\n${link.crypto_link}`}
                        >
                          {link.crypto_link}
                        </a>
                        <div className="link-actions">
                          <a
                            href={link.crypto_link}
                            target={link.crypto_link.startsWith('http') ? '_blank' : '_self'}
                            rel="noopener noreferrer"
                            className="btn-small btn-open"
                            title="Открыть в Happ"
                          >
                            ↗ Перейти
                          </a>
                          <button
                            className={`btn-small ${copiedId === link.id + '-happ' ? 'btn-copied' : 'btn-happ'}`}
                            onClick={() => handleCopyLink(link.crypto_link, link.id + '-happ')}
                          >
                            {copiedId === link.id + '-happ' ? '✓ Скопировано' : 'Копировать Happ'}
                          </button>
                          <button
                            className="btn-small btn-happ-refresh"
                            onClick={() => handleGenerateCryptoLink(link.id)}
                            disabled={generatingHappId === link.id}
                            title="Обновить крипто-ссылку Happ"
                          >
                            {generatingHappId === link.id ? '...' : '🔄'}
                          </button>
                        </div>
                      </div>
                    ) : (
                      <div className="happ-empty-state">
                        <span>Крипто-ссылка Happ еще не создана</span>
                        <button
                          className="btn-small btn-happ-generate"
                          onClick={() => handleGenerateCryptoLink(link.id)}
                          disabled={generatingHappId === link.id}
                        >
                          {generatingHappId === link.id ? '⚡ Генерация Happ...' : '⚡ Создать Happ ссылку'}
                        </button>
                      </div>
                    )}
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
