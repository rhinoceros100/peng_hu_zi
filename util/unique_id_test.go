package util

import "testing"

func TestUniqueId(t *testing.T) {
	ids := map[uint64]bool{}
	for i := 0; i < 100000; i++ {
		id := UniqueId()
		if _, ok := ids[id]; ok {
			t.Fatal("it should not same", i, id)
		}
		ids[id] = true
	}
}
