import type { Question } from './quizzes'

export type Result = 'got' | 'missed'

// The latest result for each question id, kept in the browser between visits.
export type Progress = Record<string, Result>

const storageKey = 'quiz-progress'

export const loadProgress = (): Progress => {
  try {
    return JSON.parse(localStorage.getItem(storageKey) ?? '{}')
  } catch {
    return {}
  }
}

export const saveProgress = (progress: Progress) => {
  try {
    localStorage.setItem(storageKey, JSON.stringify(progress))
  } catch {
    // Storage can be blocked (e.g. private browsing); progress just isn't kept.
  }
}

export const tally = (questions: Question[], results: Record<string, Result>) => {
  const got = questions.filter((q) => results[q.id] === 'got').length
  const missed = questions.filter((q) => results[q.id] === 'missed').length
  return { got, missed, unseen: questions.length - got - missed }
}
