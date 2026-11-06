import styles from './MasteryBar.module.css'

type MasteryBarProps = {
  got: number
  missed: number
  unseen: number
  thin?: boolean
}

// Green for got it, red for missed, grey for not tried yet.
const MasteryBar = ({ got, missed, unseen, thin = false }: MasteryBarProps) => {
  const total = got + missed + unseen || 1
  return (
    <div className={`${styles.bar} ${thin ? styles.thin : ''}`} aria-label={`${got} got, ${missed} missed, ${unseen} not tried`}>
      <span className={styles.got} style={{ width: `${(got / total) * 100}%` }} />
      <span className={styles.missed} style={{ width: `${(missed / total) * 100}%` }} />
    </div>
  )
}

export default MasteryBar
