import { parserFor } from '../highlight'
import type { QuizId } from '../quizzes'
import styles from './Markdown.module.css'

type MarkdownProps = {
  text: string
  quiz: QuizId
  inline?: boolean
  className?: string
}

// Renders text from the repo's own QUIZ*.md files, so the HTML is trusted.
const Markdown = ({ text, quiz, inline = false, className = '' }: MarkdownProps) => {
  const parser = parserFor(quiz)
  if (inline) {
    return <span className={`${styles.markdown} ${className}`} dangerouslySetInnerHTML={{ __html: parser.parseInline(text, { async: false }) }} />
  }
  return <div className={`${styles.markdown} ${className}`} dangerouslySetInnerHTML={{ __html: parser.parse(text, { async: false }) }} />
}

export default Markdown
