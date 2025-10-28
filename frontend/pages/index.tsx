import { useState, FormEvent } from 'react'
import { useRouter } from 'next/router'

const API = process.env.NEXT_PUBLIC_SCRYPTS_API || 'http://localhost:8080'

export default function Home() {
  const [user, setUser] = useState('')
  const [pass, setPass] = useState('')
  const [msg, setMsg] = useState('')
  const router = useRouter()

  async function register(e: FormEvent) {
    e.preventDefault()
    const res = await fetch(`${API}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, password: pass })
    })
    if (res.status === 201) {
      setMsg('Registered — now login')
    } else {
      const txt = await res.text()
      setMsg(`Register failed: ${txt}`)
    }
  }

  async function login(e: FormEvent) {
    e.preventDefault()
    const res = await fetch(`${API}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, password: pass })
    })
    if (res.ok) {
      const data = await res.json()
      localStorage.setItem('scrypts_token', data.token)
      localStorage.setItem('scrypts_user', user)
      router.push('/notes')
    } else {
      const txt = await res.text()
      setMsg(`Login failed: ${txt}`)
    }
  }

  return (
    <div style={{ maxWidth: 640, margin: '4rem auto', fontFamily: 'system-ui, sans-serif' }}>
      <h1>Scrypts — Encrypted Notes</h1>
      <p style={{ margin: '1rem 0', color: msg.includes('failed') ? '#d00' : '#080' }}>{msg}</p>
      <form onSubmit={login} style={{ display: 'grid', gap: 8 }}>
        <input
          placeholder="username"
          value={user}
          onChange={(e) => setUser(e.target.value)}
          style={{ padding: '0.5rem', border: '1px solid #ccc', borderRadius: 4 }}
        />
        <input
          placeholder="password"
          value={pass}
          type="password"
          onChange={(e) => setPass(e.target.value)}
          style={{ padding: '0.5rem', border: '1px solid #ccc', borderRadius: 4 }}
        />
        <div style={{ display: 'flex', gap: 8 }}>
          <button
            onClick={login}
            type="button"
            style={{ flex: 1, padding: '0.5rem', cursor: 'pointer', background: '#0070f3', color: '#fff', border: 'none', borderRadius: 4 }}
          >
            Login
          </button>
          <button
            onClick={register}
            type="button"
            style={{ flex: 1, padding: '0.5rem', cursor: 'pointer', background: '#10b981', color: '#fff', border: 'none', borderRadius: 4 }}
          >
            Register
          </button>
          <button
            onClick={() => router.push('/notes')}
            type="button"
            style={{ flex: 1, padding: '0.5rem', cursor: 'pointer', background: '#666', color: '#fff', border: 'none', borderRadius: 4 }}
          >
            Notes
          </button>
        </div>
      </form>
      <p style={{ marginTop: 24, color: '#666', fontSize: '0.875rem' }}>
        This is a minimal client for local dev. Set <code>NEXT_PUBLIC_SCRYPTS_API</code> to point to your API if needed.
      </p>
    </div>
  )
}
