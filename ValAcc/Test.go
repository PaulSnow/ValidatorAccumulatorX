package main

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/FactomProject/factomd/common/primitives/random"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/merkleDag"
	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/accumulator"
)

func main() {
	var Accumulator accumulator.Accumulator
	var chains []types.Hash

	EntryFeed, Control, MDFeed := Accumulator.Init()
	go Accumulator.Run()
	go func() {
		for { // Process Blocks
			fmt.Println("EOB")
			time.Sleep(10 * time.Second) // Create a block for some period of time.
			Control <- true              // Send true to Control to end the block
		}
	}()

	seedHash := sha256.Sum256([]byte("Seed Hash"))
	for {
		chain := random.RandInt() % 50
		if chain >= len(chains) {
			var h types.Hash
			h.Extract(seedHash[:])
			seedHash = sha256.Sum256(seedHash[:])
			chain = len(chains)
			chains = append(chains, h)
		}
		var eh merkleDag.EntryHash
		eh.ChainID = chains[chain]
		eh.EntryHash.Extract(seedHash[:])
		seedHash = sha256.Sum256(seedHash[:])
		<-MDFeed
		EntryFeed <- eh
	}
}
