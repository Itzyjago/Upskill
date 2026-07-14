package algo

import "testing"

func TestStack(t *testing.T) {
	var s Stack[int]
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on empty stack should return ok=false")
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if got := s.Len(); got != 3 {
		t.Fatalf("Len() = %d, want 3", got)
	}
	if v, ok := s.Peek(); !ok || v != 3 {
		t.Fatalf("Peek() = %d, %v, want 3, true", v, ok)
	}
	for _, want := range []int{3, 2, 1} {
		v, ok := s.Pop()
		if !ok || v != want {
			t.Fatalf("Pop() = %d, %v, want %d, true", v, ok, want)
		}
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop after draining should return ok=false")
	}
}
