# Roadmap

Work through these in order. Don't add several at once.

1. **CRUD** (done)
2. **HTTP fundamentals:** status codes, JSON encode/decode, validation, path/query params, consistent errors, Content-Type, request lifecycle, `context.Context`
3. **Authentication / authorization:** users, password hashing (bcrypt or Argon2id), sessions + cookies (sessions before JWTs), auth middleware, ownership checks (403 vs 404), database constraints
4. **Transactions:** a dedicated MoveCard operation that updates the card's column and writes a movement-history row in one transaction (`BeginTx` / `Commit` / `Rollback`)
5. **Concurrency:** race conditions (two users moving the same card), `SELECT ... FOR UPDATE`, isolation levels, optimistic vs pessimistic concurrency, goroutines, mutexes, channels
6. **Database tooling:** what `database/sql` boilerplate costs, then the alternatives: sqlx (maps SQL results into structs), sqlc (generates Go code from SQL), ORMs (generate SQL from Go code). Convert one store file to sqlc and compare.
7. **WebSockets:** real-time board updates, hub/broadcast pattern, connection lifecycle, heartbeats, disconnect cleanup, backpressure
8. **Background workers:** in-process job queue (`chan Job` + worker goroutines) for audit events and notifications
9. **Caching / Redis**
10. **Observability:** structured logging with `log/slog`, request IDs, latency and query timing, error tracking, metrics
11. **Testing**
12. **Docker + deployment**
13. **Production architecture**
