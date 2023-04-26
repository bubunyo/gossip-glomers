package main

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertBitPlaces(t *testing.T) {
	assert.Equal(t, 64, TIMESTAMP_BITS+SHARD_BITS+AUTO_SEQ_MOD_BITS)
}

var checker = map[uint64]int{}

func insertKey() chan uint64 {
	c := make(chan uint64, 100000)
	go func() {
		for v := range c {
			if _, ok := checker[v]; ok {
				checker[v] += 1
			} else {
				checker[v] = 1
			}
		}
	}()
	return c
}

func printChecker() {
	for i, v := range checker {
		if v > 1 {
			fmt.Println(i, "->", v)
		}
	}
}

func TestUniqueId(t *testing.T) {
	g := NewIdGenerator()
	c := insertKey()
	for i := 0; i < 5000000; i++ {
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
	}
	close(c)
	printChecker()
}

type timeGenTester struct{}

func (g *timeGenTester) CustomEpoch() uint64 {
	return 1257894000000
}

func TestIdGenerationLogic(t *testing.T) {

	ig := &IdGenerator{
		autoIncSeq: atomic.Int64{},
		timeGen:    &timeGenTester{},
		seqMod:     AUT0_SEQ_MOD,
	}

	for i := 0; i < 1; i++ {
		fmt.Println(">>>>>>>>>>>>", ig.GenerateId(MAX_SHARD_COUNT))
	}
}
