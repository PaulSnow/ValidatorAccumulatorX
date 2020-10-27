package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"math/rand"
	"time"

	router2 "github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/router"

	"github.com/dustin/go-humanize"

	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/node"
	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/types"
)

const help = "\n\nUsage:     VarAcc [ Entry Limit [ChainsInBlock Limit  [tps Limit] ]\n" +
	"Examples:  VarAcc 100000000          # Sets the Entry Limit to 100 M\n" +
	"           VarAcc 100000000 50000    # Sets the Entry Limit to 100 M and the chain Limit to 50k\n" +
	"           VarAcc 100000 50000 1000  # Sets the Entry Limit to 100k, ChainsInBlock Limit to 50k and tps limit to 1000k"

func main() {
	types.StartApp = time.Now()
	var chains []types.Hash

	EntryLimitPtr := flag.Int64("e", 1000000, "the number of entries to be processed in this test")
	ChainLimitPtr := flag.Int64("c", 1000, "The number of chains updated while processing this test")
	TpsLimitPtr := flag.Int64("t", -1, "the tps limit of data generated to run this test. if t < 0, no limit")
	AccNumberPtr := flag.Int64("a", 1, "the number of accumulator instances used in this test")
	flag.Parse()
	EntryLimit := *EntryLimitPtr
	ChainLimit := *ChainLimitPtr
	TpsLimit := *TpsLimitPtr
	AccNumber := *AccNumberPtr

	tpsstr := humanize.Comma(TpsLimit)
	if TpsLimit < 0 {
		tpsstr = "none"
	}

	fmt.Println("=========================")
	fmt.Println(" -e <number of entries>")
	fmt.Println(" -c <number of chains>")
	fmt.Println(" -t <tps limit ( -1 is none)>")
	fmt.Println(" -a <number of accumulators>")
	fmt.Println("=========================")
	fmt.Printf(
		"Entry limit of     %15s\n"+
			"Chain limit of     %15s\n"+
			"TPS limit of       %15s\n"+
			"# of Accumulators  %15d\n",
		humanize.Comma(int64(EntryLimit)),
		humanize.Comma(int64(ChainLimit)),
		tpsstr,
		AccNumber)
	fmt.Println("=========================")
	fmt.Println()

	router := new(router2.Router)
	EntryFeed := make(chan node.EntryHash, 10000)
	router.Init(EntryFeed, int(AccNumber))
	go router.Run()

	// Validator implementation
	// Just create a series of hashes to be recorded.
	seedHash := sha256.Sum256([]byte(fmt.Sprint("Accumulator", rand.Int())))
	blockCount := 0
	time.Sleep(time.Second * 2)
	total := int64(0)
	for i := 0; i < int(EntryLimit); i++ {
		chain := rand.Int63() % ChainLimit
		if int(chain) >= len(chains) {
			var h types.Hash
			h.Extract(seedHash[:])
			seedHash = sha256.Sum256(seedHash[:])
			chain = int64(len(chains))
			chains = append(chains, h)
		}
		var eh node.EntryHash
		eh.ChainID = chains[chain]
		eh.EntryHash.Extract(seedHash[:])
		seedHash = sha256.Sum256(seedHash[:])
		EntryFeed <- eh
		total++
		if i&0xFF == 0 {
			tps := total / (time.Now().Unix() - types.StartApp.Unix() + 1)
			for tps > TpsLimit {
				time.Sleep(time.Second)
				tps = total / (time.Now().Unix() - types.StartApp.Unix() + 1)
			}
		}
	}
	// Wait for the accumulator to eat up the Entries
	for len(EntryFeed) == 0 {
		time.Sleep(100 * time.Millisecond)
	}

	// We have submitted the entries requested, so indicate it is time to end.
	fmt.Printf("\n====================\nRecorded %s Entries in %s Blocks\n",
		humanize.Comma(int64(EntryLimit)),
		humanize.Comma(int64(blockCount)))
	fmt.Println("Test complete.")

	time.Sleep(1 * time.Second)
}
