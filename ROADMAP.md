# Roadmap

Work through these in order. Don't add several at once.

1. **CRUD** ✅
2. **HTTP fundamentals:** status codes, JSON encode/decode, validation, path/query params, consistent errors, Content-Type, request lifecycle, `context.Context` ✅
3. **Authentication / authorization:** users, password hashing (bcrypt or Argon2id), sessions + cookies (sessions before JWTs), auth middleware, ownership checks (403 vs 404), database constraints ✅
4. **Transactions:** a dedicated MoveCard operation that updates the card's column and writes a movement-history row in one transaction (`BeginTx` / `Commit` / `Rollback`)
5. **Concurrency:** race conditions (two users moving the same card), `SELECT ... FOR UPDATE`, isolation levels, optimistic vs pessimistic concurrency, goroutines, mutexes, channels
6. **Database tooling:** what `database/sql` boilerplate costs, then the alternatives: sqlx (maps SQL results into structs), sqlc (generates Go code from SQL), ORMs (generate SQL from Go code). Convert one store file to sqlc and compare. Plus **migrations**: numbered SQL files (e.g. with goose or golang-migrate) that build every database the same way, replacing the hand-run SQL in `README.md` and preventing schema drift.
7. **WebSockets:** real-time board updates, hub/broadcast pattern, connection lifecycle, heartbeats, disconnect cleanup, backpressure
8. **Background workers:** in-process job queue (`chan Job` + worker goroutines) for audit events and notifications
9. **Caching / Redis**
10. **Observability:** structured logging with `log/slog`, request IDs, latency and query timing, error tracking, metrics
11. **Testing**
12. **Docker + deployment**
13. **Production architecture:** includes service-to-service communication with gRPC (HTTP/2, Protocol Buffers, `.proto` contracts) vs REST. Exercise: move the notification worker (step 8) into its own service and call it over gRPC.

## Schedule

Goal: finish by the end of April 2027, about two weeks per step. Step sizes vary, so dates will shift; April is buffer.

| Dates | Step |
| --- | --- |
| 21 Sep – 26 Oct 2026 | 1–3 ✅ |
| 27 Oct – 1 Nov | Revision: step 3, SQL, Go weak spots |
| 2 – 15 Nov | 4. Transactions |
| 16 – 29 Nov | 5. Concurrency |
| 30 Nov – 13 Dec | 6. Database tooling + migrations |
| 14 Dec – 3 Jan 2027 | 7. WebSockets (3 weeks, includes holidays) |
| 4 – 17 Jan | 8. Background workers |
| 18 – 31 Jan | 9. Caching / Redis |
| 1 – 14 Feb | 10. Observability |
| 15 – 28 Feb | 11. Testing |
| 1 – 14 Mar | 12. Docker + deployment |
| 15 – 28 Mar | 13. Production architecture |
| April | Buffer and full revision |
