# Algorithms notes

## Big-O, intuitively
- Describes growth as input n grows, not wall-clock time.
- `O(1)` hash lookup · `O(log n)` binary search · `O(n)` scan · `O(n log n)`
  good sorts · `O(n²)` nested loops · `O(2ⁿ)` brute-force subsets.
- Also reason about **space** complexity (recursion stack counts).

## Sorting
- Built-in sorts are `O(n log n)` (Timsort, introsort) — use them.
- Know *why*: merge sort = divide + merge (stable, O(n) extra); quicksort =
  partition around a pivot (in-place, worst-case O(n²) on bad pivots).

## Binary search (sorted input!)
```python
lo, hi = 0, len(a) - 1
while lo <= hi:
    mid = (lo + hi) // 2
    if a[mid] == target: return mid
    if a[mid] < target:  lo = mid + 1
    else:                hi = mid - 1
return -1
```
- Off-by-one bugs live here. Also useful: "search for the boundary" variant.

## Patterns worth recognizing
- **Two pointers** — pairs in a sorted array, in-place dedupe. O(n), O(1) space.
- **Sliding window** — longest/`k`-sized subarray problems.
- **Hash map for O(1) lookup** — two-sum, counting, dedupe.
- **BFS/DFS** — graph/tree traversal; BFS for shortest unweighted path.
- **Dynamic programming** — overlapping subproblems + optimal substructure;
  memoize (top-down) or tabulate (bottom-up).

## Mindset
- Clarify constraints first (n size, sorted?, duplicates?) — they pick the tool.
- Brute force, then optimize the bottleneck. Correct + clear beats clever.
