'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'

export default function Home() {
  const router = useRouter()
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  useEffect(() => {
    if (isAuthenticated()) {
      router.replace('/dashboard')
    } else {
      router.replace('/login')
    }
  }, [router, isAuthenticated])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-accent font-mono">Loading...</div>
    </div>
  )
}
