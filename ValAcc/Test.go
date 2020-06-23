package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/accumulator"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/database"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/node"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

const help = "\n\nUsage:     VarAcc [ Entry Limit [Chain Limit  [tps Limit] ]\n" +
	"Examples:  VarAcc 100000000          # Sets the Entry Limit to 100 M\n" +
	"           VarAcc 100000000 50000    # Sets the Entry Limit to 100 M and the chain Limit to 50k\n" +
	"           VarAcc 100000 50000 1000  # Sets the Entry Limit to 100k, Chain Limit to 50k and tps limit to 1000k"

func main() {
	var Accumulator accumulator.Accumulator
	var chains []types.Hash
	DB := new(database.DB)
	DB.Init(0)

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
	fmt.Printf("Running a test with: entry limit of %s chain limit of %s tps limit of %s\n",
		humanize.Comma(int64(EntryLimit)),
		humanize.Comma(int64(ChainLimit)),
		humanize.Comma(int64(TpsLimit)))
	fmt.Println("=========================")
	fmt.Println()
	// Calculate and package the AccDID for the Accumulator
	AccDIDHash := sha256.Sum256([]byte("AccVal TestChain"))
	var AccDID types.Hash
	AccDID.Extract(AccDIDHash[:])

	// Initialize the Accumulator
	EntryFeed, Control, MDFeed := Accumulator.Init(DB, &AccDID)

	// Start the Accumulator running
	go Accumulator.Run()

	// These channels are used to end the test after creating at least 100 Million Entries and recording them in a
	// block.
	endBlocks := make(chan bool, 1)
	endTest := make(chan bool, 1)
	blocks := make(chan int, 1)
	// Set of a timer to mark the end of blocks as they are processed
	go func() {
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
			Control <- true  // Send true to Control to end the block
			Control <- false // Block timing until the block is processed.
			Control <- false // Block timing until the block is processed.
			blkCnt++
			blocks <- blkCnt
			select { // When we are signaled to end the block, then wait for this last block to process
			case <-endBlocks: // then signal that we can end the test, and kill this go routine.
				endTest <- true
				return
			default:
			}
		}
	}()

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
		select {
		case blockCount = <-blocks:
		case md := <-MDFeed:
			if md == nil {
				fmt.Println("No MDRoot reported")
			} else {
				fmt.Printf("MD: %x\n", md.Bytes())
			}
		default:
		}
		time.Sleep(time.Duration(((int64(time.Second) * 4 / 7) / int64(TpsLimit))))

	}
	// Wait for the accumulator to eat up the Entries
	for len(EntryFeed) == 0 {
		time.Sleep(100 * time.Millisecond)
	}
	endBlocks <- true
	<-endTest

	// We have submitted 100 M entries, so indicate it is time to end.
	fmt.Printf("\n====================\nRecorded %s Entries in %s Blocks\n",
		humanize.Comma(int64(EntryLimit)),
		humanize.Comma(int64(blockCount)))
	fmt.Println("Test complete.")
}
