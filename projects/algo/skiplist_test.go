package algo

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSkipListInsertAndContains(t *testing.T) {
	sl := NewSkipList(rand.New(rand.NewSource(1)))
	vals := []int{5, 2, 8, 1, 9, 3}
	for _, v := range vals {
		sl.Insert(v)
	}
	for _, v := range vals {
		if !sl.Contains(v) {
			t.Errorf("Contains(%d) = false, want true", v)
		}
	}
	for _, v := range []int{0, 4, 6, 7, 10} {
		if sl.Contains(v) {
			t.Errorf("Contains(%d) = true, want false (never inserted)", v)
		}
	}
}

func TestSkipListToSliceIsSorted(t *testing.T) {
	sl := NewSkipList(rand.New(rand.NewSource(2)))
	for _, v := range []int{7, 3, 9, 1, 5, 2, 8, 4, 6} {
		sl.Insert(v)
	}
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if got := sl.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Errorf("ToSlice() = %v, want %v", got, want)
	}
}

func TestSkipListDuplicateInsertIgnored(t *testing.T) {
	sl := NewSkipList(rand.New(rand.NewSource(3)))
	sl.Insert(5)
	sl.Insert(5)
	sl.Insert(5)
	want := []int{5}
	if got := sl.ToSlice(); !reflect.DeepEqual(got, want) {
		t.Errorf("ToSlice() = %v, want %v (duplicates should not add nodes)", got, want)
	}
}

func TestSkipListManyInsertsStaySorted(t *testing.T) {
	sl := NewSkipList(rand.New(rand.NewSource(4)))
	seen := make(map[int]bool)
	for i := 0; i < 500; i++ {
		v := (i * 37) % 500 // a deterministic pseudo-shuffle, not sorted input
		if !seen[v] {
			seen[v] = true
			sl.Insert(v)
		}
	}
	got := sl.ToSlice()
	for i := 1; i < len(got); i++ {
		if got[i-1] >= got[i] {
			t.Fatalf("ToSlice() not sorted at index %d: %d >= %d", i, got[i-1], got[i])
		}
	}
	if len(got) != len(seen) {
		t.Errorf("len(ToSlice()) = %d, want %d", len(got), len(seen))
	}
}
