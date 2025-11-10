'use client'

import { useState } from 'react'
import { Eye, EyeOff, Lock, User } from 'lucide-react'

interface AuthFormProps {
  mode: 'login' | 'register'
  onSubmit: (username: string, password: string) => Promise<void>
  onModeChange: () => void
  loading: boolean
}

export function AuthForm({ mode, onSubmit, onModeChange, loading }: AuthFormProps) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (username.trim() && password.trim()) {
      await onSubmit(username.trim(), password)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-background">
      <div className="w-full max-w-md">
        <div className="bg-slate-800/50 backdrop-blur-sm rounded-lg border border-slate-700/50 p-8 glow">
          {/* Logo */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold font-mono text-accent glow">
              Scrypts
            </h1>
            <p className="text-slate-400 text-sm mt-2 font-mono">
              Your thoughts, encrypted.
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Username */}
            <div>
              <label htmlFor="username" className="sr-only">
                Username
              </label>
              <div className="relative">
                <User className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                <input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="Username"
                  required
                  disabled={loading}
                  className="w-full pl-10 pr-4 py-3 bg-slate-900/50 border border-slate-600 rounded-lg focus:border-accent focus:ring-1 focus:ring-accent text-foreground placeholder-slate-400 font-mono transition-colors disabled:opacity-50"
                />
              </div>
            </div>

            {/* Password */}
            <div>
              <label htmlFor="password" className="sr-only">
                Password
              </label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Password"
                  required
                  disabled={loading}
                  className="w-full pl-10 pr-12 py-3 bg-slate-900/50 border border-slate-600 rounded-lg focus:border-accent focus:ring-1 focus:ring-accent text-foreground placeholder-slate-400 font-mono transition-colors disabled:opacity-50"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  disabled={loading}
                  className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-accent transition-colors disabled:opacity-50"
                >
                  {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                  ) : (
                    <Eye className="w-5 h-5" />
                  )}
                </button>
              </div>
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading || !username.trim() || !password.trim()}
              className="w-full py-3 bg-accent hover:bg-accent-dark disabled:bg-slate-700 disabled:text-slate-400 text-slate-900 font-mono font-medium rounded-lg transition-colors disabled:cursor-not-allowed"
            >
              {loading ? 'Processing...' : mode === 'login' ? 'Sign In' : 'Create Account'}
            </button>
          </form>

          {/* Mode Switch */}
          <div className="mt-6 text-center">
            <p className="text-slate-400 text-sm font-mono">
              {mode === 'login' ? "Don't have an account?" : 'Already have an account?'}
            </p>
            <button
              onClick={onModeChange}
              disabled={loading}
              className="mt-1 text-accent hover:text-accent-dark font-mono text-sm underline transition-colors disabled:opacity-50"
            >
              {mode === 'login' ? 'Create one' : 'Sign in instead'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}