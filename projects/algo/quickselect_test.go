package algo

import (
	"reflect"
	"sort"
	"testing"
)

func TestQuickSelectAgreesWithSort(t *testing.T) {
	nums := []int{7, 2, 9, 4, 1, 8, 3, 5, 6}
	sorted := append([]int(nil), nums...)
	sort.Ints(sorted)

	for k := 0; k < len(nums); k++ {
		if got := QuickSelect(nums, k); got != sorted[k] {
			t.Errorf("QuickSelect(nums, %d) = %d, want %d", k, got, sorted[k])
		}
	}
}

func TestQuickSelectDoesNotMutateInput(t *testing.T) {
	nums := []int{5, 3, 1, 4, 2}
	before := append([]int(nil), nums...)
	QuickSelect(nums, 2)
	if !reflect.DeepEqual(nums, before) {
		t.Errorf("QuickSelect mutated input: got %v, want %v", nums, before)
	}
}

func TestQuickSelectSingleElement(t *testing.T) {
	if got := QuickSelect([]int{42}, 0); got != 42 {
		t.Errorf("QuickSelect single element = %d, want 42", got)
	}
}
