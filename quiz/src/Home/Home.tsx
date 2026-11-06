import Legend from '../Legend'
import MasteryBar from '../MasteryBar'
import Markdown from '../Markdown'
import { groupBy, quizzes } from '../quizzes'
import type { Question, Quiz } from '../quizzes'
import { tally } from '../progress'
import type { Progress } from '../progress'
import { deckStyle } from '../theme'
import styles from './Home.module.css'

type HomeProps = {
  progress: Progress
  shuffle: boolean
  onShuffleChange: (shuffle: boolean) => void
  onStart: (quiz: Quiz, label: string, questions: Question[]) => void
}

const Home = ({ progress, shuffle, onShuffleChange, onStart }: HomeProps) => (
  <main className={styles.home}>
    <header className={styles.header}>
      <h1>Recall</h1>
      <p>Pick a quiz, a section or a single topic, then choose an answer for each question.</p>
      <div className={styles.key}>
        <span>Colour key:</span>
        <Legend />
      </div>
      <label className={styles.shuffle}>
        <input type="checkbox" checked={shuffle} onChange={(event) => onShuffleChange(event.target.checked)} />
        Shuffle questions
      </label>
    </header>

    <div className={styles.decks}>
      {quizzes.map((quiz) => {
        const counts = tally(quiz.questions, progress)
        const missed = quiz.questions.filter((q) => progress[q.id] === 'missed')
        return (
          <section key={quiz.id} className={styles.deck} style={deckStyle(quiz.id)}>
            <div className={styles.deckHeader}>
              <div>
                <h2>{quiz.title}</h2>
                <span className={styles.source}>{quiz.source} · {quiz.questions.length} questions</span>
              </div>
              <div className={styles.deckActions}>
                <button type="button" className={styles.primary} onClick={() => onStart(quiz, 'All questions', quiz.questions)}>Start all</button>
                <button type="button" className={styles.review} disabled={missed.length === 0} onClick={() => onStart(quiz, 'Missed questions', missed)}>
                  Redo missed ({missed.length})
                </button>
              </div>
            </div>
            <MasteryBar {...counts} />
            <div className={styles.legend}>
              <span className={styles.gotDot}>{counts.got} got</span>
              <span className={styles.missedDot}>{counts.missed} missed</span>
              <span className={styles.unseenDot}>{counts.unseen} not tried</span>
            </div>

            {groupBy(quiz.questions, (q) => q.section).map(([section, sectionQuestions]) => (
              <div key={section} className={styles.section}>
                <button type="button" className={styles.sectionButton} onClick={() => onStart(quiz, section, sectionQuestions)}>
                  <Markdown inline quiz={quiz.id} text={section} />
                  <span className={styles.count}>{sectionQuestions.length}</span>
                </button>
                <div className={styles.topics}>
                  {groupBy(sectionQuestions, (q) => q.topic)
                    .filter(([topic]) => topic !== section)
                    .map(([topic, topicQuestions]) => {
                      const topicCounts = tally(topicQuestions, progress)
                      return (
                        <button key={topic} type="button" className={styles.topic} onClick={() => onStart(quiz, topic, topicQuestions)}>
                          <Markdown inline quiz={quiz.id} text={topic} />
                          <MasteryBar {...topicCounts} thin />
                        </button>
                      )
                    })}
                </div>
              </div>
            ))}
          </section>
        )
      })}
    </div>
  </main>
)

export default Home
