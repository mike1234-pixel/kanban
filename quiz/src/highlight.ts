import { Marked } from 'marked'
import type { QuizId } from './quizzes'

// Every piece of inline code is coloured by what kind of thing it is, so the
// same colour always means the same kind: purple HTTP, orange SQL, teal Go.
export type Kind = 'http' | 'sql' | 'go'

const httpMethods = /^(GET|POST|PUT|PATCH|DELETE)( \/.*)?$/
const httpWords = new Set(['Content-Type', 'Set-Cookie', 'Cookie', 'HttpOnly', 'Secure', 'SameSite', 'Max-Age', 'Expires', 'Authorization', 'WriteHeader'])
const sqlWords = new Set([
  'SELECT', 'INSERT', 'UPDATE', 'DELETE', 'CREATE', 'ALTER', 'DROP', 'FROM', 'WHERE', 'JOIN', 'LEFT', 'INNER', 'GROUP', 'ORDER', 'UNIQUE',
  'PRIMARY', 'NOT', 'NULL', 'REFERENCES', 'CHECK', 'DEFAULT', 'BEGIN', 'COMMIT', 'ROLLBACK', 'VARCHAR', 'TEXT', 'TIMESTAMPTZ', 'TIMESTAMP',
  'DATE', 'TIME', 'UUID', 'INTEGER', 'BIGSERIAL', 'ON', 'EXISTS', 'FOR', 'LIMIT', 'OFFSET', 'INDEX', 'RETURNING', 'READ', 'EXPLAIN', 'USING', 'SET', 'VALUES',
])
const sqlPlanWords = new Set(['Seq Scan', 'Index Cond', 'Filter', '"order"'])

const tables = /^(users|sessions|boards|columns|cards|card_movements)(\.[a-z_]+)?$|→/ // `boards`, `boards.owner_id`, `cards → columns`
// Go code looks like camelCase or PascalCase names, `sql.` calls, package paths, string literals or `nil`.
const looksLikeGo = (code: string) => /[a-z][A-Z]|^[A-Z][a-z]|^sql\.|^"|\(&|\//.test(code) || code === 'nil'

export const kindOf = (code: string, quiz: QuizId): Kind => {
  const firstWord = /^[A-Z]+/.exec(code)?.[0] ?? ''
  // A bare DELETE is SQL in the SQL quiz and an HTTP method everywhere else.
  if (quiz === 'sql' && sqlWords.has(code)) return 'sql'
  if (/^[1-5]\d\d\b/.test(code) || httpMethods.test(code) || httpWords.has(code)) return 'http'
  if (/^(application|text)\//.test(code) || /^[A-Z]+(_[A-Z]+)+$/.test(code)) return 'http' // content types, API error codes
  if (/^(\/|\?|https?:|localhost)/.test(code) || code === 'code' || code === 'message') return 'http' // URL parts, error JSON fields
  if (/^(\d{5}|23)/.test(code) || sqlPlanWords.has(code) || sqlWords.has(firstWord)) return 'sql' // Postgres error codes, SQL
  if (tables.test(code) || /^([a-z]+\.)?[a-z]+(_[a-z]+)+$/.test(code) || /^\$\d/.test(code)) return 'sql' // tables, snake_case columns, placeholders
  if (quiz === 'sql' && !looksLikeGo(code)) return 'sql' // aliases like `s` and `s.user_id`
  return 'go'
}

const escapeHtml = (text: string) =>
  text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')

const parsers = new Map<QuizId, Marked>()

// A markdown parser whose inline code carries its kind as a class.
export const parserFor = (quiz: QuizId) => {
  let parser = parsers.get(quiz)
  if (!parser) {
    parser = new Marked({
      renderer: {
        codespan: ({ text }) => `<code class="kind-${kindOf(text, quiz)}">${escapeHtml(text)}</code>`,
      },
    })
    parsers.set(quiz, parser)
  }
  return parser
}
