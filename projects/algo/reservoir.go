package algo

import "math/rand"

// ReservoirSample picks k items uniformly at random from stream (a slice
// standing in for a stream too large to hold twice), reading it exactly
// once. Algorithm R: fill the reservoir with the first k, then for each
// later item at index i, replace a random reservoir slot with probability
// k/(i+1) — the classic trick for sampling from a stream of unknown length
// without buffering it all.
func ReservoirSample(stream []int, k int, rng *rand.Rand) []int {
	if k <= 0 {
		return nil
	}
	reservoir := make([]int, 0, k)
	for i, v := range stream {
		if i < k {
			reservoir = append(reservoir, v)
			continue
		}
		j := rng.Intn(i + 1)
		if j < k {
			reservoir[j] = v
		}
	}
	return reservoir
}
