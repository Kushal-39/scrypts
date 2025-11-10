'use client'

import { Note, useNotesStore } from '@/lib/store'
import { formatDistanceToNow } from 'date-fns'

interface NoteCardProps {
  note: Note
  isActive?: boolean
  onClick: () => void
}

export function NoteCard({ note, isActive, onClick }: NoteCardProps) {
  const getDisplayNote = useNotesStore((state) => state.getDisplayNote)
  const displayNote = getDisplayNote(note)
  
  const preview = displayNote.content.slice(0, 100) + (displayNote.content.length > 100 ? '...' : '')
  const timeAgo = formatDistanceToNow(new Date(displayNote.updatedAt), { addSuffix: true })

  return (
    <button
      onClick={onClick}
      className={`
        w-full text-left p-4 rounded-lg border transition-all hover:border-accent/50 group
        ${isActive 
          ? 'bg-accent/10 border-accent/50 shadow-lg' 
          : 'bg-slate-800/30 border-slate-700/50 hover:bg-slate-800/50'
        }
      `}
    >
      <div className="space-y-2">
        {/* Title */}
        <h3 className={`
          font-mono font-medium line-clamp-1 transition-colors
          ${isActive ? 'text-accent' : 'text-foreground group-hover:text-accent'}
        `}>
          {displayNote.title || 'Untitled'}
        </h3>
        
        {/* Preview */}
        <p className="text-sm text-slate-400 line-clamp-2 font-mono leading-relaxed">
          {preview || 'No content'}
        </p>
        
        {/* Timestamp */}
        <p className="text-xs text-slate-500 font-mono">
          {timeAgo}
        </p>
      </div>
    </button>
  )
}

// Install date-fns for time formatting
// npm install date-fns