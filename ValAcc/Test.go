package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/accumulator"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/database"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/node"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

func main() {
	var Accumulator accumulator.Accumulator
	var chains []types.Hash
	DB := new(database.DB)
	DB.Init(0)

	// Calculate and package the AccDID for the Accumulator
	AccDIDHash := sha256.Sum256([]byte("AccVal TestChain"))
	var AccDID types.Hash
	AccDID.Extract(AccDIDHash[:])

	// Initialize the Accumulator
	EntryFeed, Control, MDFeed := Accumulator.Init(DB, &AccDID)

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
	seedHash := sha256.Sum256([]byte(fmt.Sprint("Accumulator", rand.Int())))
	for {
		chain := rand.Int() % 50
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
