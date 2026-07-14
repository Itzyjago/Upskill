package algo

import (
	"math/rand"
	"testing"
)

func TestReservoirSampleSizeAndMembership(t *testing.T) {
	stream := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rng := rand.New(rand.NewSource(42))
	sample := ReservoirSample(stream, 3, rng)

	if len(sample) != 3 {
		t.Fatalf("len(sample) = %d, want 3", len(sample))
	}
	inStream := make(map[int]bool, len(stream))
	for _, v := range stream {
		inStream[v] = true
	}
	for _, v := range sample {
		if !inStream[v] {
			t.Errorf("sample contains %d, not present in the original stream", v)
		}
	}
}

func TestReservoirSampleKLargerThanStream(t *testing.T) {
	stream := []int{1, 2, 3}
	rng := rand.New(rand.NewSource(1))
	sample := ReservoirSample(stream, 10, rng)
	if len(sample) != len(stream) {
		t.Errorf("len(sample) = %d, want %d (k > len(stream) should just return everything)", len(sample), len(stream))
	}
}

func TestReservoirSampleRoughlyUniform(t *testing.T) {
	// Every index of a 10-item stream should end up in a size-1 sample
	// with roughly equal frequency over many trials — a real statistical
	// check of Algorithm R's core guarantee, not just "it returns k items".
	const n, trials = 10, 20000
	stream := make([]int, n)
	for i := range stream {
		stream[i] = i
	}
	counts := make([]int, n)
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < trials; i++ {
		sample := ReservoirSample(stream, 1, rng)
		counts[sample[0]]++
	}
	expected := float64(trials) / float64(n)
	for i, c := range counts {
		if deviation := float64(c) - expected; deviation < -expected*0.25 || deviation > expected*0.25 {
			t.Errorf("index %d picked %d times, expected ~%.0f (+/-25%%)", i, c, expected)
		}
	}
}
