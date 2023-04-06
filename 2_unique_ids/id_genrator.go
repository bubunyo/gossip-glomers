package main

import (
	"sync/atomic"
	"time"
)

const (
	EPOCH               = 1640991600000 // 1st January 2022, time start
	TIMESTAMP_BITS      = 41
	SHARD_BITS          = 13
	AUTO_SEQ_MOD_BITS   = 10
	MAX_SHARD_COUNT     = 1 << SHARD_BITS        // there can only be 2^13 number of shards
	AUT0_SEQ_MOD        = 1 << AUTO_SEQ_MOD_BITS // there can only be 2^10 number of shards
	TIMESTAMP_BIT_SHIFT = 64 - TIMESTAMP_BITS
	SHARD_ID_BIT_SHIFT  = TIMESTAMP_BIT_SHIFT - SHARD_BITS
)

var (
	autoIncSeq int64
)

/*
Ids are 64 bit long. and segregated into the following
11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111
|____________________________________________||____________||_________|
                      |                              |           |
               Custom Timestamp                   ShardId   Auto Inc Seq No

Custom Timestamp: [41 bits] This is the number of milliseconds since EPOCH.
ShardId: [13 bits] Shard bucket of the server
    Shard id is genrated by apply 2^13 to server number
Auto Inc Seq No: [10 bits] Ever increase counter per ig generator instance.
    Mod 2^10 is applied onto the number to make it fit into 10 bits
*/

type IdGenerator struct {
	// nid        int
	autoIncSeq uint64
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		// nid: nid + 1
	}
}

func (g *IdGenerator) GenerateId(nid int) uint64 {
	now := uint64(time.Now().UnixMilli() - EPOCH)
	id := now << (TIMESTAMP_BIT_SHIFT)
	shardId := int64(nid % MAX_SHARD_COUNT)
	id |= uint64(shardId << (SHARD_ID_BIT_SHIFT))
	seq := atomic.AddUint64(&g.autoIncSeq, 1)
	// log.Println(">>>>>>>>>>>>>", seq)
	id |= uint64(seq % AUT0_SEQ_MOD)
	return id
}
