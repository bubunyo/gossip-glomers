package main

import (
	"log"
	"os"
	"sync/atomic"
	"time"
)

const (
	EPOCH               = 1640991600000 // 1st January 2022, time start
	TIMESTAMP_BITS      = 41
	SHARD_BITS          = 5
	AUTO_SEQ_MOD_BITS   = 18
	MAX_SHARD_COUNT     = (1 << SHARD_BITS) - 1  // there can only be 2^13 number of shards
	AUT0_SEQ_MOD        = 1 << AUTO_SEQ_MOD_BITS // sequence counters
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
	seqMod     int64
	timeGen    TimeGenerator

	lt int64
}

type TimeGenerator interface {
	CustomEpoch() int64
}

type timeGen struct{}

func (tg *timeGen) CustomEpoch() int64 {
	return time.Now().UnixMilli() - EPOCH
}

func (g *IdGenerator) GetTimeStamp() (ts int64, seq int64) {
	seq = g.autoIncSeq.Add(1) % g.seqMod
	now := g.timeGen.CustomEpoch()

	for now == g.lt && seq == 0 {
		time.Sleep(time.Millisecond)
		now = g.timeGen.CustomEpoch()
	}
	g.lt = now

	return now, seq
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		autoIncSeq: atomic.Int64{},
		timeGen:    &timeGen{},
		seqMod:     AUT0_SEQ_MOD,
	}
}

func (g *IdGenerator) GenerateId(nid int) int64 {
	f, err := os.OpenFile("test.log", os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	if nid > MAX_SHARD_COUNT {
		// max shard count reached
		return 0
	}
	now, seq := g.GetTimeStamp()
	id := now << TIMESTAMP_BIT_SHIFT
	id |= int64(nid << SHARD_ID_BIT_SHIFT)
	id |= seq
	return id
}
