import { useState } from 'react'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import type { Card } from '@/types'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreHorizontal, Trash2, Flag } from 'lucide-react'
import { cn } from '@/lib/utils'

interface KanbanCardProps {
  card: Card
  columnId?: number
  isDragging?: boolean
  onDelete: (cardId: number) => void
  onUpdate: (card: Card) => void
}

const priorityConfig = {
  low: {
    label: 'Low',
    badgeClass: 'priority-low',
    btnClass: 'priority-btn-low',
    dotColor: 'bg-green-500',
  },
  medium: {
    label: 'Medium',
    badgeClass: 'priority-medium',
    btnClass: 'priority-btn-medium',
    dotColor: 'bg-yellow-500',
  },
  high: {
    label: 'High',
    badgeClass: 'priority-high',
    btnClass: 'priority-btn-high',
    dotColor: 'bg-red-500',
  },
}

export function KanbanCard({
  card,
  isDragging,
  onDelete,
  onUpdate,
}: KanbanCardProps) {
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editTitle, setEditTitle] = useState(card.title)
  const [editDescription, setEditDescription] = useState(card.description)
  const [editPriority, setEditPriority] = useState(card.priority)

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging: isSortableDragging,
  } = useSortable({
    id: card.id,
    data: {
      type: 'card',
      card,
    },
  })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  const handleSave = () => {
    onUpdate({
      ...card,
      title: editTitle,
      description: editDescription,
      priority: editPriority,
    })
    setDialogOpen(false)
  }

  const priority = priorityConfig[card.priority]

  return (
    <>
      <div
        ref={setNodeRef}
        style={style}
        {...attributes}
        {...listeners}
        className={cn(
          'group bg-card border rounded-lg p-3 cursor-grab active:cursor-grabbing transition-all',
          'hover:shadow-md hover:border-primary/30',
          (isDragging || isSortableDragging) && 'opacity-50 shadow-xl rotate-2 scale-105'
        )}
        onClick={() => setDialogOpen(true)}
      >
        {/* Title & Menu */}
        <div className="flex items-start justify-between gap-2">
          <p className="text-sm font-medium leading-tight">{card.title}</p>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => e.stopPropagation()}
              >
                <MoreHorizontal className="h-3 w-3" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="animate-scale-in">
              <DropdownMenuItem
                className="text-destructive focus:text-destructive cursor-pointer"
                onClick={(e) => {
                  e.stopPropagation()
                  onDelete(card.id)
                }}
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        {/* Description */}
        {card.description && (
          <p className="text-xs text-muted-foreground mt-1.5 line-clamp-2">
            {card.description}
          </p>
        )}

        {/* Priority Badge */}
        <div className="mt-3">
          <span
            className={cn(
              'inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-medium uppercase tracking-wide',
              priority.badgeClass
            )}
          >
            <Flag className="h-2.5 w-2.5" />
            {priority.label}
          </span>
        </div>
      </div>

      {/* Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Card</DialogTitle>
          </DialogHeader>
          <div className="space-y-5 mt-4">
            <div className="space-y-2">
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                className="h-11"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <textarea
                id="description"
                value={editDescription}
                onChange={(e) => setEditDescription(e.target.value)}
                className="flex min-h-[100px] w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-none"
                placeholder="Add a description..."
              />
            </div>
            <div className="space-y-2">
              <Label>Priority</Label>
              <div className="flex gap-2">
                {(['low', 'medium', 'high'] as const).map((p) => {
                  const config = priorityConfig[p]
                  const isSelected = editPriority === p
                  return (
                    <Button
                      key={p}
                      type="button"
                      variant={isSelected ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setEditPriority(p)}
                      className={cn(
                        'flex-1 capitalize',
                        isSelected && config.btnClass
                      )}
                    >
                      <div
                        className={cn(
                          'w-2 h-2 rounded-full mr-2',
                          isSelected ? 'bg-white/90' : config.dotColor
                        )}
                      />
                      {config.label}
                    </Button>
                  )
                })}
              </div>
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button variant="outline" onClick={() => setDialogOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleSave}>Save changes</Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}
