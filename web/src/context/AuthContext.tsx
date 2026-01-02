import { createContext, useContext, useState, useEffect } from 'react'
import type { User } from '@/types'
import { authApi } from '@/lib/api'

interface AuthContextValue {
  user: User | null
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    try {
      const response = await authApi.me()
      if (response.success && response.data) {
        setUser(response.data)
      }
    } catch {
      // Not authenticated
    } finally {
      setLoading(false)
    }
  }

  const login = async (email: string, password: string) => {
    const response = await authApi.login({ email, password })
    if (response.success && response.data) {
      setUser(response.data)
    } else {
      throw new Error(response.error || 'Login failed')
    }
  }

  const register = async (email: string, password: string, name: string) => {
    const response = await authApi.register({ email, password, name })
    if (response.success && response.data) {
      // Auto-login after registration
      await login(email, password)
    } else {
      throw new Error(response.error || 'Registration failed')
    }
  }

  const logout = async () => {
    await authApi.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
