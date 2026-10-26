import { useQuery, useQueryClient } from '@tanstack/react-query'
import { api, type User } from '../api'
import App from '../App'
import LoginPage from '../LoginPage'
import styles from './AuthGate.module.css'

// AuthGate asks the API who is logged in (GET /me) and shows the login page
// until there's a user.
const AuthGate = () => {
  const queryClient = useQueryClient()
  const meQuery = useQuery({ queryKey: ['me'], queryFn: api.me, retry: false })

  const onLogin = (user: User) => {
    queryClient.setQueryData(['me'], user)
  }

  const onLogout = async () => {
    await api.logout()
    // Update ['me'] in place first, so AuthGate re-renders with the login page.
    // (clear() would discard the query this component is subscribed to.)
    queryClient.setQueryData(['me'], null)
    // Then forget the previous user's boards, so the next user never sees them.
    queryClient.removeQueries({ predicate: (query) => query.queryKey[0] !== 'me' })
  }

  if (meQuery.isLoading) {
    return <div className={styles.center}>Loading…</div>
  }
  if (meQuery.isError) {
    return <div className={styles.center}>Couldn’t reach the API. Make sure it’s running on port 8080.</div>
  }
  if (!meQuery.data) {
    return <LoginPage onLogin={onLogin} />
  }
  return <App user={meQuery.data} onLogout={onLogout} />
}

export default AuthGate
