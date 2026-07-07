# Data structures notes

## Hash maps
- Average O(1) insert/lookup by hashing the key to a bucket.
- Collisions handled by chaining (linked list/tree per bucket) or open
  addressing (probe for the next slot).
- Worst case O(n) if many keys collide; a good hash + resizing keeps it rare.
- Unordered by design — don't rely on iteration order (Python dict insertion
  order is a CPython guarantee, not a hash-map property in general).

## Trees
- **Binary search tree**: left < node < right → O(log n) search *when balanced*,
  O(n) when it degenerates into a list.
- **Self-balancing** (AVL, red-black) keeps height ~log n on every insert/delete;
  most "ordered map" implementations are red-black trees.
- **B-tree / B+ tree**: wide, shallow — fewer disk reads, which is why databases
  use them for indexes (see [sql-indexing.md](sql-indexing.md)).
- **Heap**: not for searching; gives O(1) peek-min/max and O(log n) push/pop —
  the backbone of priority queues.

## Picking one
| Need | Reach for |
|------|-----------|
| Fast key lookup, order doesn't matter | Hash map |
| Keys kept sorted / range queries | Balanced BST / B-tree |
| "Give me the smallest next" repeatedly | Heap |
| Membership at scale, OK with rare false positives | Bloom filter |

## Mindset
- Big-O is about *growth*, not raw speed. For small n, a flat array often beats
  a "smarter" structure thanks to cache locality.

## A real hash map, not a kata (wordcount's `metrics.go`)
This sat as an abstract 🟡 for too long — the fix wasn't reading more theory,
it was noticing this repo already has a hash map doing real work.
`metrics` (`metrics.go`) is a Prometheus registry keyed by `labelKey{method,
path, status}` — `map[labelKey]int64`, `map[labelKey][]int64`, and so on.
Two things fall directly out of the "average O(1), unordered by design" bullet
above, made concrete:
- **A struct as a map key** works because Go structs are comparable
  (hashable) when every field is — three strings, no pointers/slices, so
  `labelKey` is a valid key with zero extra work. This is *why* the registry
  reaches for a map instead of, say, a slice of structs scanned linearly on
  every `observe()` call: O(1) amortized insert/update per request instead
  of an O(n) scan through every label combination seen so far.
- **"Unordered by design" is not academic here** — it's a bug that got
  fixed. `render()` (the `/metrics` exposition endpoint) iterates the
  registry, and Go's map iteration order is *randomized per run*, not just
  "unspecified once." Without correcting for it, `/metrics` would return
  the same data in a different byte order on every scrape — cosmetically
  fine for Prometheus (it parses by label, not position) but bad for anyone
  diffing scrape output or writing a golden-file test against it.
  `sortedKeys()` collects the map's keys into a slice and sorts them (path,
  then method, then status) before rendering — turning an unordered
  structure into a deterministic *view* of it, without changing the
  underlying map or paying map-ordering cost on every single insert.
- The general lesson: reach for a hash map when lookup/update by key is the
  hot path (`observe()` runs on every request); reach for a sorted slice (or
  re-sort on read, like here) only where something downstream — a
  human, a diff, a test — actually needs the order, and only at the point
  that needs it.
