package main

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/node"

	"github.com/FactomProject/factomd/common/primitives/random"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/accumulator"
)

func main() {
	var Accumulator accumulator.Accumulator
	var chains []types.Hash

	// Calculate and package the AccDID for the Accumulator
	AccDIDHash := sha256.Sum256([]byte("AccVal TestChain"))
	var AccDID types.Hash
	AccDID.Extract(AccDIDHash[:])

	// Initialize the Accumulator
	EntryFeed, Control, MDFeed := Accumulator.Init(&AccDID)

	// Start the Accumulator running
	go Accumulator.Run()

	// Set of a timer to mark the end of blocks as they are processed
	go func() {
		for { // Process Blocks
			fmt.Println("EOB")
			time.Sleep(10 * time.Second) // Create a block for some period of time.
			Control <- true              // Send true to Control to end the block
		}
	}()

	// Validator implementation
	// Just create a series of hashes to be recorded.
	seedHash := sha256.Sum256([]byte("Seed Hash"))
	for {
		chain := random.RandInt() % 50000
		if chain >= len(chains) {
			var h types.Hash
			h.Extract(seedHash[:])
			seedHash = sha256.Sum256(seedHash[:])
			chain = len(chains)
			chains = append(chains, h)
		}
		var eh node.EntryHash
		eh.ChainID = chains[chain]
		eh.EntryHash.Extract(seedHash[:])
		seedHash = sha256.Sum256(seedHash[:])
		EntryFeed <- eh
		select {
		case md := <-MDFeed:
			if md == nil {
				fmt.Println("No MDRoot reported")
			} else {
				fmt.Printf("MD: %x\n", md.Bytes())
			}
		default:
		}
	}
}
