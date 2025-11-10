'use client'

import { ReactNode } from 'react'
import { X } from 'lucide-react'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm?: () => void
  title: string
  children: ReactNode
  confirmText?: string
  confirmVariant?: 'danger' | 'primary'
}

export function Modal({ 
  isOpen, 
  onClose, 
  onConfirm, 
  title, 
  children, 
  confirmText = 'Confirm',
  confirmVariant = 'primary'
}: ModalProps) {
  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      />
      
      {/* Modal */}
      <div className="relative w-full max-w-md bg-slate-800 rounded-lg shadow-xl border border-slate-700 fade-in">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-slate-700">
          <h3 className="text-lg font-semibold font-mono text-foreground">{title}</h3>
          <button
            onClick={onClose}
            className="p-1 rounded hover:bg-slate-700 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
        
        {/* Content */}
        <div className="p-6">
          {children}
        </div>
        
        {/* Footer */}
        {onConfirm && (
          <div className="flex gap-3 p-6 pt-0">
            <button
              onClick={onClose}
              className="flex-1 px-4 py-2 text-sm font-mono font-medium text-slate-300 bg-slate-700 hover:bg-slate-600 rounded transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={() => {
                onConfirm()
                onClose()
              }}
              className={`
                flex-1 px-4 py-2 text-sm font-mono font-medium rounded transition-colors
                ${confirmVariant === 'danger'
                  ? 'text-white bg-red-600 hover:bg-red-700'
                  : 'text-slate-900 bg-accent hover:bg-accent-dark'
                }
              `}
            >
              {confirmText}
            </button>
          </div>
        )}
      </div>
    </div>
  )
}