package main

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	EPOCH               = 1640991600000 // 1st January 2022, time start
	TIMESTAMP_BITS      = 41
	SHARD_BITS          = 5
	AUTO_SEQ_MOD_BITS   = 18
	MAX_SHARD_COUNT     = (1 << SHARD_BITS) - 1        // there can only be 2^13 number of shards
	AUT0_SEQ_MOD        = (1 << AUTO_SEQ_MOD_BITS) - 1 // sequence counters
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

	mu sync.Mutex
	lt uint64
}

type TimeGenerator interface {
	CustomEpoch() uint64
}

type timeGen struct{}

func (tg *timeGen) CustomEpoch() uint64 {
	return uint64(time.Now().UnixMilli() - EPOCH)
}

func (tg *IdGenerator) GetTimeStamp() (ts uint64, seq uint64) {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	seq = uint64(tg.autoIncSeq.Add(1) % int64(tg.seqMod))
	now := tg.timeGen.CustomEpoch()

	for now == tg.lt && seq == 0 {
		time.Sleep(1 * time.Millisecond)
		now = tg.timeGen.CustomEpoch()
		seq = uint64(tg.autoIncSeq.Add(1) % int64(tg.seqMod))
	}
	tg.lt = now

	return uint64(now), seq
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		autoIncSeq: atomic.Int64{},
		timeGen:    &timeGen{},
		seqMod:     AUT0_SEQ_MOD,
	}
}

func (g *IdGenerator) GenerateId(nid int) uint64 {
	if nid > MAX_SHARD_COUNT {
		// max shard count reached
		return 0
	}
	now, seq := g.GetTimeStamp()
	id := now << (TIMESTAMP_BIT_SHIFT)
	id |= uint64(nid << (SHARD_ID_BIT_SHIFT))
	id |= seq
	return id
}
