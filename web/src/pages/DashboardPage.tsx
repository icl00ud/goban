import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import {
  DndContext,
  DragOverlay,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragStartEvent,
  type DragEndEvent,
} from '@dnd-kit/core'
import {
  SortableContext,
  sortableKeyboardCoordinates,
  rectSortingStrategy,
  useSortable,
  arrayMove,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { boardApi } from '@/lib/api'
import type { Board } from '@/types'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Plus, Kanban, ArrowRight, Loader2, GripVertical } from 'lucide-react'
import { cn } from '@/lib/utils'

const BOARD_COLORS = [
  '#2dd4bf', // teal (primary)
  '#3b82f6', // blue
  '#8b5cf6', // violet
  '#ec4899', // pink
  '#f97316', // orange
  '#22c55e', // green
  '#eab308', // yellow
  '#6366f1', // indigo
]

// Sortable Board Card Component
function SortableBoardCard({ board, t }: { board: Board; t: (key: string, options?: Record<string, unknown>) => string }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: board.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  const columnCount = board.columns?.length || 0
  const cardCount = board.columns?.reduce((acc, col) => acc + (col.cards?.length || 0), 0) || 0

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        'group relative',
        isDragging && 'z-50'
      )}
    >
      <div
        className={cn(
          'relative bg-card border rounded-xl p-5 h-full transition-all duration-200',
          'hover:shadow-lg hover:shadow-black/5 dark:hover:shadow-black/20 hover:border-primary/30',
          isDragging && 'shadow-2xl scale-105 rotate-2 opacity-90'
        )}
      >
        {/* Drag Handle */}
        <button
          {...attributes}
          {...listeners}
          className="absolute top-3 right-3 p-1.5 rounded-md opacity-0 group-hover:opacity-100 hover:bg-muted transition-all cursor-grab active:cursor-grabbing"
          onClick={(e) => e.preventDefault()}
        >
          <GripVertical className="h-4 w-4 text-muted-foreground" />
        </button>

        {/* Color accent */}
        <div
          className="absolute top-0 left-0 right-0 h-1 rounded-t-xl"
          style={{ backgroundColor: board.color }}
        />

        <Link to={`/boards/${board.id}`} className="block mt-2">
          <div className="flex items-start justify-between gap-3 pr-8">
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-lg truncate group-hover:text-primary transition-colors">
                {board.name}
              </h3>
              {board.description && (
                <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                  {board.description}
                </p>
              )}
            </div>
          </div>

          {/* Board stats */}
          <div className="flex items-center justify-between mt-4 pt-4 border-t">
            <div className="flex items-center gap-3 text-xs text-muted-foreground">
              <span>{t('dashboard.columns', { count: columnCount })}</span>
              <span>{t('dashboard.cards', { count: cardCount })}</span>
            </div>
            <ArrowRight className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 group-hover:text-primary transition-all" />
          </div>
        </Link>
      </div>
    </div>
  )
}

// Board Card for Drag Overlay
function BoardCardOverlay({ board }: { board: Board }) {
  return (
    <div className="bg-card border rounded-xl p-5 shadow-2xl scale-105 rotate-3 opacity-95 w-full max-w-sm">
      <div
        className="absolute top-0 left-0 right-0 h-1 rounded-t-xl"
        style={{ backgroundColor: board.color }}
      />
      <div className="mt-2">
        <h3 className="font-semibold text-lg truncate">{board.name}</h3>
        {board.description && (
          <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
            {board.description}
          </p>
        )}
      </div>
    </div>
  )
}

