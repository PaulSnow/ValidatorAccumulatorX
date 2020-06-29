package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	router2 "github.com/PaulSnow/ValidatorAccumulator/ValAcc/router"

	"github.com/dustin/go-humanize"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/node"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

const help = "\n\nUsage:     VarAcc [ Entry Limit [ChainsInBlock Limit  [tps Limit] ]\n" +
	"Examples:  VarAcc 100000000          # Sets the Entry Limit to 100 M\n" +
	"           VarAcc 100000000 50000    # Sets the Entry Limit to 100 M and the chain Limit to 50k\n" +
	"           VarAcc 100000 50000 1000  # Sets the Entry Limit to 100k, ChainsInBlock Limit to 50k and tps limit to 1000k"

func main() {

	var chains []types.Hash

	// Default limits
	EntryLimit := 50000000 // EntryLimit where we stop. 50 million is enough to run, not too long a run
	TpsLimit := 1000000000 // A billion is pretty much unlimited, by default.
	ChainLimit := 1000     // How many chains do we spread the entries over.
	if len(os.Args) > 1 {
		var err error
		EntryLimit, err = strconv.Atoi(os.Args[1])
		if err != nil {
			println(help)
			return
		}
		if len(os.Args) > 2 {
			ChainLimit, err = strconv.Atoi(os.Args[2])
			if err != nil {
				println(help)
				return
			}
			if len(os.Args) > 3 {
				TpsLimit, err = strconv.Atoi(os.Args[3])
				if err != nil {
					println(help)
					return
				}
			}
		}
	}
	fmt.Println()
	fmt.Println("=========================")
	fmt.Printf("Entry limit of   %s\n"+
		"Chain limit of   %s\n"+
		"TPS limit of     %s\n",
		humanize.Comma(int64(EntryLimit)),
		humanize.Comma(int64(ChainLimit)),
		humanize.Comma(int64(TpsLimit)))
	fmt.Println("=========================")
	fmt.Println()

	router := new(router2.Router)
	EntryFeed := make(chan node.EntryHash, 10000)
	router.Init(EntryFeed, 4)
	go router.Run()

	// Validator implementation
	// Just create a series of hashes to be recorded.
	seedHash := sha256.Sum256([]byte(fmt.Sprint("Accumulator", rand.Int())))
	blockCount := 0
	time.Sleep(time.Second * 2)

	for i := 0; i < EntryLimit; i++ {
		chain := rand.Int() % ChainLimit
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
		time.Sleep(time.Duration((int64(time.Second) * 4 / 7) / int64(TpsLimit)))

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
