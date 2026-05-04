import React, { createContext, useContext, useState, useEffect } from 'react'
import request from '../utils/request'

const AuthContext = createContext()

export const AuthProvider = ({ children }) => {
  const [admin, setAdmin] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('admin_token')
    const savedAdmin = localStorage.getItem('admin')
    
    if (token && savedAdmin) {
      setAdmin(JSON.parse(savedAdmin))
    }
    setLoading(false)
  }, [])

  const login = async (username, password) => {
    const res = await request.post('/admin/login', { username, password })
    const { token, admin } = res.data
    localStorage.setItem('admin_token', token)
    localStorage.setItem('admin', JSON.stringify(admin))
    setAdmin(admin)
    return res
  }

  const logout = () => {
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin')
    setAdmin(null)
  }

  return (
    <AuthContext.Provider value={{
      admin,
      loading,
      login,
      logout,
      isAuthenticated: !!admin,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
