import CardEditor from '../CardEditor'
import ColumnView from '../ColumnView'
import BoardToolbar from '../BoardToolbar'
import EmptyWorkspace from '../EmptyWorkspace'
import Icon from '../Icon'
import Modal from '../Modal'
import TitleEditor from '../TitleEditor'
import { useBoardWorkspace } from './AppHook'
import type { User } from '../api'
import styles from './App.module.css'

type AppProps = { user: User; onLogout: () => void }

const App = ({ user, onLogout }: AppProps) => {
  const workspace = useBoardWorkspace()
  const {
    board, boardQuery, boards, boardsQuery, cardCount, closeModal, columns,
    deleteCard, draggingCardId, editingCard, modal, moveCard, mutation, notice,
    notify, openCard, saveCard, saveTitle, setActiveId, setDraggingCardId, setModal,
  } = workspace

  return (
    <main className={styles.appShell}>
      <aside className={styles.sidebar}>
        <a className={styles.brand} href="#home" aria-label="Boardroom home">
          <span className={styles.brandMark}><span /><span /><span /><span /></span>
          <span>boardroom<span className={styles.brandDot}>.</span></span>
        </a>
        <div className={styles.workspaceLabel}>WORKSPACE</div>
        <button className={`${styles.navItem} ${styles.navItemActive}`} type="button">
          <Icon>▦</Icon><span>My boards</span><span className={styles.navCount}>{boards.length}</span>
        </button>
        <div className={styles.sidebarSection}>
          <div className={styles.sidebarHeading}>
            <span>YOUR BOARDS</span>
            <button type="button" className={styles.tinyButton} onClick={() => setModal('board')} aria-label="Create a board">＋</button>
          </div>
          <div className={styles.boardList}>
            {boards.map((item, index) => (
              <button
                key={item.id}
                type="button"
                onClick={() => setActiveId(item.id)}
                className={`${styles.boardLink} ${board?.id === item.id ? styles.boardLinkActive : ''}`}
              >
                <span className={`${styles.boardGlyph} ${styles[`glyph${index % 4}`]}`}>{['◈', '◉', '◇', '▧'][index % 4]}</span>
                <span className={styles.boardName}>{item.title}</span>
                {board?.id === item.id && <span className={styles.activeDot} />}
              </button>
            ))}
            {boards.length === 0 && !boardsQuery.isLoading && <p className={styles.sidebarEmpty}>Your boards will show up here.</p>}
          </div>
        </div>
        <div className={styles.sidebarBottom}>
          <div className={styles.avatar}>{user.name.slice(0, 2).toUpperCase()}</div>
          <div className={styles.userInfo}><strong>{user.name}</strong><span>{user.email}</span></div>
          <button type="button" className={styles.logoutButton} onClick={onLogout}>Log out</button>
        </div>
      </aside>

      <section className={styles.mainArea}>
        <header className={styles.topbar}>
          <div className={styles.breadcrumb}>
            <span>Workspace</span><span className={styles.crumbSlash}>/</span><strong>{board?.title ?? 'My boards'}</strong>
          </div>
          <div className={styles.topActions}>
            <span className={styles.savedIndicator}><i /> All changes saved</span>
            <button className={styles.helpButton} type="button" aria-label="Help">?</button>
          </div>
        </header>

        {boardsQuery.isLoading ? (
          <div className={styles.centerState}><span className={styles.spinner} />Loading your workspace…</div>
        ) : boardsQuery.isError ? (
          <div className={styles.centerState}>
            <div className={styles.stateIcon}>!</div><h2>Couldn’t load your boards</h2>
            <p>Make sure the Go API is running on port 8080, then try again.</p>
            <button className={styles.primaryButton} onClick={() => boardsQuery.refetch()}>Try again</button>
          </div>
        ) : !board ? (
          <EmptyWorkspace onCreate={() => setModal('board')} />
        ) : (
          <div className={styles.boardPage}>
            <div className={styles.boardHeader}>
              <div>
                <div className={styles.eyebrow}>YOUR WORKSPACE <span>·</span> BOARD</div>
                <h1>{board.title}</h1>
                <p className={styles.boardSubtitle}>A little progress each day adds up to big results.</p>
              </div>
              <div className={styles.boardControls}>
                <div className={styles.memberStack}><span className={styles.memberOne}>M</span><span className={styles.memberTwo}>A</span><span className={styles.memberCount}>+2</span></div>
                <button className={styles.shareButton} type="button" onClick={() => {
                  navigator.clipboard?.writeText(window.location.href)
                  notify('Board link copied')
                }}><Icon>↗</Icon> Share</button>
                <button className={styles.moreButton} type="button" aria-label="More board actions">···</button>
              </div>
            </div>
            <BoardToolbar cardCount={cardCount} onAddColumn={() => setModal('column')} />
            {boardQuery.isError && <div className={styles.inlineError}>Could not refresh board details. Showing the latest board list data.</div>}
            <div className={styles.columnsViewport}>
              <div className={styles.columnsGrid}>
                {columns.length > 0 ? columns.slice().sort((a, b) => a.position - b.position).map((column, index) => (
                  <ColumnView
                    key={column.id}
                    column={column}
                    index={index}
                    draggingId={draggingCardId}
                    onAdd={() => openCard(column.id)}
                    onEdit={(card) => openCard(column.id, card)}
                    onDelete={deleteCard}
                    onDragStart={setDraggingCardId}
                    onDrop={(cardId) => moveCard(cardId, column.id)}
                  />
                )) : (
                  <div className={styles.noColumns}>
                    <span>＋</span><h3>Start with a column</h3>
                    <p>Break your board into steps like To do, In progress, and Done.</p>
                    <button className={styles.outlineButton} onClick={() => setModal('column')}>Add your first column</button>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </section>

      {notice && <div className={styles.toast} role="status">{notice}<button onClick={() => notify('')} aria-label="Dismiss">×</button></div>}
      {modal && board && (
        <Modal
          title={modal === 'card' ? (editingCard ? 'Edit task' : 'Add a task') : modal === 'column' ? 'Add a column' : 'Create a board'}
          onClose={closeModal}
        >
          {modal === 'card' ? (
            <CardEditor key={editingCard?.id ?? 'new-card'} initial={editingCard} busy={mutation.isPending} onSubmit={saveCard} onCancel={closeModal} />
          ) : (
            <TitleEditor
              key={modal}
              label={modal === 'column' ? 'Column name' : 'Board name'}
              placeholder={modal === 'column' ? 'e.g. In progress' : 'e.g. Product launch'}
              busy={mutation.isPending}
              onSubmit={saveTitle}
              onCancel={closeModal}
            />
          )}
        </Modal>
      )}
      {modal && !board && modal === 'board' && (
        <Modal title="Create a board" onClose={closeModal}>
          <TitleEditor label="Board name" placeholder="e.g. Product launch" busy={mutation.isPending} onSubmit={saveTitle} onCancel={closeModal} />
        </Modal>
      )}
    </main>
  )
}

export default App
