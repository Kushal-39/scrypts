'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { AuthForm } from '@/components/AuthForm'
import { useToast } from '@/components/Toast'
import { useAuthStore } from '@/lib/store'

export default function RegisterPage() {
  const [mode, setMode] = useState<'login' | 'register'>('register')
  const [loading, setLoading] = useState(false)
  const router = useRouter()
  const login = useAuthStore((state) => state.login)
  const register = useAuthStore((state) => state.register)
  const { showToast, Toast } = useToast()

  const handleSubmit = async (username: string, password: string) => {
    setLoading(true)
    try {
      if (mode === 'login') {
        await login(username, password)
        showToast('Welcome back!', 'success')
        router.push('/dashboard')
      } else {
        await register(username, password)
        showToast('Account created! Please sign in.', 'success')
        router.push('/login')
      }
    } catch (error: any) {
      showToast(error.message, 'error')
    } finally {
      setLoading(false)
    }
  }

  const handleModeChange = () => {
    if (mode === 'register') {
      router.push('/login')
    } else {
      setMode('register')
    }
  }

  return (
    <>
      <AuthForm
        mode={mode}
        onSubmit={handleSubmit}
        onModeChange={handleModeChange}
        loading={loading}
      />
      {Toast}
    </>
  )
}