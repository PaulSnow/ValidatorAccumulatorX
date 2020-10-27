package router

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/accumulator"
	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/database"
	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/node"
	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/types"
)

// The router is used to configure a set of accumulators to distribute the construction of merkle DAGs.  Here
// we are distributing the work over go routines, but in a production system we would distribute the load over
// a network.

type Router struct {
	EntryHashStream chan node.EntryHash        // Stream of hashes to record
	DBs             []*database.DB             // Databases where hashes are recorded
	ACCs            []*accumulator.Accumulator // Accumulators to record hashes
	EntryFeeds      []chan node.EntryHash
	Controls        []chan bool
	MDFeeds         []chan *types.Hash
}

func (r *Router) blockTimer() {
	blkCnt := 1
	for { // Process Blocks
		time.Sleep(10 * time.Second) // Create a block for some period of time.
		fmt.Println("EOB", blkCnt)
		// The Control channel only has one space.  Sending true indicates to the accumulator that it is
		// time to seal off a block, do all that indexing, and start the next block.  Sending an immediate
		// false will stall this loop until the accumulator comes around to pull that false out of the
		// channel, but will not stall this loop.  So we send another false immediately so this loop will
		// also stall.  False in the Control panel is really just a noop, so sending the two false indicators
		// keeps the Process Blocks loop here and the block generator in the accumulator in sync.
		for _, ctl := range r.Controls {
			ctl <- true // Send true to Control to end the block
		}
		for _, ctl := range r.Controls {
			ctl <- false // Block timing until the block is processed.
		}
		for _, ctl := range r.Controls {
			ctl <- false // Block timing until the block is processed.
		}
		for i, mdFeed := range r.MDFeeds {
			fmt.Printf("Merkle DAG Root hash for %d is %x\n", i, *<-mdFeed)
		}
		blkCnt++
		var totalEntries, totalChains int64

		for _, acc := range r.ACCs {
			totalEntries += acc.EntryCnt.Load()
			totalChains += acc.ChainCnt.Load()
		}
		secs := time.Now().Unix() - types.StartApp.Unix() + 1
		fmt.Printf("Total Entries Written %s to %s total chains, @ %s tps\n",
			humanize.Comma(totalEntries),
			humanize.Comma(totalChains),
			humanize.Comma(totalEntries/secs))
	}
}

// Init
// Allocate a given number of accumulators to record hashes
func (r *Router) Init(entryHashStream chan node.EntryHash, NumAccumulator int) {
	r.EntryHashStream = entryHashStream
	for i := 0; i < NumAccumulator; i++ {
		acc := new(accumulator.Accumulator)
		r.ACCs = append(r.ACCs, acc)
		db := new(database.DB)
		r.DBs = append(r.DBs, db)
		db.Init(i)
		chainID := types.Hash(sha256.Sum256([]byte(fmt.Sprintf("Accumulator %d", i))))
		entryFeed, control, mdHashes := acc.Init(db, &chainID)
		r.EntryFeeds = append(r.EntryFeeds, entryFeed)
		r.Controls = append(r.Controls, control)
		r.MDFeeds = append(r.MDFeeds, mdHashes)
		go acc.Run()
	}
}

func (r *Router) Run() {
	go r.blockTimer()
	cnt := len(r.ACCs) // Count of accumulators
	for {
		entry := <-r.EntryHashStream
		chainNumber := int(entry.ChainID[0])<<8 + int(entry.ChainID[1])
		idx := chainNumber % cnt
		r.ACCs[idx].GetEntryFeed() <- entry
	}
}
