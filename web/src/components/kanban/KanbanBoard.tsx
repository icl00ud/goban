import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  DndContext,
  DragOverlay,
  closestCorners,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragStartEvent,
  type DragEndEvent,
  type DragOverEvent,
} from '@dnd-kit/core'
import {
  SortableContext,
  horizontalListSortingStrategy,
  arrayMove,
} from '@dnd-kit/sortable'
import type { Board, Column, Card } from '@/types'
import { columnApi, cardApi } from '@/lib/api'
import { KanbanColumn } from './KanbanColumn'
import { KanbanCard } from './KanbanCard'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Plus, X } from 'lucide-react'

interface KanbanBoardProps {
  board: Board
  onBoardUpdate?: (board: Board) => void
}

export function KanbanBoard({ board }: KanbanBoardProps) {
  const { t } = useTranslation()
  const [columns, setColumns] = useState<Column[]>(board.columns || [])
  const [activeCard, setActiveCard] = useState<Card | null>(null)
  const [newColumnTitle, setNewColumnTitle] = useState('')
  const [addingColumn, setAddingColumn] = useState(false)

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor)
  )

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event
    const activeData = active.data.current

    if (activeData?.type === 'card') {
      setActiveCard(activeData.card)
    }
  }

  const handleDragOver = (event: DragOverEvent) => {
    const { active, over } = event
    if (!over) return

    const activeData = active.data.current
    const overData = over.data.current

    if (activeData?.type !== 'card') return

    const activeCard = activeData.card as Card
    const activeColumnId = activeCard.column_id

    let overColumnId: number

    if (overData?.type === 'card') {
      overColumnId = (overData.card as Card).column_id
    } else if (overData?.type === 'column') {
      overColumnId = overData.column.id
    } else {
      return
    }

    if (activeColumnId === overColumnId) return

    // Move card to different column
    setColumns((prev) => {
      const newColumns = prev.map((col) => {
        if (col.id === activeColumnId) {
          return {
            ...col,
            cards: col.cards?.filter((c) => c.id !== activeCard.id) || [],
          }
        }
        if (col.id === overColumnId) {
          const updatedCard = { ...activeCard, column_id: overColumnId }
          return {
            ...col,
            cards: [...(col.cards || []), updatedCard],
          }
        }
        return col
      })
      return newColumns
    })
  }

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event
    setActiveCard(null)

    if (!over) return

    const activeData = active.data.current
    const overData = over.data.current

    if (activeData?.type === 'card') {
      const card = activeData.card as Card
      const overColumnId = overData?.type === 'card'
        ? (overData.card as Card).column_id
        : overData?.type === 'column'
          ? overData.column.id
          : card.column_id

      // Find position
      const targetColumn = columns.find((c) => c.id === overColumnId)
      const cards = targetColumn?.cards || []
      const overIndex = overData?.type === 'card'
        ? cards.findIndex((c) => c.id === (overData.card as Card).id)
        : cards.length

      try {
        if (card.column_id !== overColumnId) {
          // Move to different column
          await cardApi.move(card.id, {
            target_column_id: overColumnId,
            position: overIndex >= 0 ? overIndex : 0,
          })
        } else {
          // Reorder within same column
          const cardIds = cards.map((c) => c.id)
          const oldIndex = cardIds.indexOf(card.id)
          const newCardIds = arrayMove(cardIds, oldIndex, overIndex >= 0 ? overIndex : cardIds.length - 1)

          await cardApi.reorder({
            column_id: overColumnId,
            card_ids: newCardIds,
          })
        }
      } catch (err) {
        console.error('Failed to move card:', err)
      }
    }
  }

  const handleAddColumn = async () => {
    if (!newColumnTitle.trim()) return

    try {
      const response = await columnApi.create(board.id, { title: newColumnTitle })
      if (response.success && response.data) {
        setColumns([...columns, { ...response.data, cards: [] }])
        setNewColumnTitle('')
        setAddingColumn(false)
      }
    } catch (err) {
      console.error('Failed to create column:', err)
    }
  }

  const handleDeleteColumn = async (columnId: number) => {
    try {
      const response = await columnApi.delete(columnId)
      if (response.success) {
        setColumns(columns.filter((c) => c.id !== columnId))
      }
    } catch (err) {
      console.error('Failed to delete column:', err)
    }
  }

  const handleAddCard = async (columnId: number, title: string) => {
    try {
      const response = await cardApi.create(columnId, { title })
      if (response.success && response.data) {
        setColumns(columns.map((col) => {
          if (col.id === columnId) {
            return {
              ...col,
              cards: [...(col.cards || []), response.data!],
            }
          }
          return col
        }))
      }
    } catch (err) {
      console.error('Failed to create card:', err)
    }
  }

  const handleDeleteCard = async (cardId: number, columnId: number) => {
    try {
      const response = await cardApi.delete(cardId)
      if (response.success) {
        setColumns(columns.map((col) => {
          if (col.id === columnId) {
            return {
              ...col,
              cards: col.cards?.filter((c) => c.id !== cardId) || [],
            }
          }
          return col
        }))
      }
    } catch (err) {
      console.error('Failed to delete card:', err)
    }
  }

  const handleUpdateCard = async (card: Card) => {
    try {
      const response = await cardApi.update(card.id, {
        title: card.title,
        description: card.description,
        priority: card.priority,
      })
      if (response.success && response.data) {
        setColumns(columns.map((col) => {
          if (col.id === card.column_id) {
            return {
              ...col,
              cards: col.cards?.map((c) =>
                c.id === card.id ? response.data! : c
              ) || [],
            }
          }
          return col
        }))
      }
    } catch (err) {
      console.error('Failed to update card:', err)
    }
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCorners}
      onDragStart={handleDragStart}
      onDragOver={handleDragOver}
      onDragEnd={handleDragEnd}
    >
      <div className="flex gap-4 p-4 h-full overflow-x-auto">
        <SortableContext
          items={columns.map((c) => `column-${c.id}`)}
          strategy={horizontalListSortingStrategy}
        >
          {columns.map((column, index) => (
            <div
              key={column.id}
              className="animate-slide-up"
              style={{ animationDelay: `${index * 50}ms` }}
            >
              <KanbanColumn
                column={column}
                onDeleteColumn={handleDeleteColumn}
                onAddCard={handleAddCard}
                onDeleteCard={handleDeleteCard}
                onUpdateCard={handleUpdateCard}
              />
            </div>
          ))}
        </SortableContext>

        {/* Add Column */}
        <div className="flex-shrink-0 w-80">
          {addingColumn ? (
            <div className="bg-[var(--color-column)] rounded-xl p-3 space-y-2 animate-scale-in">
              <Input
                value={newColumnTitle}
                onChange={(e) => setNewColumnTitle(e.target.value)}
                placeholder={t('board.columnPlaceholder')}
                autoFocus
                className="bg-background"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') handleAddColumn()
                  if (e.key === 'Escape') {
                    setAddingColumn(false)
                    setNewColumnTitle('')
                  }
                }}
              />
              <div className="flex gap-2">
                <Button size="sm" onClick={handleAddColumn} className="flex-1">
                  {t('board.add')}
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => {
                    setAddingColumn(false)
                    setNewColumnTitle('')
                  }}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ) : (
            <Button
              variant="ghost"
              className="w-full justify-start h-12 border-2 border-dashed border-border hover:border-primary/50 hover:bg-primary/5 rounded-xl transition-colors"
              onClick={() => setAddingColumn(true)}
            >
              <Plus className="mr-2 h-4 w-4" />
              {t('board.addColumn')}
            </Button>
          )}
        </div>
      </div>

      <DragOverlay>
        {activeCard && (
          <div className="rotate-3 opacity-90">
            <KanbanCard
              card={activeCard}
              columnId={activeCard.column_id}
              isDragging
              onDelete={() => {}}
              onUpdate={() => {}}
            />
          </div>
        )}
      </DragOverlay>
    </DndContext>
  )
}
