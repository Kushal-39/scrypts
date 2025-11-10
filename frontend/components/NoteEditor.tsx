'use client'

import { useState, useEffect, useRef } from 'react'
import { Save, Trash2, CheckCircle } from 'lucide-react'
import { Note, useNotesStore } from '@/lib/store'

interface NoteEditorProps {
  note: Note | null
  onSave: (title: string, content: string) => Promise<void>
  onDelete?: (id: string) => Promise<void>
  loading?: boolean
  className?: string
}

export function NoteEditor({ note, onSave, onDelete, loading, className }: NoteEditorProps) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [hasChanges, setHasChanges] = useState(false)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const titleRef = useRef<HTMLInputElement>(null)
  const contentRef = useRef<HTMLTextAreaElement>(null)
  const getDisplayNote = useNotesStore((state) => state.getDisplayNote)

  // Update form when note changes
  useEffect(() => {
    if (note) {
      const displayNote = getDisplayNote(note)
      setTitle(displayNote.title)
      setContent(displayNote.content)
      setHasChanges(false)
    } else {
      setTitle('')
      setContent('')
      setHasChanges(false)
    }
  }, [note, getDisplayNote])

  // Track changes
  useEffect(() => {
    if (note) {
      const displayNote = getDisplayNote(note)
      setHasChanges(title !== displayNote.title || content !== displayNote.content)
    } else {
      setHasChanges(title.trim() !== '' || content.trim() !== '')
    }
  }, [title, content, note, getDisplayNote])

  // Auto-save on blur
  const handleAutoSave = async () => {
    if (hasChanges && (title.trim() || content.trim())) {
      await handleSave()
    }
  }

  const handleSave = async () => {
    try {
      setSaving(true)
      await onSave(title, content)
      setSaved(true)
      setHasChanges(false)
      setTimeout(() => setSaved(false), 2000)
    } catch (error) {
      console.error('Failed to save:', error)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (note && onDelete) {
      try {
        await onDelete(note.id)
      } catch (error) {
        console.error('Failed to delete:', error)
      }
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    // Cmd/Ctrl + S to save
    if ((e.metaKey || e.ctrlKey) && e.key === 's') {
      e.preventDefault()
      if (hasChanges) {
        handleSave()
      }
    }
  }

  return (
    <div className={`flex flex-col h-full ${className}`} onKeyDown={handleKeyDown}>
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-slate-700/50">
        <div className="flex items-center gap-3">
          <h2 className="font-mono text-lg font-medium text-foreground">
            {note ? 'Edit Note' : 'New Note'}
          </h2>
          {saved && (
            <div className="flex items-center gap-2 text-green-400">
              <CheckCircle className="w-4 h-4" />
              <span className="text-sm font-mono">Saved securely</span>
            </div>
          )}
        </div>
        
        <div className="flex items-center gap-2">
          {/* Save Button */}
          <button
            onClick={handleSave}
            disabled={!hasChanges || saving || loading}
            className="flex items-center gap-2 px-3 py-1.5 text-sm font-mono font-medium bg-accent hover:bg-accent-dark disabled:bg-slate-700 disabled:text-slate-400 text-slate-900 rounded transition-colors disabled:cursor-not-allowed"
          >
            <Save className="w-4 h-4" />
            {saving ? 'Saving...' : 'Save'}
          </button>

          {/* Delete Button */}
          {note && onDelete && (
            <button
              onClick={handleDelete}
              disabled={loading}
              className="flex items-center gap-2 px-3 py-1.5 text-sm font-mono font-medium bg-red-600 hover:bg-red-700 disabled:bg-slate-700 disabled:text-slate-400 text-white rounded transition-colors disabled:cursor-not-allowed"
            >
              <Trash2 className="w-4 h-4" />
              Delete
            </button>
          )}
        </div>
      </div>

      {/* Editor */}
      <div className="flex-1 flex flex-col p-4 space-y-4">
        {/* Title Input */}
        <input
          ref={titleRef}
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          onBlur={handleAutoSave}
          placeholder="Note title..."
          disabled={loading}
          className="w-full text-xl font-mono font-medium bg-transparent border-none outline-none text-foreground placeholder-slate-500 resize-none disabled:opacity-50"
        />

        {/* Content Textarea */}
        <textarea
          ref={contentRef}
          value={content}
          onChange={(e) => setContent(e.target.value)}
          onBlur={handleAutoSave}
          placeholder="Start writing..."
          disabled={loading}
          className="flex-1 w-full bg-transparent border-none outline-none text-foreground placeholder-slate-500 font-mono leading-relaxed resize-none disabled:opacity-50"
        />
      </div>

      {/* Status Bar */}
      <div className="flex items-center justify-between p-4 border-t border-slate-700/50 text-xs font-mono text-slate-500">
        <div className="flex items-center gap-4">
          <span>{content.length} characters</span>
          <span>{content.split(/\s+/).filter(word => word.length > 0).length} words</span>
        </div>
        
        <div className="flex items-center gap-2">
          {hasChanges && (
            <span className="text-yellow-400">Unsaved changes</span>
          )}
          <span>Cmd+S to save</span>
        </div>
      </div>
    </div>
  )
}