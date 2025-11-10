'use client'

import { useRouter, useParams } from 'next/navigation'
import { useEffect, useState } from 'react'
import { ArrowLeft } from 'lucide-react'
import { useAuthStore, useNotesStore, Note } from '@/lib/store'
import { NoteEditor } from '@/components/NoteEditor'
import { Modal } from '@/components/Modal'
import { useToast } from '@/components/Toast'

export default function NotePage() {
  const router = useRouter()
  const params = useParams()
  const noteId = params.id as string
  
  const [note, setNote] = useState<Note | null>(null)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [loading, setLoading] = useState(true)
  
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const { notes, fetchNotes, updateNote, deleteNote } = useNotesStore()
  const { showToast, Toast } = useToast()

  useEffect(() => {
    if (!isAuthenticated()) {
      router.replace('/login')
      return
    }

    // Find note in store or fetch notes
    const findNote = () => {
      const foundNote = notes.find(n => n.id === noteId)
      if (foundNote) {
        setNote(foundNote)
        setLoading(false)
      } else if (notes.length === 0) {
        // Notes not loaded yet, fetch them
        fetchNotes().then(() => {
          const freshNote = notes.find(n => n.id === noteId)
          if (freshNote) {
            setNote(freshNote)
          } else {
            showToast('Note not found', 'error')
            router.push('/dashboard')
          }
          setLoading(false)
        }).catch(() => {
          showToast('Failed to load note', 'error')
          router.push('/dashboard')
          setLoading(false)
        })
      } else {
        // Notes loaded but note not found
        showToast('Note not found', 'error')
        router.push('/dashboard')
        setLoading(false)
      }
    }

    findNote()
  }, [noteId, notes, isAuthenticated, router, fetchNotes, showToast])

  const handleBack = () => {
    router.push('/dashboard')
  }

  const handleSaveNote = async (title: string, content: string) => {
    if (!note) return
    
    try {
      await updateNote(note.id, title, content)
      // Refetch notes to get updated state
      await fetchNotes()
      const updatedNote = notes.find(n => n.id === note.id)
      if (updatedNote) {
        setNote(updatedNote)
      }
      showToast('Note updated securely', 'success')
    } catch (error: any) {
      showToast(error.message, 'error')
    }
  }

  const handleDeleteNote = async () => {
    if (!note) return
    
    try {
      await deleteNote(note.id)
      showToast('Note deleted', 'success')
      router.push('/dashboard')
    } catch (error: any) {
      showToast(error.message, 'error')
    }
  }

  const confirmDelete = () => {
    setDeleteModalOpen(true)
  }

  if (!isAuthenticated()) {
    return null
  }

  if (loading) {
    return (
      <div className="h-screen flex items-center justify-center bg-background">
        <div className="text-accent font-mono">Loading note...</div>
      </div>
    )
  }

  if (!note) {
    return (
      <div className="h-screen flex items-center justify-center bg-background">
        <div className="text-center space-y-4">
          <div className="text-4xl">📝</div>
          <div>
            <h3 className="text-lg font-mono font-medium text-foreground mb-2">
              Note not found
            </h3>
            <p className="text-slate-400 font-mono text-sm mb-4">
              The note you're looking for doesn't exist or has been deleted.
            </p>
            <button
              onClick={handleBack}
              className="flex items-center gap-2 px-4 py-2 bg-accent hover:bg-accent-dark text-slate-900 font-mono font-medium rounded transition-colors mx-auto"
            >
              <ArrowLeft className="w-4 h-4" />
              Back to Dashboard
            </button>
          </div>
        </div>
      </div>
    )
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
          Edit Note
        </h1>
      </div>

      {/* Editor */}
      <div className="flex-1">
        <NoteEditor
          note={note}
          onSave={handleSaveNote}
          onDelete={async () => confirmDelete()}
          className="h-full"
        />
      </div>

      {/* Delete confirmation modal */}
      <Modal
        isOpen={deleteModalOpen}
        onClose={() => setDeleteModalOpen(false)}
        onConfirm={handleDeleteNote}
        title="Delete Note"
        confirmText="Delete"
        confirmVariant="danger"
      >
        <p className="text-foreground font-mono">
          Are you sure you want to delete "{note ? useNotesStore.getState().getDisplayNote(note).title || 'Untitled' : 'Untitled'}"?
          <br />
          <span className="text-slate-400 text-sm">This action cannot be undone.</span>
        </p>
      </Modal>

      {Toast}
    </div>
  )
}