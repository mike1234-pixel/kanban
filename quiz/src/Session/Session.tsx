import { useEffect, useMemo, useState } from 'react'
import Legend from '../Legend'
import Markdown from '../Markdown'
import { shuffled } from '../quizzes'
import type { Question, Quiz } from '../quizzes'
import { tally } from '../progress'
import type { Result } from '../progress'
import { deckStyle } from '../theme'
import styles from './Session.module.css'

type SessionProps = {
  quiz: Quiz
  label: string
  questions: Question[]
  onRecord: (id: string, result: Result) => void
  onFinish: (results: Record<string, Result>) => void
  onExit: () => void
}

const letters = ['A', 'B', 'C', 'D']

const Session = ({ quiz, label, questions, onRecord, onFinish, onExit }: SessionProps) => {
  const [index, setIndex] = useState(0)
  const [picked, setPicked] = useState<number | null>(null)
  const [results, setResults] = useState<Record<string, Result>>({})
  const question = questions[index]
  // Shuffle the options for each question, so their position isn't a clue.
  const choices = useMemo(() => shuffled(question.choices), [question])
  const correctChoice = choices.find((c) => c.correct)
  const answered = picked !== null
  const { got, missed } = tally(questions, results)

  const pick = (i: number) => {
    if (answered || i >= choices.length) return
    const result: Result = choices[i].correct ? 'got' : 'missed'
    setPicked(i)
    setResults({ ...results, [question.id]: result })
    onRecord(question.id, result)
  }

  const next = () => {
    if (!answered) return
    if (index + 1 === questions.length) {
      onFinish(results)
    } else {
      setIndex(index + 1)
      setPicked(null)
    }
  }

  // 1–4 or A–D picks an option; Enter goes to the next question.
  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      const key = event.key.toUpperCase()
      if (['1', '2', '3', '4'].includes(key)) pick(Number(key) - 1)
      else if (letters.includes(key)) pick(letters.indexOf(key))
      else if (key === 'ENTER') { event.preventDefault(); next() }
      else if (key === 'ESCAPE') onExit()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  })

  const optionClass = (i: number) => {
    if (!answered) return styles.option
    if (choices[i].correct) return `${styles.option} ${styles.right}`
    if (i === picked) return `${styles.option} ${styles.wrong}`
    return `${styles.option} ${styles.faded}`
  }

  return (
    <main className={styles.session} style={deckStyle(quiz.id)}>
      <header className={styles.topBar}>
        <button type="button" className={styles.back} onClick={onExit}>← Quizzes</button>
        <div className={styles.score} aria-label={`${got} right, ${missed} wrong`}>
          <span className={styles.gotScore}>✓ {got}</span>
          <span className={styles.missedScore}>✗ {missed}</span>
        </div>
      </header>

      <div className={styles.where}>
        <span className={styles.deckName}>{quiz.title}</span>
        <Markdown inline quiz={quiz.id} text={label} className={styles.label} />
        <span className={styles.counter}>Question {index + 1} of {questions.length}</span>
      </div>
      <div className={styles.progress}>
        <span style={{ width: `${((index + (answered ? 1 : 0)) / questions.length) * 100}%` }} />
      </div>

      <article className={styles.card}>
        <div className={styles.meta}>
          <span className={styles.number}>Q{question.number}</span>
          <Markdown inline quiz={quiz.id} text={question.topic} className={styles.topic} />
        </div>
        <Markdown quiz={quiz.id} text={question.question} className={styles.question} />

        <ol className={styles.options}>
          {choices.map((choice, i) => (
            <li key={choice.text}>
              <button type="button" className={optionClass(i)} onClick={() => pick(i)} disabled={answered}>
                <span className={styles.letter}>{answered && choice.correct ? '✓' : answered && i === picked ? '✗' : letters[i]}</span>
                <Markdown inline quiz={quiz.id} text={choice.text} />
              </button>
            </li>
          ))}
        </ol>

        {answered && correctChoice && (
          <section className={`${styles.explain} ${choices[picked].correct ? styles.explainRight : styles.explainWrong}`}>
            <p className={styles.verdict}>{choices[picked].correct ? '✓ Correct' : '✗ Not quite'}</p>
            <p className={styles.keyIdea}>
              <Markdown inline quiz={quiz.id} text={correctChoice.text} className={styles.keyText} />
            </p>
            <Markdown quiz={quiz.id} text={question.answer} className={styles.answer} />
          </section>
        )}
      </article>

      <footer className={styles.footer}>
        <Legend />
        <button type="button" className={styles.next} onClick={next} disabled={!answered}>
          {index + 1 === questions.length ? 'See results' : 'Next'} →
        </button>
      </footer>
    </main>
  )
}

export default Session
