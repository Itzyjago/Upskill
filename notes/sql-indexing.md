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
