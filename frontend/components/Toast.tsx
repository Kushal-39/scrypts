'use client'

import { useState, useEffect } from 'react'
import { CheckCircle, XCircle, X } from 'lucide-react'

interface ToastProps {
  message: string
  type: 'success' | 'error'
  onClose: () => void
}

export function Toast({ message, type, onClose }: ToastProps) {
  useEffect(() => {
    const timer = setTimeout(onClose, 4000)
    return () => clearTimeout(timer)
  }, [onClose])

  return (
    <div className="fixed top-4 right-4 z-50 fade-in">
      <div className={`
        flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg backdrop-blur-sm
        ${type === 'success' 
          ? 'bg-green-900/80 border border-green-700/50 text-green-100' 
          : 'bg-red-900/80 border border-red-700/50 text-red-100'
        }
      `}>
        {type === 'success' ? (
          <CheckCircle className="w-5 h-5 text-green-400" />
        ) : (
          <XCircle className="w-5 h-5 text-red-400" />
        )}
        <span className="text-sm font-mono">{message}</span>
        <button
          onClick={onClose}
          className="ml-2 p-1 rounded hover:bg-white/10 transition-colors"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
    </div>
  )
}

// Toast manager hook
export function useToast() {
  const [toast, setToast] = useState<{
    message: string
    type: 'success' | 'error'
  } | null>(null)

  const showToast = (message: string, type: 'success' | 'error') => {
    setToast({ message, type })
  }

  const hideToast = () => {
    setToast(null)
  }

  const ToastComponent = toast ? (
    <Toast
      message={toast.message}
      type={toast.type}
      onClose={hideToast}
    />
  ) : null

  return {
    showToast,
    hideToast,
    Toast: ToastComponent,
  }
}