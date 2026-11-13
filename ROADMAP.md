# Roadmap

Work through these in order. Don't add several at once.

1. **CRUD** ✅
2. **HTTP fundamentals:** status codes, JSON encode/decode, validation, path/query params, consistent errors, Content-Type, request lifecycle, `context.Context` ✅
3. **Authentication / authorization:** users, password hashing (bcrypt or Argon2id), sessions + cookies (sessions before JWTs), auth middleware, ownership checks (403 vs 404), database constraints ✅
4. **Transactions:** a dedicated MoveCard operation that updates the card's column and writes a movement-history row in one transaction (`BeginTx` / `Commit` / `Rollback`) ✅
5. **Concurrency:** race conditions (two users moving the same card), `SELECT ... FOR UPDATE`, isolation levels, optimistic vs pessimistic concurrency, goroutines, mutexes, channels ✅
6. **Database tooling:** what `database/sql` boilerplate costs, then the alternatives: sqlx (maps SQL results into structs), sqlc (generates Go code from SQL), ORMs (generate SQL from Go code). Convert one store file to sqlc and compare. Plus **migrations**: numbered SQL files (e.g. with goose or golang-migrate) that build every database the same way, replacing the hand-run SQL in `README.md` and preventing schema drift. Plus **indexes and query cost**: Big O (a full table scan is O(n), an index lookup is O(log n)), why Postgres indexes are B-trees, reading `EXPLAIN ANALYZE`, and adding the missing `cards.column_id` index in a migration.
7. **Collaboration:** workspaces (teams) own boards; users join them through `workspace_members` with a role (`owner`, `editor`, `viewer`). Every user gets a personal workspace. Role-based access control: non-members get `404`, a viewer who tries to edit gets `403`. Invite and remove members. Move `boards.owner_id` to `boards.workspace_id` with a data migration. Two different people can now move the same card, so step 5's race becomes real.
8. **WebSockets:** real-time board updates, hub/broadcast pattern, connection lifecycle, heartbeats, disconnect cleanup, backpressure
9. **Background workers:** the queue data structure (first in, first out), then an in-process job queue (`chan Job` + worker goroutines) for audit events and notifications
10. **Caching / Redis**
11. **Observability:** structured logging with `log/slog`, request IDs, latency and query timing, error tracking, metrics
12. **Testing**
13. **Docker + deployment**
14. **Production architecture:** system design for scaling: load balancers, running several API servers (horizontal scaling, and why sessions in Postgres make that easy), database read replicas, and rate limiting. Also service-to-service communication with gRPC (HTTP/2, Protocol Buffers, `.proto` contracts) vs REST. Exercise: move the notification worker (step 9) into its own service and call it over gRPC.

## Schedule

Goal: finish by the end of April 2027, about two weeks per step. Step sizes vary, so dates will shift; April is buffer.

| Dates | Step |
| --- | --- |
| 21 Sep – 26 Oct 2026 | 1–3 ✅ |
| 27 Oct – 1 Nov | Revision: step 3, SQL, Go weak spots ✅ |
| 2 – 15 Nov | 4. Transactions ✅ |
| 16 – 29 Nov | 5. Concurrency ✅ |
| 30 Nov – 13 Dec | 6. Database tooling + migrations |
| 14 Dec – 3 Jan 2027 | 7. Collaboration (3 weeks, includes holidays) |
| 4 – 17 Jan | 8. WebSockets |
| 18 – 31 Jan | 9. Background workers |
| 1 – 14 Feb | 10. Caching / Redis |
| 15 – 28 Feb | 11. Observability |
| 1 – 14 Mar | 12. Testing |
| 15 – 28 Mar | 13. Docker + deployment |
| 29 Mar – 11 Apr | 14. Production architecture |
| 12 – 30 Apr | Buffer and full revision |
