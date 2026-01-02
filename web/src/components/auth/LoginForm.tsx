import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AlertCircle, ArrowRight, Loader2 } from 'lucide-react'

export function LoginForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await login(email, password)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full max-w-md animate-fade-in">
      {/* Logo & Branding */}
      <div className="text-center mb-8">
        <Link to="/" className="inline-flex flex-col items-center gap-3 group">
          <img
            src="/logo.png"
            alt="GoBan"
            className="h-16 w-auto drop-shadow-lg transition-transform group-hover:scale-105"
          />
          <div>
            <h1 className="text-3xl font-bold tracking-tight">GoBan</h1>
            <p className="text-sm text-muted-foreground">
              Your minimalist Kanban board
            </p>
          </div>
        </Link>
      </div>

      {/* Form Card */}
      <div className="bg-card border rounded-xl p-8 shadow-lg shadow-black/5 dark:shadow-black/20">
        <div className="mb-6">
          <h2 className="text-xl font-semibold">Welcome back</h2>
          <p className="text-sm text-muted-foreground mt-1">
            Sign in to continue to your boards
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div className="flex items-center gap-2 p-3 text-sm text-destructive bg-destructive/10 border border-destructive/20 rounded-lg animate-slide-up">
              <AlertCircle className="h-4 w-4 flex-shrink-0" />
              <span>{error}</span>
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
              className="h-11"
              autoComplete="email"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              required
              className="h-11"
              autoComplete="current-password"
            />
          </div>

          <Button
            type="submit"
            className="w-full h-11 text-base font-medium"
            disabled={loading}
          >
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Signing in...
              </>
            ) : (
              <>
                Sign in
                <ArrowRight className="ml-2 h-4 w-4" />
              </>
            )}
          </Button>
        </form>

        <div className="mt-6 pt-6 border-t text-center">
          <p className="text-sm text-muted-foreground">
            Don't have an account?{' '}
            <Link
              to="/register"
              className="text-primary font-medium hover:underline underline-offset-4"
            >
              Create one
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
