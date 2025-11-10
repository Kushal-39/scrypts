'use client'

import { useRouter } from 'next/navigation'
import { useAuthStore, useNotesStore } from '@/lib/store'
import { NoteEditor } from '@/components/NoteEditor'
import { useToast } from '@/components/Toast'
import { useEffect } from 'react'
import { ArrowLeft } from 'lucide-react'

export default function NewNotePage() {
  const router = useRouter()
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const createNote = useNotesStore((state) => state.createNote)
  const { showToast, Toast } = useToast()

  useEffect(() => {
    if (!isAuthenticated()) {
      router.replace('/login')
    }
  }, [isAuthenticated, router])

  const handleSaveNote = async (title: string, content: string) => {
    try {
      await createNote(title, content)
      showToast('Note created securely', 'success')
      router.push('/dashboard')
    } catch (error: any) {
      showToast(error.message, 'error')
    }
  }

  const handleBack = () => {
    router.push('/dashboard')
  }

  if (!isAuthenticated()) {
    return null
  }

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Header */}
      <div className="flex items-center gap-4 p-4 border-b border-slate-700/50 bg-slate-800/30 backdrop-blur-sm">
        <button
          onClick={handleBack}
          className="flex items-center gap-2 px-3 py-1.5 text-sm font-mono text-slate-400 hover:text-accent transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to Dashboard
        </button>
        
        <h1 className="text-lg font-mono font-medium text-foreground">
          Create New Note
        </h1>
      </div>

      {/* Editor */}
      <div className="flex-1">
        <NoteEditor
          note={null}
          onSave={handleSaveNote}
          className="h-full"
        />
      </div>

      {Toast}
    </div>
  )
}