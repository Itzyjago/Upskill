package algo

import "testing"

func TestKnapsack01(t *testing.T) {
	// Classic textbook case: items (weight, value) = (1,1),(3,4),(4,5),(5,7),
	// capacity 7 -> best is items 2+3 (weight 3+4=7, value 4+5=9).
	weights := []int{1, 3, 4, 5}
	values := []int{1, 4, 5, 7}
	if got := Knapsack01(weights, values, 7); got != 9 {
		t.Errorf("Knapsack01(...) = %d, want 9", got)
	}
}

func TestKnapsack01ZeroCapacity(t *testing.T) {
	if got := Knapsack01([]int{1, 2}, []int{10, 20}, 0); got != 0 {
		t.Errorf("Knapsack01 with capacity 0 = %d, want 0", got)
	}
}

func TestKnapsack01SingleItemFits(t *testing.T) {
	if got := Knapsack01([]int{5}, []int{100}, 5); got != 100 {
		t.Errorf("Knapsack01 single fitting item = %d, want 100", got)
	}
}

func TestKnapsack01SingleItemTooHeavy(t *testing.T) {
	if got := Knapsack01([]int{10}, []int{100}, 5); got != 0 {
		t.Errorf("Knapsack01 single too-heavy item = %d, want 0", got)
	}
}
