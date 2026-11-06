import Markdown from '../Markdown'
import MasteryBar from '../MasteryBar'
import { groupBy } from '../quizzes'
import type { Question, Quiz } from '../quizzes'
import { tally } from '../progress'
import type { Result } from '../progress'
import { deckStyle } from '../theme'
import styles from './Results.module.css'

type ResultsProps = {
  quiz: Quiz
  label: string
  questions: Question[]
  results: Record<string, Result>
  onRetry: (questions: Question[]) => void
  onHome: () => void
}

const radius = 54
const circumference = 2 * Math.PI * radius

const Results = ({ quiz, label, questions, results, onRetry, onHome }: ResultsProps) => {
  const { got } = tally(questions, results)
  const missed = questions.filter((q) => results[q.id] === 'missed')
  const percent = Math.round((got / questions.length) * 100)

  return (
    <main className={styles.results} style={deckStyle(quiz.id)}>
      <section className={styles.summary}>
        <svg className={styles.ring} viewBox="0 0 128 128" aria-hidden="true">
          <circle cx="64" cy="64" r={radius} className={styles.track} />
          <circle cx="64" cy="64" r={radius} className={styles.progress} strokeDasharray={circumference} strokeDashoffset={circumference * (1 - got / questions.length)} />
        </svg>
        <div className={styles.percent}>{percent}%</div>
        <h1>{got} / {questions.length}</h1>
        <p className={styles.label}>
          {quiz.title} · <Markdown inline quiz={quiz.id} text={label} />
        </p>
        <div className={styles.actions}>
          {missed.length > 0 && (
            <button type="button" className={styles.retry} onClick={() => onRetry(missed)}>Redo the {missed.length} missed</button>
          )}
          <button type="button" className={styles.home} onClick={onHome}>Back to quizzes</button>
        </div>
      </section>

      <section className={styles.panel}>
        <h2>By topic</h2>
        {groupBy(questions, (q) => q.topic).map(([topic, topicQuestions]) => {
          const counts = tally(topicQuestions, results)
          return (
            <div key={topic} className={styles.topicRow}>
              <Markdown inline quiz={quiz.id} text={topic} className={styles.topicName} />
              <MasteryBar {...counts} />
              <span className={styles.topicScore}>{counts.got}/{topicQuestions.length}</span>
            </div>
          )
        })}
      </section>

      {missed.length > 0 && (
        <section className={styles.panel}>
          <h2>Worth another look</h2>
          {missed.map((q) => (
            <article key={q.id} className={styles.missedCard}>
              <div className={styles.missedMeta}>
                <span className={styles.number}>Q{q.number}</span>
                <Markdown inline quiz={quiz.id} text={q.topic} className={styles.topicName} />
              </div>
              <Markdown quiz={quiz.id} text={q.question} className={styles.missedQuestion} />
              <div className={styles.missedAnswer}>
                <p className={styles.keyIdea}>
                  <Markdown inline quiz={quiz.id} text={q.choices.find((c) => c.correct)?.text ?? ''} className={styles.keyText} />
                </p>
                <Markdown quiz={quiz.id} text={q.answer} />
              </div>
            </article>
          ))}
        </section>
      )}
    </main>
  )
}

export default Results
