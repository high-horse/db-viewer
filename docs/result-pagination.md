# Result pagination

Supported SELECT queries and read-only CTE shapes use counted cursor pagination.
The backend runs `COUNT(*)` over the original result to return an exact row and
page count on the first response. User filters, grouping, and explicit limits
remain part of that result. A conservative capability check excludes commands,
data-modifying CTEs, locking reads, multiple statements, and unsupported quoting.
It is not a security boundary or a guarantee that user-defined functions have no
side effects; counted navigation is intended for repeatable reads.

The initial response contains 100 rows. Next fetches another batch from the same
`database/sql.Rows` stream, with one row of lookahead. Ten recent pages are cached
per tab. Previous uses the cache or a direct offset query when necessary. Last
jumps to the final page with dialect-appropriate OFFSET (and the required
unbounded LIMIT in SQLite/MySQL), without fetching all intermediate pages through
the application. The resulting cursor continues at its absolute page position.

Clicking a column header cycles ascending, descending, and the original query
order; the context menu offers the same controls. Sorting wraps the original
result in an outer `ORDER BY` using a validated column ordinal and direction,
resets to page one, and clears stale cached pages. It sorts the complete query
result in the database, not just the displayed page. Internal column keys avoid
collisions with duplicate labels or real columns named `id` and `sn`.

The row-number gutter grows to fit the total number of digits. The footer shows
`Page X of Y`, total rows, and a separate `100 rows/page` batch-size label.

Non-pageable statements retain execute-once streaming, including mutations with
RETURNING. These results have no global-sort controls, derive their total only
when exhausted, and use a stoppable sequential fetch for Last. Stop takes effect
after the current batch so the stream remains resumable.

Exact counts and deep offsets can still be expensive inside the database. Counts
are refreshed on new query executions, sort changes, and offset jumps. Count and
data queries are separate statements, and cached pages do not form a transaction
snapshot: concurrent changes may affect totals or page contents. Use an explicit
unique ordering for stable navigation; sorting on non-unique values alone does
not guarantee tie order across re-executions.

Each live cursor pins a database connection. Cursors close on exhaustion, SQL or
scan failure, rerun, tab close, disconnect, connection replacement, application
shutdown, or after a ten-minute maximum lifetime. At most sixteen streams may be
open. Run the query again using the refresh button after expiry. Refresh is an
explicit re-execution, including for statements with side effects. Closing an
unconsumed result is cancellation, not a guarantee that its statement rolled back.

Opaque cursor tokens are scoped to the connection. Fetch requests include the next
page sequence; repeating the most recent fetch while its cursor is open returns
the cached batch instead of advancing again. A Wails request's completion does not
cancel the cursor; the first request's cancellation still interrupts query startup.

Application memory is bounded by page row counts, not total result size. Individual
large cells, driver buffers, and the database's own execution plan can still use
substantial memory. This is driver-stream pagination, not PostgreSQL `DECLARE` /
`FETCH` server-cursor SQL or keyset query rewriting. Streaming does not make expensive
sorts or scans cheap and may hold server resources until the cursor closes.

Validation includes SQLite integration tests for large streams, exact counts,
direct jumps and backwards navigation, global sorting, original query limits,
empty and exact page boundaries, mutation execution once, request lifetime,
replay, invalid tokens/pages, connection isolation, and cleanup. Browser checks
cover resizing, scrolling, multi-digit row numbers, counts, direct-page and sort
requests, and cursor Next navigation with simulated backend batches. Live
PostgreSQL and MySQL integration require running database servers.