export function DashboardPage() {
  const { t } = useTranslation()
  const [boards, setBoards] = useState<Board[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [newBoardName, setNewBoardName] = useState('')
  const [newBoardDescription, setNewBoardDescription] = useState('')
  const [newBoardColor, setNewBoardColor] = useState(BOARD_COLORS[0])
  const [creating, setCreating] = useState(false)
  const [activeBoard, setActiveBoard] = useState<Board | null>(null)

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  )

  useEffect(() => {
    loadBoards()
  }, [])

  const loadBoards = async () => {
    try {
      const response = await boardApi.list()
      if (response.success && response.data) {
        setBoards(response.data)
      }
    } catch (err) {
      console.error('Failed to load boards:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateBoard = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newBoardName.trim()) return

    setCreating(true)
    try {
      const response = await boardApi.create({
        name: newBoardName,
        description: newBoardDescription,
        color: newBoardColor,
      })
      if (response.success && response.data) {
        setBoards([...boards, response.data])
        setNewBoardName('')
        setNewBoardDescription('')
        setNewBoardColor(BOARD_COLORS[0])
        setDialogOpen(false)
      }
    } catch (err) {
      console.error('Failed to create board:', err)
    } finally {
      setCreating(false)
    }
  }

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event
    const board = boards.find((b) => b.id === active.id)
    if (board) {
      setActiveBoard(board)
    }
  }

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event
    setActiveBoard(null)

    if (!over || active.id === over.id) return

    const oldIndex = boards.findIndex((b) => b.id === active.id)
    const newIndex = boards.findIndex((b) => b.id === over.id)

    if (oldIndex === -1 || newIndex === -1) return

    // Optimistic update
    const newBoards = arrayMove(boards, oldIndex, newIndex)
    setBoards(newBoards)

    // Persist to backend
    try {
      await boardApi.reorder({
        board_ids: newBoards.map((b) => b.id),
      })
    } catch (err) {
      console.error('Failed to reorder boards:', err)
      // Revert on error
      setBoards(boards)
    }
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-3">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="text-sm text-muted-foreground">{t('dashboard.loadingBoards')}</p>
      </div>
    )
  }

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{t('dashboard.title')}</h1>
          <p className="text-muted-foreground mt-1">
            {boards.length === 0
              ? t('dashboard.emptyState')
              : `${t('dashboard.boardCount', { count: boards.length })} • ${t('dashboard.dragHint')}`}
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              {t('dashboard.newBoard')}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>{t('board.create')}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateBoard} className="space-y-5 mt-4">
              <div className="space-y-2">
                <Label htmlFor="name">{t('board.name')}</Label>
                <Input
                  id="name"
                  value={newBoardName}
                  onChange={(e) => setNewBoardName(e.target.value)}
                  placeholder={t('board.namePlaceholder')}
                  required
                  className="h-11"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">{t('board.description')}</Label>
                <Input
                  id="description"
                  value={newBoardDescription}
                  onChange={(e) => setNewBoardDescription(e.target.value)}
                  placeholder={t('board.descriptionPlaceholder')}
                  className="h-11"
                />
              </div>
              <div className="space-y-2">
                <Label>{t('board.color')}</Label>
                <div className="flex gap-2 flex-wrap">
                  {BOARD_COLORS.map((color) => (
                    <button
                      key={color}
                      type="button"
                      onClick={() => setNewBoardColor(color)}
                      className="w-8 h-8 rounded-lg transition-all hover:scale-110 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      style={{
                        backgroundColor: color,
                        boxShadow: newBoardColor === color ? `0 0 0 2px var(--color-background), 0 0 0 4px ${color}` : 'none',
                      }}
                    />
                  ))}
                </div>
              </div>
              <Button type="submit" className="w-full h-11" disabled={creating}>
                {creating ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    {t('board.creating')}
                  </>
                ) : (
                  t('board.createButton')
                )}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Empty State */}
      {boards.length === 0 ? (
        <div className="text-center py-16 animate-fade-in">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-primary/10 mb-4">
            <Kanban className="h-8 w-8 text-primary" />
          </div>
          <h2 className="text-xl font-semibold mb-2">{t('dashboard.noBoardsTitle')}</h2>
          <p className="text-muted-foreground mb-6 max-w-sm mx-auto">
            {t('dashboard.noBoardsDescription')}
          </p>
          <Button onClick={() => setDialogOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            {t('dashboard.createFirst')}
          </Button>
        </div>
      ) : (
        /* Board Grid with DnD */
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <SortableContext items={boards.map((b) => b.id)} strategy={rectSortingStrategy}>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {boards.map((board, index) => (
                <div
                  key={board.id}
                  className="animate-slide-up"
                  style={{ animationDelay: `${index * 50}ms` }}
                >
                  <SortableBoardCard board={board} t={t} />
                </div>
              ))}
            </div>
          </SortableContext>

          <DragOverlay>
            {activeBoard && <BoardCardOverlay board={activeBoard} />}
          </DragOverlay>
        </DndContext>
      )}
    </div>
  )
}
