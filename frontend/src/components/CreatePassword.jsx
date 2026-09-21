import { useState } from 'react'

function CreatePassword({ onSuccess }) {
  const [password, setPassword] = useState('')
  const [passwordConfirm, setPasswordConfirm] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    if (!password || !passwordConfirm) {
      setError('Заполните все поля')
      return
    }

    if (password !== passwordConfirm) {
      setError('Пароли не совпадают')
      return
    }

    setLoading(true)

    try {
      const response = await fetch('/api/create-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password, password_confirm: passwordConfirm }),
      })

      if (!response.ok) {
        const data = await response.text()
        throw new Error(data || 'Ошибка создания пароля')
      }

      onSuccess()
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container">
      <h1>VPN Управление</h1>
      <h2>Создание пароля</h2>

      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Пароль</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={loading}
            autoFocus
          />
        </div>

        <div className="form-group">
          <label>Подтверждение пароля</label>
          <input
            type="password"
            value={passwordConfirm}
            onChange={(e) => setPasswordConfirm(e.target.value)}
            disabled={loading}
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? 'Создание...' : 'Создать пароль'}
        </button>
      </form>
    </div>
  )
}

export default CreatePassword
