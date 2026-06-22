# SQL notes

## Joins
- `INNER JOIN` — rows that match in both tables.
- `LEFT JOIN` — all left rows, NULLs where the right has no match.
- A `LEFT JOIN` + `WHERE right.id IS NULL` finds left rows with no match.

## Filtering vs grouping
```sql
SELECT customer_id, COUNT(*) AS orders
FROM   orders
WHERE  status = 'paid'      -- filters rows BEFORE grouping
GROUP  BY customer_id
HAVING COUNT(*) > 5;        -- filters groups AFTER grouping
```

## Indexes
- An index speeds up reads on the columns you filter/join/sort by.
- It costs write speed and disk — index what you query, not everything.
- Run `EXPLAIN` (or `EXPLAIN ANALYZE`) to see if a query uses the index or
  falls back to a full table scan.

## Gotchas
- `NULL` is not equal to anything, including `NULL`. Use `IS NULL`.
- Wrapping an indexed column in a function (`WHERE LOWER(email) = ...`) usually
  disables the index — store/normalize the value instead.
