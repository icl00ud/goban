import type {
  ApiResponse,
  User,
  Board,
  Column,
  Card,
  LoginRequest,
  RegisterRequest,
  CreateBoardRequest,
  UpdateBoardRequest,
  ReorderBoardsRequest,
  CreateColumnRequest,
  UpdateColumnRequest,
  ReorderColumnsRequest,
  CreateCardRequest,
  UpdateCardRequest,
  MoveCardRequest,
  ReorderCardsRequest,
} from '@/types'

const API_BASE = '/api/v1'

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include',
  })

  const data = await response.json()
  return data
}

// Auth API
export const authApi = {
  register: (data: RegisterRequest) =>
    request<User>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  login: (data: LoginRequest) =>
    request<User>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  logout: () =>
    request<void>('/auth/logout', {
      method: 'POST',
    }),

  me: () => request<User>('/auth/me'),
}

// Board API
export const boardApi = {
  list: () => request<Board[]>('/boards'),

  get: (id: number) => request<Board>(`/boards/${id}`),

  create: (data: CreateBoardRequest) =>
    request<Board>('/boards', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: number, data: UpdateBoardRequest) =>
    request<Board>(`/boards/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<void>(`/boards/${id}`, {
      method: 'DELETE',
    }),

  reorder: (data: ReorderBoardsRequest) =>
    request<void>('/boards/reorder', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
}

// Column API
export const columnApi = {
  create: (boardId: number, data: CreateColumnRequest) =>
    request<Column>(`/boards/${boardId}/columns`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: number, data: UpdateColumnRequest) =>
    request<Column>(`/columns/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<void>(`/columns/${id}`, {
      method: 'DELETE',
    }),

  reorder: (data: ReorderColumnsRequest) =>
    request<void>('/columns/reorder', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
}

// Card API
export const cardApi = {
  get: (id: number) => request<Card>(`/cards/${id}`),

  create: (columnId: number, data: CreateCardRequest) =>
    request<Card>(`/columns/${columnId}/cards`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: number, data: UpdateCardRequest) =>
    request<Card>(`/cards/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request<void>(`/cards/${id}`, {
      method: 'DELETE',
    }),

  move: (id: number, data: MoveCardRequest) =>
    request<Card>(`/cards/${id}/move`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  reorder: (data: ReorderCardsRequest) =>
    request<void>('/cards/reorder', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
}
