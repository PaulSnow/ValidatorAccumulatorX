// In this test package for the testing of the limits of how much data can be collected
// sorted, and put into Merkle trees,

package main

import (
	"fmt"
	"time"

	"github.com/PaulSnow/LoadTest/accumulator"
	"github.com/PaulSnow/LoadTest/txGenerators"
	"github.com/dustin/go-humanize"
)

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

// Hashes of transactions
var txs chan []byte

func main() {

	router1 := new(accumulator.Router)
	router2 := new(accumulator.Router)
	router3 := new(accumulator.Router)
	router1.Init()
	router2.Init()
	router3.Init()

	Seconds := 60

	go ShowTime()
	go txGenerators.GenSimpleTokenTxs(*router1)
	go txGenerators.GenSimpleTokenTxs(*router2)
	go txGenerators.GenSimpleTokenTxs(*router3)

	time.Sleep(time.Duration(Seconds) * time.Second)

	chains1, count1, pending1 := router1.CloseAll()
	chains2, count2, pending2 := router2.CloseAll()
	chains3, count3, pending3 := router3.CloseAll()
	chains := chains1 + chains2 + chains3
	count := count1 + count2 + count3
	pending := pending1 + pending2 + pending3
	fmt.Printf("\nHashes: %s %s h/s chains: %s pending: %s\n",
		humanize.Comma(int64(count)),
		humanize.Comma(count/int64(Seconds)),
		humanize.Comma(chains),
		humanize.Comma(pending))
}
