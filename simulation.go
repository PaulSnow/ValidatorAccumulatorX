// In this test package for the testing of the limits of how much data can be collected
// sorted, and put into Merkle trees,

package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/PaulSnow/LoadTest/accumulator"
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

	txs = make(chan []byte, 1000000)

	Seconds := 30

	chain := new(accumulator.Chain)

	go ShowTime()
	go chain.Run(txs)
	go genTransactions(txs)
	go genTransactions(txs)
	go genTransactions(txs)

	end := time.Now().Add(time.Duration(Seconds) * time.Second)
	cnt := 0
	for time.Now().Before(end) {
		cnt += 100
	}
	fmt.Printf("\nHashes: %d %4.2f h/s \n", cnt, float64(cnt)/float64(Seconds))
	chain.CloseMR()
	fmt.Printf("\nDone:   %d %4.2f h/s \n", cnt, float64(cnt)/float64(Seconds))
}

func genTransactions(txs accumulator.HashStream) {
	// Addresses transacting
	// An initial balance
	addresses := append([]float64{}, 1000000)

	for {
		// pick two addresses
		SAdr := rand.Int() % len(addresses)
		DAdr := rand.Int() % len(addresses)
		if (len(addresses) < 1000 || rand.Float32() < .1) && len(addresses) < 50000 {
			SAdr = 0
			DAdr = len(addresses)
			addresses = append(addresses, 0)
		}
		// Make transfers to different addresses
		if SAdr == DAdr {
			continue
		}

		amt := addresses[SAdr] * (rand.Float64() + 1) / 10

		tx := fmt.Sprintf("src: %d  dest: %d  amt: %f ", SAdr, DAdr, amt)

		addresses[SAdr] -= amt
		addresses[DAdr] += amt

		h := sha256.Sum256([]byte(tx))
		txs <- h[:]
	}

}
