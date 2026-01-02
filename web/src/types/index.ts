// User types
export interface User {
  id: number
  email: string
  name: string
}

// Board types
export interface Board {
  id: number
  name: string
  description: string
  color: string
  position: number
  columns?: Column[]
  created_at: string
}

// Column types
export interface Column {
  id: number
  title: string
  position: number
  board_id: number
  cards?: Card[]
}

// Card types
export interface Card {
  id: number
  title: string
  description: string
  position: number
  priority: 'low' | 'medium' | 'high'
  column_id: number
}

// API Response types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// Auth request types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

// Board request types
export interface CreateBoardRequest {
  name: string
  description?: string
  color?: string
}

export interface UpdateBoardRequest {
  name?: string
  description?: string
  color?: string
}

export interface ReorderBoardsRequest {
  board_ids: number[]
}

// Column request types
export interface CreateColumnRequest {
  title: string
}

export interface UpdateColumnRequest {
  title?: string
}

export interface ReorderColumnsRequest {
  board_id: number
  column_ids: number[]
}

// Card request types
export interface CreateCardRequest {
  title: string
  description?: string
  priority?: 'low' | 'medium' | 'high'
}

export interface UpdateCardRequest {
  title?: string
  description?: string
  priority?: 'low' | 'medium' | 'high'
}

export interface MoveCardRequest {
  target_column_id: number
  position: number
}

export interface ReorderCardsRequest {
  column_id: number
  card_ids: number[]
}
