import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { boardApi } from '@/lib/api'
import type { Board } from '@/types'
import { Button } from '@/components/ui/button'
import { KanbanBoard } from '@/components/kanban/KanbanBoard'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { ArrowLeft, Trash2, MoreVertical, Loader2, AlertCircle } from 'lucide-react'

export function BoardPage() {
  const { t } = useTranslation()
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [board, setBoard] = useState<Board | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadBoard()
  }, [id])

  const loadBoard = async () => {
    if (!id) return

    try {
      const response = await boardApi.get(parseInt(id))
      if (response.success && response.data) {
        setBoard(response.data)
      } else {
        setError(response.error || t('errors.loadBoard'))
      }
    } catch {
      setError(t('errors.loadBoard'))
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteBoard = async () => {
    if (!board) return
    if (!confirm(t('column.deleteConfirm'))) return

    try {
      const response = await boardApi.delete(board.id)
      if (response.success) {
        navigate('/')
      }
    } catch (err) {
      console.error('Failed to delete board:', err)
    }
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-3 animate-fade-in">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="text-sm text-muted-foreground">{t('board.loadingBoard')}</p>
      </div>
    )
  }

  if (error || !board) {
    return (
      <div className="container mx-auto p-6 max-w-md animate-fade-in">
        <div className="text-center py-12">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-destructive/10 mb-4">
            <AlertCircle className="h-6 w-6 text-destructive" />
          </div>
          <h2 className="text-lg font-semibold mb-2">{t('board.notFound')}</h2>
          <p className="text-muted-foreground mb-6">
            {error || t('board.notFoundDescription')}
          </p>
          <Button onClick={() => navigate('/')}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            {t('board.backToDashboard')}
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="h-[calc(100vh-4rem)] flex flex-col animate-fade-in">
      {/* Board Header */}
      <div className="border-b bg-background/50 backdrop-blur-sm px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate('/')}
            className="h-9 w-9"
          >
            <ArrowLeft className="h-5 w-5" />
          </Button>

          <div className="flex items-center gap-3">
            <div
              className="w-3 h-3 rounded-full ring-2 ring-background shadow-sm"
              style={{ backgroundColor: board.color }}
            />
            <div>
              <h1 className="text-lg font-semibold leading-tight">{board.name}</h1>
              {board.description && (
                <p className="text-xs text-muted-foreground">{board.description}</p>
              )}
            </div>
          </div>
        </div>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-9 w-9">
              <MoreVertical className="h-5 w-5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="animate-scale-in">
            <DropdownMenuItem
              className="text-destructive focus:text-destructive cursor-pointer"
              onClick={handleDeleteBoard}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              {t('board.delete')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Kanban Board */}
      <div className="flex-1 overflow-hidden">
        <KanbanBoard board={board} onBoardUpdate={setBoard} />
      </div>
    </div>
  )
}
