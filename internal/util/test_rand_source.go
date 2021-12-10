package util

import (
	"hash/fnv"
	"math/rand"
	"os"
	"testing"
)

const DefaultTestRandSeed = 0xCAFE

// TestRandSource creates a new rand.Source for use in tests and benchmarks.
// It uses a hash of the environment variable TEST_RAND_SEED as the seed if present, or DefaultTestRandSeed otherwise.
// If t is non-nil, its Name is hashed and mixed into the seed value, so that different tests will get different
// random values.
func TestRandSource(t testing.TB) rand.Source {
	var seed int64 = DefaultTestRandSeed

	if str, ok := os.LookupEnv("TEST_RAND_SEED"); ok {
		h := fnv.New64()
		_, err := h.Write([]byte(str))
		if err != nil {
			panic(err)
		}
		seed = int64(h.Sum64())
	}

	// mix in the test name so that different tests don't get identical sequences.
	if t != nil {
		h := fnv.New64()
		_, err := h.Write([]byte(t.Name()))
		if err != nil {
			panic(err)
		}
		seed = int64(h.Sum64() ^ uint64(seed))
	}

	return rand.NewSource(seed)
}
