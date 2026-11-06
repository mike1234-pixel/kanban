import type { CSSProperties } from 'react'
import type { QuizId } from './quizzes'

// One colour per quiz, matching the highlight colours: Backend purple (like
// HTTP), Go teal, SQL orange.
const decks: Record<QuizId, { accent: string; tint: string; ink: string }> = {
  backend: { accent: '#5b43d6', tint: '#efebff', ink: '#4a33c4' },
  go: { accent: '#06889b', tint: '#e0f4f7', ink: '#035f6d' },
  sql: { accent: '#d26a22', tint: '#fdeedf', ink: '#9c4410' },
}

// CSS variables for a quiz's colours; components read var(--accent) etc.
export const deckStyle = (quiz: QuizId) =>
  ({ '--accent': decks[quiz].accent, '--tint': decks[quiz].tint, '--ink': decks[quiz].ink }) as CSSProperties
