import { create } from 'zustand'
import axios from 'axios'

export interface Note {
  id: string
  owner: string
  content: string
  created: number
  modified: number
}

export interface NoteDisplay {
  id: string
  title: string
  content: string
  createdAt: string
  updatedAt: string
}

interface AuthState {
  token: string | null
  user: string | null
  login: (username: string, password: string) => Promise<void>
  register: (username: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: () => boolean
}

interface NotesState {
  notes: Note[]
  currentNote: Note | null
  loading: boolean
  fetchNotes: () => Promise<void>
  createNote: (title: string, content: string) => Promise<void>
  updateNote: (id: string, title: string, content: string) => Promise<void>
  deleteNote: (id: string) => Promise<void>
  setCurrentNote: (note: Note | null) => void
  getDisplayNote: (note: Note) => NoteDisplay
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'

// Configure axios defaults
axios.defaults.baseURL = API_BASE

// Axios interceptor to add token to requests
axios.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Axios interceptor to handle 401 responses
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
    }
    return Promise.reject(error)
  }
)

export const useAuthStore = create<AuthState>((set, get) => ({
  token: null,
  user: null,
  
  login: async (username: string, password: string) => {
    try {
      const response = await axios.post('/login', { username, password })
      const { token } = response.data
      set({ token, user: username })
    } catch (error: any) {
      throw new Error(error.response?.data || 'Login failed')
    }
  },
  
  register: async (username: string, password: string) => {
    try {
      await axios.post('/register', { username, password })
    } catch (error: any) {
      throw new Error(error.response?.data || 'Registration failed')
    }
  },
  
  logout: () => {
    set({ token: null, user: null })
    useNotesStore.getState().notes = []
    useNotesStore.getState().currentNote = null
  },
  
  isAuthenticated: () => !!get().token,
}))

export const useNotesStore = create<NotesState>((set, get) => ({
  notes: [],
  currentNote: null,
  loading: false,
  
  fetchNotes: async () => {
    set({ loading: true })
    try {
      const response = await axios.get('/notes')
      set({ notes: response.data, loading: false })
    } catch (error) {
      set({ loading: false })
      throw error
    }
  },
  
  createNote: async (title: string, content: string) => {
    try {
      // Combine title and content for backend
      const fullContent = title ? `${title}\n${content}` : content
      const response = await axios.post('/notes', { content: fullContent })
      // Backend returns {id: "uuid"}, we need to fetch notes to get the full note
      await get().fetchNotes()
    } catch (error) {
      throw error
    }
  },
  
  updateNote: async (id: string, title: string, content: string) => {
    try {
      // Combine title and content for backend
      const fullContent = title ? `${title}\n${content}` : content
      await axios.put('/notes', { id, content: fullContent })
      // Backend returns {status: "updated"}, we need to fetch notes to get updated data
      await get().fetchNotes()
      // Update current note if it was the one being edited
      const notes = get().notes
      const updatedNote = notes.find(note => note.id === id)
      if (updatedNote && get().currentNote?.id === id) {
        set({ currentNote: updatedNote })
      }
    } catch (error) {
      throw error
    }
  },
  
  deleteNote: async (id: string) => {
    try {
      await axios.delete('/notes', { data: { id } })
      set({ 
        notes: get().notes.filter(note => note.id !== id),
        currentNote: get().currentNote?.id === id ? null : get().currentNote
      })
    } catch (error) {
      throw error
    }
  },
  
  setCurrentNote: (note: Note | null) => {
    set({ currentNote: note })
  },

  getDisplayNote: (note: Note): NoteDisplay => {
    // Parse title from content (first line)
    const lines = note.content.split('\n')
    const title = lines[0] || 'Untitled'
    const content = lines.slice(1).join('\n')
    
    return {
      id: note.id,
      title,
      content,
      createdAt: new Date(note.created * 1000).toISOString(),
      updatedAt: new Date(note.modified * 1000).toISOString()
    }
  },
}))