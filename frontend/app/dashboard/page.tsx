'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Plus, LogOut, Search, Menu, X } from 'lucide-react'
import { useAuthStore, useNotesStore, Note } from '@/lib/store'
import { NoteCard } from '@/components/NoteCard'
import { NoteEditor } from '@/components/NoteEditor'
import { Modal } from '@/components/Modal'
import { useToast } from '@/components/Toast'

export default function DashboardPage() {
  const router = useRouter()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [noteToDelete, setNoteToDelete] = useState<Note | null>(null)
  
  // Auth store
  const { user, logout, isAuthenticated } = useAuthStore()
  
  // Notes store
  const {
    notes,
    currentNote,
    loading,
    fetchNotes,
    createNote,
    updateNote,
    deleteNote,
    setCurrentNote,
  } = useNotesStore()
  
  const { showToast, Toast } = useToast()

  // Check authentication and fetch notes
  useEffect(() => {
    if (!isAuthenticated()) {
      router.replace('/login')
      return
    }
    
    fetchNotes().catch((error) => {
      showToast('Failed to load notes', 'error')
    })
  }, [isAuthenticated, router, fetchNotes, showToast])

  // Filter notes by search term
  const filteredNotes = notes.filter(note => {
    const displayNote = useNotesStore.getState().getDisplayNote(note)
    return (
      displayNote.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      displayNote.content.toLowerCase().includes(searchTerm.toLowerCase())
    )
  })

  const handleLogout = () => {
    logout()
    router.push('/login')
    showToast('Signed out successfully', 'success')
  }

  const handleNewNote = () => {
    setCurrentNote(null)
    setSidebarOpen(false)
  }

  const handleNoteSelect = (note: Note) => {
    setCurrentNote(note)
    setSidebarOpen(false)
  }

  const handleSaveNote = async (title: string, content: string) => {
    try {
      if (currentNote) {
        await updateNote(currentNote.id, title, content)
        showToast('Note updated securely', 'success')
      } else {
        await createNote(title, content)
        showToast('Note created securely', 'success')
      }
    } catch (error: any) {
      showToast(error.message, 'error')
    }
  }

  const handleDeleteNote = async (id: string) => {
    try {
      await deleteNote(id)
      showToast('Note deleted', 'success')
    } catch (error: any) {
      showToast(error.message, 'error')
    }
  }

  const confirmDelete = (note: Note) => {
    setNoteToDelete(note)
    setDeleteModalOpen(true)
  }

  if (!isAuthenticated()) {
    return null
  }

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Navbar */}
      <nav className="flex items-center justify-between p-4 border-b border-slate-700/50 bg-slate-800/30 backdrop-blur-sm">
        <div className="flex items-center gap-4">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="lg:hidden p-2 rounded hover:bg-slate-700/50 transition-colors"
          >
            <Menu className="w-5 h-5" />
          </button>
          
          <h1 className="text-xl font-bold font-mono text-accent glow">
            Scrypts
          </h1>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={handleNewNote}
            className="flex items-center gap-2 px-3 py-1.5 bg-accent hover:bg-accent-dark text-slate-900 font-mono font-medium text-sm rounded transition-colors"
          >
            <Plus className="w-4 h-4" />
            New Note
          </button>
          
          <div className="flex items-center gap-2 text-sm font-mono text-slate-400">
            <span>{user}</span>
            <button
              onClick={handleLogout}
              className="p-2 rounded hover:bg-slate-700/50 hover:text-accent transition-colors"
              title="Sign out"
            >
              <LogOut className="w-4 h-4" />
            </button>
          </div>
        </div>
      </nav>

      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar */}
        <aside className={`
          w-80 border-r border-slate-700/50 bg-slate-800/20 backdrop-blur-sm flex flex-col
          lg:relative lg:translate-x-0
          ${sidebarOpen ? 'fixed inset-y-0 left-0 z-40 translate-x-0' : 'fixed -translate-x-full lg:translate-x-0'}
          transition-transform duration-200 ease-in-out
        `}>
          <div className="p-4 border-b border-slate-700/50">
            <div className="flex items-center justify-between mb-3">
              <h2 className="font-mono font-medium text-foreground">Notes</h2>
              <button
                onClick={() => setSidebarOpen(false)}
                className="lg:hidden p-1 rounded hover:bg-slate-700/50 transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <input
                type="text"
                placeholder="Search notes..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-3 py-2 bg-slate-900/50 border border-slate-600 rounded text-sm font-mono text-foreground placeholder-slate-500 focus:border-accent focus:ring-1 focus:ring-accent transition-colors"
              />
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-2">
            {loading ? (
              <div className="text-center text-slate-400 font-mono text-sm py-8">
                Loading notes...
              </div>
            ) : filteredNotes.length === 0 ? (
              <div className="text-center text-slate-400 font-mono text-sm py-8">
                {searchTerm ? 'No notes found' : 'No notes yet'}
                <br />
                <button
                  onClick={handleNewNote}
                  className="text-accent hover:underline mt-2"
                >
                  Create your first note
                </button>
              </div>
            ) : (
              filteredNotes.map((note) => (
                <NoteCard
                  key={note.id}
                  note={note}
                  isActive={currentNote?.id === note.id}
                  onClick={() => handleNoteSelect(note)}
                />
              ))
            )}
          </div>
        </aside>

        {/* Sidebar overlay on mobile */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 bg-black/50 z-30 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
        )}

        {/* Main editor */}
        <main className="flex-1 flex flex-col">
          {currentNote || (!currentNote && notes.length === 0) ? (
            <NoteEditor
              note={currentNote}
              onSave={handleSaveNote}
              onDelete={currentNote ? async (id: string) => confirmDelete(currentNote) : undefined}
              loading={loading}
              className="h-full"
            />
          ) : (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center space-y-4">
                <div className="text-6xl">📝</div>
                <div>
                  <h3 className="text-lg font-mono font-medium text-foreground mb-2">
                    Select a note to edit
                  </h3>
                  <p className="text-slate-400 font-mono text-sm">
                    Choose a note from the sidebar or create a new one
                  </p>
                </div>
                <button
                  onClick={handleNewNote}
                  className="flex items-center gap-2 px-4 py-2 bg-accent hover:bg-accent-dark text-slate-900 font-mono font-medium rounded transition-colors mx-auto"
                >
                  <Plus className="w-4 h-4" />
                  New Note
                </button>
              </div>
            </div>
          )}
        </main>
      </div>

      {/* Delete confirmation modal */}
      <Modal
        isOpen={deleteModalOpen}
        onClose={() => setDeleteModalOpen(false)}
        onConfirm={() => {
          if (noteToDelete) {
            handleDeleteNote(noteToDelete.id)
          }
        }}
        title="Delete Note"
        confirmText="Delete"
        confirmVariant="danger"
      >
        <p className="text-foreground font-mono">
          Are you sure you want to delete "{noteToDelete ? useNotesStore.getState().getDisplayNote(noteToDelete).title || 'Untitled' : 'Untitled'}"?
          <br />
          <span className="text-slate-400 text-sm">This action cannot be undone.</span>
        </p>
      </Modal>

      {Toast}
    </div>
  )
}