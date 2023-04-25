package main

import (
	"fmt"
	"sync"
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

type timeGenTester struct {
	m       sync.Mutex
	counter uint64
}

func (g *timeGenTester) CustomEpoch() uint64 {
	g.m.Lock()
	g.m.Unlock()
	n := g.counter
	if n%2 == 0 {
		n -= 1
	}
	// fmt.Println(">>>>>>>>>>> nnnn", n)
	g.counter += 1
	return n
}

func TestIdGenerationLogic(t *testing.T) {

	ig := &IdGenerator{
		autoIncSeq: 0,
		timeGen:    &timeGenTester{counter: 0},
		seqMod:     3,
	}

	for i := 0; i < 10; i++ {
		fmt.Println(">>>>>>>>>>>>", ig.GenerateId(0))
	}
}
