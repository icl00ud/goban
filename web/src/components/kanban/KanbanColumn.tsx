import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useDroppable } from '@dnd-kit/core'
import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import type { Column, Card } from '@/types'
import { KanbanCard } from './KanbanCard'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreHorizontal, Plus, Trash2, X } from 'lucide-react'
import { cn } from '@/lib/utils'

interface KanbanColumnProps {
  column: Column
  onDeleteColumn: (columnId: number) => void
  onAddCard: (columnId: number, title: string) => void
  onDeleteCard: (cardId: number, columnId: number) => void
  onUpdateCard: (card: Card) => void
}

export function KanbanColumn({
  column,
  onDeleteColumn,
  onAddCard,
  onDeleteCard,
  onUpdateCard,
}: KanbanColumnProps) {
  const { t } = useTranslation()
  const [addingCard, setAddingCard] = useState(false)
  const [newCardTitle, setNewCardTitle] = useState('')

  const { setNodeRef, isOver } = useDroppable({
    id: column.id,
    data: {
      type: 'column',
      column,
    },
  })

  const cards = column.cards || []

  const handleAddCard = () => {
    if (!newCardTitle.trim()) return
    onAddCard(column.id, newCardTitle)
    setNewCardTitle('')
    setAddingCard(false)
  }

  return (
    <div
      ref={setNodeRef}
      className={cn(
        'flex-shrink-0 w-80 rounded-xl flex flex-col max-h-full transition-colors',
        'bg-[var(--color-column)]',
        isOver && 'bg-[var(--color-column-hover)] ring-2 ring-primary/20'
      )}
    >
      {/* Column Header */}
      <div className="p-3 flex items-center justify-between border-b border-border/50">
        <h3 className="font-semibold text-sm flex items-center gap-2">
          <span>{column.title}</span>
          <span className="inline-flex items-center justify-center min-w-[1.25rem] h-5 px-1.5 rounded-full bg-muted text-xs text-muted-foreground font-medium">
            {cards.length}
          </span>
        </h3>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-muted-foreground hover:text-foreground"
            >
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="animate-scale-in">
            <DropdownMenuItem
              className="text-destructive focus:text-destructive cursor-pointer"
              onClick={() => onDeleteColumn(column.id)}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              {t('column.delete')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Cards */}
      <div className="flex-1 overflow-y-auto p-2 space-y-2">
        <SortableContext
          items={cards.map((c) => c.id)}
          strategy={verticalListSortingStrategy}
        >
          {cards.map((card) => (
            <KanbanCard
              key={card.id}
              card={card}
              columnId={column.id}
              onDelete={(cardId) => onDeleteCard(cardId, column.id)}
              onUpdate={onUpdateCard}
            />
          ))}
        </SortableContext>
      </div>

      {/* Add Card */}
      <div className="p-2 border-t border-border/50">
        {addingCard ? (
          <div className="space-y-2 animate-slide-up">
            <Input
              value={newCardTitle}
              onChange={(e) => setNewCardTitle(e.target.value)}
              placeholder={t('card.titlePlaceholder')}
              autoFocus
              className="bg-background"
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleAddCard()
                if (e.key === 'Escape') {
                  setAddingCard(false)
                  setNewCardTitle('')
                }
              }}
            />
            <div className="flex gap-2">
              <Button size="sm" onClick={handleAddCard} className="flex-1">
                {t('board.add')}
              </Button>
              <Button
                size="sm"
                variant="ghost"
                onClick={() => {
                  setAddingCard(false)
                  setNewCardTitle('')
                }}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </div>
        ) : (
          <Button
            variant="ghost"
            size="sm"
            className="w-full justify-start text-muted-foreground hover:text-foreground"
            onClick={() => setAddingCard(true)}
          >
            <Plus className="mr-2 h-4 w-4" />
            {t('column.addCard')}
          </Button>
        )}
      </div>
    </div>
  )
}
