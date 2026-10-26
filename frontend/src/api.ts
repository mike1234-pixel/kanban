export type Card = {
  id: string
  column_id: string
  title: string
  description: string
  order: number
}

export type Column = {
  id: string
  board_id: string
  title: string
  position: number
  cards?: Card[]
}

export type Board = {
  id: string
  title: string
  created_at?: string
  columns: Column[]
}

export type User = {
  id: string
  name: string
  email: string
}

// ApiError carries the API's error code so callers can branch on it,
// e.g. `err instanceof ApiError && err.code === 'BOARD_NOT_FOUND'`.
export class ApiError extends Error {
  status: number
  code?: string

  constructor(message: string, status: number, code?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

// The API's error format: {"error": {"code": "...", "message": "..."}}
type ErrorBody = { error?: { code?: string; message?: string } }

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
  })
  if (!response.ok) {
    // Not every error is guaranteed to be JSON (e.g. a proxy failure), so fall back to a generic message.
    const body: ErrorBody | null = await response.json().catch(() => null)
    throw new ApiError(
      body?.error?.message ?? `Request failed (${response.status})`,
      response.status,
      body?.error?.code,
    )
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}

export const api = {
  // me returns the logged-in user, or null if there's no valid session (401).
  me: async () => {
    try {
      return await request<User>('/me')
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) return null
      throw error
    }
  },
  login: (email: string, password: string) => request<User>('/login', {
    method: 'POST', body: JSON.stringify({ email, password }),
  }),
  logout: () => request<void>('/logout', { method: 'POST' }),
  boards: () => request<Board[]>('/boards'),
  board: (id: string) => request<Board>(`/boards/${id}`),
  createBoard: (title: string) => request<Board>('/boards', {
    method: 'POST', body: JSON.stringify({ id: crypto.randomUUID(), title, columns: [] }),
  }),
  createColumn: (boardId: string, title: string, position: number) => request<Column>('/columns', {
    method: 'POST', body: JSON.stringify({ id: crypto.randomUUID(), board_id: boardId, title, position }),
  }),
  createCard: (card: Omit<Card, 'id'>) => request<Card>('/cards', {
    method: 'POST', body: JSON.stringify({ ...card, id: crypto.randomUUID() }),
  }),
  updateCard: (id: string, card: Omit<Card, 'id'>) => request<Card>(`/cards/${id}`, {
    method: 'PUT', body: JSON.stringify({ ...card, id }),
  }),
  deleteCard: (id: string) => request<void>(`/cards/${id}`, { method: 'DELETE' }),
  moveCard: (id: string, columnId: string) => request<void>(`/cards/${id}/move`, {
    method: 'POST', body: JSON.stringify({ column_id: columnId }),
  }),
}
