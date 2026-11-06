import styles from './Legend.module.css'

// What each highlight colour means. Shown wherever highlighted text appears.
const Legend = () => (
  <div className={styles.legend} aria-label="Colour key">
    <code className="kind-http">HTTP</code>
    <code className="kind-sql">SQL</code>
    <code className="kind-go">Go</code>
    <strong>Key idea</strong>
  </div>
)

export default Legend
