package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

var HashList [][]byte
var MR [][]byte

func AddToMR(start int, hash []byte) {
	i := start
	if len(MR) == i {
		MR = append(MR, hash)
		return
	}
	if MR[i] == nil {
		MR[i] = hash
		return
	}
	h := sha256.New()
	h.Write(MR[i][:])
	h.Write(hash)
	MR[i] = nil
	AddToMR(i+1, h.Sum(nil))
}

func CloseMR() []byte {
	lmr := len(MR)
	var bits uint
	for lmr > 0 {
		bits++
		lmr >>= 1
	}
	lmr = len(MR)
	po2 := 1 << (bits - 1)
	if po2 == lmr {
		if po2 == 0 {
			return nil
		}
		return MR[lmr-1]
	}
	po2 <<= 1
	fmt.Println("length of MR", len(MR))
	fmt.Println("Power of 2", po2)
	for len(MR)*2 < po2 {
		AddToMR(0, MR[len(MR)-1])
	}
	return MR[len(MR)-1]
}

func ShowTime() {
	println()
	println()
	second := 0
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("\rElapsed Time: %02d:%02d   ", second/60, second%60)
		second++
	}
}

// Hashes of transactions
var txs chan []byte

func main() {

	txs = make(chan []byte, 1000)

	Seconds := 30

	go ShowTime()
	go genTransactions()
	go genTransactions()
	go genTransactions()

	end := time.Now().Add(time.Duration(Seconds) * time.Second)
	cnt := 0
	for time.Now().Before(end) {
		for i := 0; i < 100; i++ {
			AddToMR(0, <-txs)
		}
		cnt += 100
	}
	fmt.Printf("\nHashes: %d %4.2f h/s \n", cnt, float64(cnt)/float64(Seconds))
	CloseMR()
	fmt.Printf("\nDone:   %d %4.2f h/s \n", cnt, float64(cnt)/float64(Seconds))
}

func genTransactions() {
	// Addresses transacting
	// An initial balance
	addresses := append([]float64{}, 1000000)

	for {
		// pick two addresses
		s_adr := rand.Int() % len(addresses)
		d_adr := rand.Int() % len(addresses)
		if (len(addresses) < 1000 || rand.Float32() < .1) && len(addresses) < 50000 {
			s_adr = 0
			d_adr = len(addresses)
			addresses = append(addresses, 0)
		}
		// Make transfers to different addresses
		if s_adr == d_adr {
			continue
		}

		amt := addresses[s_adr] * (rand.Float64() + 1) / 10

		tx := fmt.Sprintf("src: %d  dest: %d  amt: %f ", s_adr, d_adr, amt)

		addresses[s_adr] -= amt
		addresses[d_adr] += amt

		h := sha256.Sum256([]byte(tx))
		txs <- h[:]
	}

}
