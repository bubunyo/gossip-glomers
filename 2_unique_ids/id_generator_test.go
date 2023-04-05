package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertBitPlaces(t *testing.T) {
	assert.Equal(t, 64, TIMESTAMP_BITS+SHARD_BITS+AUTO_SEQ_MOD_BITS)
}

// func TestUniqueId(t *testing.T) {
// 	g := NewIdGenerator()
// 	for i := 0; i < 100000; i++ {
// 		id := g.GenerateId(i % 10)
// 		fmt.Printf(">>> id=%d - mod=%d\n", id, i%10)
// 	}
// }
