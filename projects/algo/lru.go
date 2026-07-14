package algo

// LRUCache is a fixed-capacity cache with O(1) Get/Put: a hash map for
// lookup plus a doubly linked list to track recency, so eviction never
// has to scan.
type LRUCache struct {
	capacity   int
	items      map[int]*lruNode
	head, tail *lruNode // head = most recently used, tail = least
}

type lruNode struct {
	key, val   int
	prev, next *lruNode
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{capacity: capacity, items: make(map[int]*lruNode)}
}

func (c *LRUCache) Get(key int) (int, bool) {
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.moveToFront(n)
	return n.val, true
}

func (c *LRUCache) Put(key, val int) {
	if n, ok := c.items[key]; ok {
		n.val = val
		c.moveToFront(n)
		return
	}
	if len(c.items) >= c.capacity && c.tail != nil {
		evict := c.tail
		c.remove(evict)
		delete(c.items, evict.key)
	}
	n := &lruNode{key: key, val: val}
	c.items[key] = n
	c.pushFront(n)
}

func (c *LRUCache) moveToFront(n *lruNode) {
	if c.head == n {
		return
	}
	c.remove(n)
	c.pushFront(n)
}

func (c *LRUCache) pushFront(n *lruNode) {
	n.prev, n.next = nil, c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

func (c *LRUCache) remove(n *lruNode) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
	n.prev, n.next = nil, nil
}
