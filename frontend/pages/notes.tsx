import { useState, useEffect, FormEvent } from 'react'
import { useRouter } from 'next/router'

const API = process.env.NEXT_PUBLIC_SCRYPTS_API || 'http://localhost:8080'

interface Note {
  id: string
  owner: string
  content: string
  created: number
  modified: number
}

export default function Notes() {
  const [notes, setNotes] = useState<Note[]>([])
  const [content, setContent] = useState('')
  const [editId, setEditId] = useState<string | null>(null)
  const [msg, setMsg] = useState('')
  const router = useRouter()

  const token = typeof window !== 'undefined' ? localStorage.getItem('scrypts_token') : null
  const user = typeof window !== 'undefined' ? localStorage.getItem('scrypts_user') : null

  useEffect(() => {
    if (!token) {
      router.push('/')
      return
    }
    fetchNotes()
  }, [])

  async function fetchNotes() {
    const res = await fetch(`${API}/notes`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      setNotes(data || [])
    } else {
      setMsg('Failed to fetch notes')
    }
  }

  async function createNote(e: FormEvent) {
    e.preventDefault()
    if (!content.trim()) return
    const res = await fetch(`${API}/notes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ content })
    })
    if (res.status === 201) {
      setContent('')
      setMsg('Note created')
      fetchNotes()
    } else {
      const txt = await res.text()
      setMsg(`Create failed: ${txt}`)
    }
  }

  async function updateNote(e: FormEvent) {
    e.preventDefault()
    if (!editId || !content.trim()) return
    const res = await fetch(`${API}/notes`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ id: editId, content })
    })
    if (res.ok) {
      setContent('')
      setEditId(null)
      setMsg('Note updated')
      fetchNotes()
    } else {
      const txt = await res.text()
      setMsg(`Update failed: ${txt}`)
    }
  }

  async function deleteNote(id: string) {
    if (!confirm('Delete this note?')) return
    const res = await fetch(`${API}/notes`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ id })
    })
    if (res.ok) {
      setMsg('Note deleted')
      fetchNotes()
    } else {
      const txt = await res.text()
      setMsg(`Delete failed: ${txt}`)
    }
  }

  function startEdit(note: Note) {
    setEditId(note.id)
    setContent(note.content)
  }

  function cancelEdit() {
    setEditId(null)
    setContent('')
  }

  function logout() {
    localStorage.removeItem('scrypts_token')
    localStorage.removeItem('scrypts_user')
    router.push('/')
  }

  return (
    <div style={{ maxWidth: 800, margin: '2rem auto', padding: '0 1rem', fontFamily: 'system-ui, sans-serif' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1>My Notes — {user}</h1>
        <button
          onClick={logout}
          style={{ padding: '0.5rem 1rem', cursor: 'pointer', background: '#dc2626', color: '#fff', border: 'none', borderRadius: 4 }}
        >
          Logout
        </button>
      </div>

      {msg && <p style={{ marginBottom: '1rem', color: msg.includes('failed') ? '#d00' : '#080' }}>{msg}</p>}

      <form onSubmit={editId ? updateNote : createNote} style={{ marginBottom: '2rem' }}>
        <textarea
          placeholder="Write your note here..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={4}
          style={{ width: '100%', padding: '0.5rem', border: '1px solid #ccc', borderRadius: 4, fontFamily: 'inherit', fontSize: '1rem' }}
        />
        <div style={{ display: 'flex', gap: 8, marginTop: '0.5rem' }}>
          <button
            type="submit"
            style={{ flex: 1, padding: '0.5rem', cursor: 'pointer', background: editId ? '#f59e0b' : '#0070f3', color: '#fff', border: 'none', borderRadius: 4 }}
          >
            {editId ? 'Update Note' : 'Create Note'}
          </button>
          {editId && (
            <button
              onClick={cancelEdit}
              type="button"
              style={{ flex: 1, padding: '0.5rem', cursor: 'pointer', background: '#666', color: '#fff', border: 'none', borderRadius: 4 }}
            >
              Cancel
            </button>
          )}
        </div>
      </form>

      <div style={{ display: 'grid', gap: '1rem' }}>
        {notes.length === 0 && <p style={{ color: '#666' }}>No notes yet. Create your first note above!</p>}
        {notes.map((note) => (
          <div
            key={note.id}
            style={{
              padding: '1rem',
              border: '1px solid #e5e7eb',
              borderRadius: 8,
              background: '#f9fafb'
            }}
          >
            <p style={{ whiteSpace: 'pre-wrap', marginBottom: '0.5rem' }}>{note.content}</p>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', fontSize: '0.75rem', color: '#666' }}>
              <span>
                Modified: {new Date(note.modified * 1000).toLocaleString()}
              </span>
              <div style={{ display: 'flex', gap: 8 }}>
                <button
                  onClick={() => startEdit(note)}
                  style={{ padding: '0.25rem 0.5rem', cursor: 'pointer', background: '#f59e0b', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.75rem' }}
                >
                  Edit
                </button>
                <button
                  onClick={() => deleteNote(note.id)}
                  style={{ padding: '0.25rem 0.5rem', cursor: 'pointer', background: '#dc2626', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.75rem' }}
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
