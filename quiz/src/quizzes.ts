import backendMd from '../../QUIZ.md?raw'
import goMd from '../../QUIZ_GO.md?raw'
import sqlMd from '../../QUIZ_SQL.md?raw'

export type QuizId = 'backend' | 'go' | 'sql'

export type Question = {
  id: string
  quiz: QuizId
  number: string
  section: string // the `##` heading, e.g. "4. Transactions"
  topic: string // the `###` heading, e.g. "BeginTx / Commit / Rollback"
  question: string // markdown
  choices: Choice[] // multiple-choice options, in the file's order
  answer: string // markdown: the explanation shown after answering
}

export type Choice = { text: string; correct: boolean }

export type Quiz = {
  id: QuizId
  title: string
  source: string
  questions: Question[]
}

// Matches a question line such as "**25a.** What does `*p` give you?"
const questionLine = /^\*\*(\d+[a-z]?)\.\*\*\s*(.*)$/
// Matches an option such as "- [x] It panics" (x marks the correct one).
const choiceLine = /^- \[( |x)\] (.*)$/

// Reads the QUIZ*.md format: `##` sections, `###` topics, then each question
// as a bold number, its options as a task list, then a <details> block
// holding the answer.
const parseQuiz = (markdown: string, quiz: QuizId): Question[] => {
  const questions: Question[] = []
  let section = ''
  let topic = ''
  let number = ''
  let question = ''
  let choices: Choice[] = []
  let mode: 'none' | 'question' | 'answer' = 'none'
  let lines: string[] = []

  for (const line of markdown.split('\n')) {
    if (mode === 'answer') {
      if (line.trim() === '</details>') {
        questions.push({ id: `${quiz}-${number}`, quiz, number, section, topic, question, choices, answer: lines.join('\n').trim() })
        mode = 'none'
      } else {
        lines.push(line)
      }
      continue
    }

    const match = questionLine.exec(line)
    if (line.startsWith('## ')) {
      section = line.slice(3).trim()
      topic = section // sections without `###` topics use their own name
    } else if (line.startsWith('### ')) {
      topic = line.slice(4).trim()
    } else if (match) {
      number = match[1]
      lines = [match[2]]
      choices = []
      mode = 'question'
    } else if (mode === 'question') {
      const choice = choiceLine.exec(line)
      if (choice) {
        choices.push({ text: choice[2], correct: choice[1] === 'x' })
      } else if (line.trim().startsWith('<details>')) {
        question = lines.join('\n').trim()
        lines = []
        mode = 'answer'
      } else {
        lines.push(line)
      }
    }
  }

  return questions
}

export const quizzes: Quiz[] = [
  { id: 'backend', title: 'Backend', source: 'QUIZ.md', questions: parseQuiz(backendMd, 'backend') },
  { id: 'go', title: 'Go', source: 'QUIZ_GO.md', questions: parseQuiz(goMd, 'go') },
  { id: 'sql', title: 'SQL', source: 'QUIZ_SQL.md', questions: parseQuiz(sqlMd, 'sql') },
]

// Groups items by key, keeping the order in which keys first appear.
export const groupBy = <T,>(items: T[], key: (item: T) => string): [string, T[]][] => {
  const groups = new Map<string, T[]>()
  for (const item of items) {
    const k = key(item)
    groups.set(k, [...(groups.get(k) ?? []), item])
  }
  return [...groups]
}

export const shuffled = <T,>(items: T[]): T[] => {
  const copy = [...items]
  for (let i = copy.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[copy[i], copy[j]] = [copy[j], copy[i]]
  }
  return copy
}
