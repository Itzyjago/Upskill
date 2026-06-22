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
