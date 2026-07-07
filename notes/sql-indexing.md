# SQL indexing cheat sheet

Built from reading real `EXPLAIN ANALYZE` output. Companion to [sql.md](sql.md).

## Read the plan bottom-up
- The innermost / most-indented node runs first; rows flow upward.
- `Seq Scan` = full table read. Fine for tiny tables, a smell on big ones.
- `Index Scan` / `Index Only Scan` = the index did its job.
- Compare **estimated** vs **actual** rows — a big gap means stale statistics
  (`ANALYZE` the table) or a bad index choice.

## When an index actually helps
- High-selectivity filters: `WHERE email = ?` (few matching rows).
- Join keys and `ORDER BY` / `GROUP BY` columns.
- Not helpful when you read most of the table anyway — a seq scan is cheaper.

## Composite index column order matters
```sql
CREATE INDEX idx_orders_cust_date ON orders (customer_id, created_at);
```
- Works for `WHERE customer_id = ?` and `WHERE customer_id = ? AND created_at > ?`.
- Does **not** help a query filtering only on `created_at` — leftmost-prefix rule.

## Things that quietly kill index use
- Functions on the column: `WHERE LOWER(email) = ?` → add a functional index or
  store normalized.
- Leading wildcard: `LIKE '%foo'` can't use a B-tree.
- Implicit type casts: comparing a `text` column to an integer literal.

## Covering indexes
- `INCLUDE` non-key columns so the query is answered from the index alone
  (`Index Only Scan`), skipping the heap fetch.

## Composite indexes, the follow-up: equality before range
Leftmost-prefix says *which* columns matter; it doesn't say what order to put
them in when a query mixes an equality filter with a range one — that's a
separate rule, and getting it backwards quietly gives up most of the index's
value instead of failing outright.
- **Equality columns go first, range columns last.** Same index as above,
  `(customer_id, created_at)`, against
  `WHERE customer_id = ? AND created_at > ?`: Postgres narrows to exactly one
  `customer_id`'s worth of the index (an equality match collapses to a single
  point), then walks a contiguous range of `created_at` *within* that
  slice — one index range scan, cheap.
- **Flip the order** — `(created_at, customer_id)` — and the same query can
  only use `created_at`'s range to narrow the index; `customer_id` inside
  that range gets filtered row-by-row after the fact (still shows as "uses
  the index," but the plan's `Rows Removed by Filter` is where the lost
  selectivity shows up). The index isn't *unused*, it's just doing much less
  work than it looks like it's doing.
- **Why this matters beyond one query**: an index also serves `ORDER BY` for
  free if the trailing columns match the sort *and* every column before them
  in the index is pinned by an equality filter — `(customer_id, created_at)`
  answers `WHERE customer_id = ? ORDER BY created_at` without a separate
  sort step. Get the equality-then-range order right and this falls out for
  free; get it backwards and neither the filter nor the sort work as well.
- **Don't over-index the fix**: it's tempting to add every filterable column
  to one composite index. Each additional column makes the index wider (more
  I/O per lookup) and every `INSERT`/`UPDATE` pays to maintain it — a
  three-or-four-column composite index earns its keep for one specific hot
  query shape; it's not a substitute for choosing which query shapes
  actually need one.
