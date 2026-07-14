package algo

import "testing"

func TestConsistentHashRingAssignsEveryKey(t *testing.T) {
	r := NewConsistentHashRing(10)
	r.AddNode("nodeA")
	r.AddNode("nodeB")
	r.AddNode("nodeC")

	for i := 0; i < 100; i++ {
		key := "key-" + string(rune('a'+i%26)) + string(rune(i))
		if _, ok := r.GetNode(key); !ok {
			t.Fatalf("GetNode(%q) ok=false, want true (ring has nodes)", key)
		}
	}
}

func TestConsistentHashRingEmptyRing(t *testing.T) {
	r := NewConsistentHashRing(5)
	if _, ok := r.GetNode("anything"); ok {
		t.Error("GetNode on empty ring should return ok=false")
	}
}

func TestConsistentHashRingStableAssignment(t *testing.T) {
	// The whole point of consistent hashing: the same key maps to the same
	// node on repeated lookups, and removing an unrelated node doesn't
	// change assignments for keys that weren't on that node.
	r := NewConsistentHashRing(20)
	r.AddNode("nodeA")
	r.AddNode("nodeB")
	r.AddNode("nodeC")

	keys := make([]string, 50)
	before := make(map[string]string, 50)
	for i := range keys {
		keys[i] = "item-" + string(rune('a'+i%26)) + string(rune(i))
		node, _ := r.GetNode(keys[i])
		before[keys[i]] = node
	}

	r.AddNode("nodeD")

	changed := 0
	for _, k := range keys {
		node, _ := r.GetNode(k)
		if node != before[k] {
			changed++
		}
	}
	// Adding a 4th of 4 nodes should remap roughly 1/4 of keys, not all of
	// them — assert well under 100%, the naive mod-n behavior this exists
	// to avoid.
	if changed > len(keys)*3/4 {
		t.Errorf("changed = %d/%d after adding one node, want well under 100%% (consistent hashing's whole point)", changed, len(keys))
	}
}
