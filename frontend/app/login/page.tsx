'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { AuthForm } from '@/components/AuthForm'
import { useToast } from '@/components/Toast'
import { useAuthStore } from '@/lib/store'

export default function LoginPage() {
  const [mode, setMode] = useState<'login' | 'register'>('login')
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
        setMode('login')
      }
    } catch (error: any) {
      showToast(error.message, 'error')
    } finally {
      setLoading(false)
    }
  }

  const handleModeChange = () => {
    if (mode === 'login') {
      router.push('/register')
    } else {
      setMode('login')
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