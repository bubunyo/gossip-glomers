package main

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	EPOCH               = 1640991600000 // 1st January 2022, time start
	TIMESTAMP_BITS      = 41
	SHARD_BITS          = 10
	AUTO_SEQ_MOD_BITS   = 13
	MAX_SHARD_COUNT     = 1 << SHARD_BITS        // there can only be 2^13 number of shards
	AUT0_SEQ_MOD        = 1 << AUTO_SEQ_MOD_BITS // there can only be 2^10 number of shards
	TIMESTAMP_BIT_SHIFT = 64 - TIMESTAMP_BITS
	SHARD_ID_BIT_SHIFT  = TIMESTAMP_BIT_SHIFT - SHARD_BITS
)

/*
Ids are 64 bit long. and segregated into the following
11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111
|____________________________________________||_________||____________|
                      |                              |           |
               Custom Timestamp                   ShardId   Auto Inc Seq No

Custom Timestamp: [41 bits] This is the number of milliseconds since EPOCH.
ShardId: [10 bits] Shard bucket of the server
    Shard id is genrated by apply 2^13 to server number
Auto Inc Seq No: [13 bits] Ever increase counter per ig generator instance.
    Mod 2^10 is applied onto the number to make it fit into 10 bits
*/

type IdGenerator struct {
	autoIncSeq atomic.Int64
	seqMod     uint64
	timeGen    TimeGenerator
}

type TimeGenerator interface {
	CustomEpoch() uint64
}

type timeGen struct {
	mu sync.Mutex
	lt int64
}

func (tg *timeGen) CustomEpoch() uint64 {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == tg.lt {
		time.Sleep(time.Millisecond)
		now = time.Now().UnixMilli()
	}
	tg.lt = now
	return uint64(now)
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		autoIncSeq: atomic.Int64{},
		timeGen:    &timeGen{},
		seqMod:     AUT0_SEQ_MOD,
	}
}

func (g *IdGenerator) GenerateId(nid int) uint64 {
	now := g.timeGen.CustomEpoch()
	id := now << (TIMESTAMP_BIT_SHIFT)
	shardId := int64(nid % MAX_SHARD_COUNT)
	id |= uint64(shardId << (SHARD_ID_BIT_SHIFT))
	id |= uint64(g.autoIncSeq.Add(1)) % g.seqMod
	g.autoIncSeq.Add(1)
	return id
}
