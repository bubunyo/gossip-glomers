package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssertBitPlaces(t *testing.T) {
	assert.Equal(t, 64, TIMESTAMP_BITS+SHARD_BITS+AUTO_SEQ_MOD_BITS)
}

var checker = map[int64]int{}

func insertKey() chan int64 {
	c := make(chan int64, 100000)
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
	g := &IdGenerator{
		timeGen: &timeGenTester{},
		seqMod:  5,
	}
	c := insertKey()
	for i := 0; i < 5000; i++ {
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
		go func() { c <- g.GenerateId(0) }()
	}
	time.Sleep(5 * time.Second)
	close(c)
	printChecker()
}

type timeGenTester struct{ c int64 }

func (g *timeGenTester) CustomEpoch() int64 {
	g.c += 1
	return g.c
}

func TestIdGenerationLogic(t *testing.T) {

	ig := &IdGenerator{
		autoIncSeq: atomic.Int64{},
		timeGen:    &timeGenTester{},
		seqMod:     5,
	}

	for i := 0; i < 30; i++ {
		fmt.Println(">>>>>>>>>>>>", ig.GenerateId(MAX_SHARD_COUNT))
	}
}
