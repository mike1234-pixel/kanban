import { useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api, ApiError, type Board, type Card } from '../api'
import type { CardFormValues } from '../CardEditor'
import type { TitleFormValues } from '../TitleEditor'

export type WorkspaceModal = 'card' | 'column' | 'board' | null

export const useBoardWorkspace = () => {
  const queryClient = useQueryClient()
  const boardsQuery = useQuery({ queryKey: ['boards'], queryFn: api.boards })
  const [activeId, setActiveId] = useState('')
  const [modal, setModal] = useState<WorkspaceModal>(null)
  const [editingCard, setEditingCard] = useState<Card | null>(null)
  const [targetColumnId, setTargetColumnId] = useState('')
  const [draggingCardId, setDraggingCardId] = useState('')
  const [notice, setNotice] = useState('')

  useEffect(() => {
    if (!notice) return
    const timer = window.setTimeout(() => setNotice(''), 2800)
    return () => window.clearTimeout(timer)
  }, [notice])

  const boards = boardsQuery.data ?? []
  const activeBoard = boards.find((board) => board.id === activeId) ?? boards[0]
  const boardQuery = useQuery({
    queryKey: ['board', activeBoard?.id],
    queryFn: () => api.board(activeBoard!.id),
    enabled: Boolean(activeBoard?.id),
  })
  const board: Board | undefined = boardQuery.data ?? activeBoard
  const columns = board?.columns ?? []
  const cardCount = useMemo(
    () => columns.reduce((count, column) => count + (column.cards?.length ?? 0), 0),
    [columns],
  )

  const refreshBoard = async () => {
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: ['boards'] }),
      queryClient.invalidateQueries({ queryKey: ['board', activeBoard?.id] }),
    ])
  }

  const mutation = useMutation({
    mutationFn: (action: () => Promise<unknown>) => action(),
    onSuccess: async (result) => {
      if (result && typeof result === 'object' && 'columns' in result && 'id' in result) {
        setActiveId(String(result.id))
      }
      await refreshBoard()
      setModal(null)
      setEditingCard(null)
      setNotice('Changes saved')
    },
    onError: async (error) => {
      // Someone else changed the card since it was opened: load the latest
      // version and close the editor, so the next edit starts from it.
      if (error instanceof ApiError && error.code === 'CARD_VERSION_CONFLICT') {
        await refreshBoard()
        setModal(null)
        setEditingCard(null)
      }
      setNotice(error instanceof Error ? error.message : 'Something went wrong')
    },
  })

  const openCard = (columnId: string, card?: Card) => {
    setTargetColumnId(columnId)
    setEditingCard(card ?? null)
    setModal('card')
  }

  const closeModal = () => {
    setModal(null)
    setEditingCard(null)
  }

  const saveCard = ({ title, description }: CardFormValues) => {
    if (editingCard) {
      mutation.mutate(() => api.updateCard(editingCard.id, {
        title,
        description,
        column_id: editingCard.column_id,
        order: editingCard.order,
        version: editingCard.version,
      }))
      return
    }

    const order = (columns.find((column) => column.id === targetColumnId)?.cards ?? []).length
    mutation.mutate(() => api.createCard({ title, description, column_id: targetColumnId, order }))
  }

  const saveTitle = ({ title }: TitleFormValues) => {
    if (modal === 'column' && board) {
      mutation.mutate(() => api.createColumn(board.id, title, columns.length))
    } else if (modal === 'board') {
      mutation.mutate(() => api.createBoard(title))
    }
  }

  const moveCard = (cardId: string, columnId: string) => {
    setDraggingCardId('')
    if (cardId && columnId) mutation.mutate(() => api.moveCard(cardId, columnId))
  }

  const deleteCard = (card: Card) => mutation.mutate(() => api.deleteCard(card.id))

  const notify = (message: string) => setNotice(message)

  return {
    activeBoard,
    board,
    boardQuery,
    boards,
    boardsQuery,
    cardCount,
    closeModal,
    columns,
    deleteCard,
    draggingCardId,
    editingCard,
    modal,
    moveCard,
    mutation,
    notice,
    notify,
    openCard,
    saveCard,
    saveTitle,
    setActiveId,
    setDraggingCardId,
    setModal,
  }
}
