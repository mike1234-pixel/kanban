import { useState } from 'react'
import Home from '../Home'
import Session from '../Session'
import Results from '../Results'
import { shuffled } from '../quizzes'
import type { Question, Quiz } from '../quizzes'
import { loadProgress, saveProgress } from '../progress'
import type { Progress, Result } from '../progress'
import styles from './App.module.css'

type Screen =
  | { kind: 'home' }
  | { kind: 'session'; quiz: Quiz; label: string; questions: Question[]; run: number }
  | { kind: 'results'; quiz: Quiz; label: string; questions: Question[]; results: Record<string, Result> }

const App = () => {
  const [screen, setScreen] = useState<Screen>({ kind: 'home' })
  const [progress, setProgress] = useState<Progress>(loadProgress)
  const [shuffle, setShuffle] = useState(false)

  // run changes on every start, so a retry gets a fresh Session (back to card 1).
  const start = (quiz: Quiz, label: string, questions: Question[]) =>
    setScreen({ kind: 'session', quiz, label, questions: shuffle ? shuffled(questions) : questions, run: Date.now() })

  const record = (id: string, result: Result) => {
    setProgress((current) => {
      const next = { ...current, [id]: result }
      saveProgress(next)
      return next
    })
  }

  return (
    <div className={styles.app}>
      {screen.kind === 'home' && <Home progress={progress} shuffle={shuffle} onShuffleChange={setShuffle} onStart={start} />}
      {screen.kind === 'session' && (
        <Session
          key={screen.run}
          quiz={screen.quiz}
          label={screen.label}
          questions={screen.questions}
          onRecord={record}
          onFinish={(results) => setScreen({ ...screen, kind: 'results', results })}
          onExit={() => setScreen({ kind: 'home' })}
        />
      )}
      {screen.kind === 'results' && (
        <Results
          quiz={screen.quiz}
          label={screen.label}
          questions={screen.questions}
          results={screen.results}
          onRetry={(missed) => start(screen.quiz, `${screen.label}: missed`, missed)}
          onHome={() => setScreen({ kind: 'home' })}
        />
      )}
    </div>
  )
}

export default App
