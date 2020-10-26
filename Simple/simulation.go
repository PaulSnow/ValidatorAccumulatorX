// In this test package for the testing of the limits of how much data can be collected
// sorted, and put into Merkle trees,

package main

import (
	"fmt"
	"time"

	"github.com/AccumulusNetwork/ValidatorAccumulator/Simple/accumulator"
	"github.com/AccumulusNetwork/ValidatorAccumulator/Simple/txGenerators"

	"github.com/dustin/go-humanize"
)

// Hashes of transactions
var txs chan []byte

func main() {

	// Allocate three routers and initalize them
	router1 := new(accumulator.Router)
	router2 := new(accumulator.Router)
	router3 := new(accumulator.Router)
	router1.Init()
	router2.Init()
	router3.Init()

	// How many seconds we will run the simulation
	Seconds := 60

	go ShowTime()                               // Start the simulator
	go txGenerators.GenSimpleTokenTxs(*router1) // Feed transactions into a router
	go txGenerators.GenSimpleTokenTxs(*router2) // the router collects the hashes of these txs and sends them to an accumulator
	go txGenerators.GenSimpleTokenTxs(*router3)

	time.Sleep(time.Duration(Seconds) * time.Second) // Simulation runs in the go routines this long.

	chains1, count1, pending1 := router1.CloseAll() // Close and collect all the data from the accumulators
	chains2, count2, pending2 := router2.CloseAll()
	chains3, count3, pending3 := router3.CloseAll()

	chains := chains1 + chains2 + chains3     // Add up all the results
	count := count1 + count2 + count3         //
	pending := pending1 + pending2 + pending3 //

	fmt.Printf("\nHashes: %s %s h/s chains: %s pending: %s\n", // Print a summery of the performance
		humanize.Comma(int64(count)),
		humanize.Comma(count/int64(Seconds)),
		humanize.Comma(chains),
		humanize.Comma(pending))
}

// ShowTime
// prints the time every second for feedback as the simulation runs.
func ShowTime() {
	println()
	println()
	second := 0
	for {
		second++
		time.Sleep(1 * time.Second)
		fmt.Printf("\rElapsed Time: %02d:%02d   ", second/60, second%60)
	}
}
