import { useState, useEffect } from 'react'
import CreatePassword from './components/CreatePassword'
import Login from './components/Login'
import Dashboard from './components/Dashboard'

function App() {
  const [loading, setLoading] = useState(true)
  const [passwordExists, setPasswordExists] = useState(false)
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    try {
      const response = await fetch('/api/check-password-exists')
      const data = await response.json()
      setPasswordExists(data.exists)

      if (data.exists) {
        const vpnResponse = await fetch('/api/list/vpn')
        if (vpnResponse.ok) {
          setIsAuthenticated(true)
        }
      }
    } catch (error) {
      console.error('Auth check failed:', error)
    } finally {
      setLoading(false)
    }
  }

  const handlePasswordCreated = () => {
    setPasswordExists(true)
    setIsAuthenticated(true)
  }

  const handleLoginSuccess = () => {
    setIsAuthenticated(true)
  }

  const handleLogout = () => {
    document.cookie = 'jwt_token=; Path=/; Max-Age=0'
    setIsAuthenticated(false)
  }

  if (loading) {
    return (
      <div className="container">
        <div className="loading">Загрузка...</div>
      </div>
    )
  }

  if (!passwordExists) {
    return <CreatePassword onSuccess={handlePasswordCreated} />
  }

  if (!isAuthenticated) {
    return <Login onSuccess={handleLoginSuccess} />
  }

  return <Dashboard onLogout={handleLogout} />
}

export default App
