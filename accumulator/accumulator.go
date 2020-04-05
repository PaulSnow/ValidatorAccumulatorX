package accumulator

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

// In this test package for the testing of the limits of how much data can be collected
// sorted, and put into Merkle trees,

package main

import (
"crypto/sha256"
"fmt"
"math/rand"
"time"
)

type HashStream chan []byte

type Chain struct {
	HashList [][]byte
	MR       [][]byte
	txs      HashStream
	mux      sync.Mutex
}

func (c *Chain) Run(txs HashStream) {
	c.mux.Lock()
	c.txs = txs
	for {
		c.AddToMR(<-txs)
	}
	c.mux.Unlock()
}

// Add a Hash to a building Merkle Tree
func (c *Chain) AddToMR(hash []byte) {
	c.addToMR2(0, hash)
}

func (c *Chain) addToMR2(start int, hash []byte) {
	i := start
	if len(c.MR) == i {
		c.MR = append(c.MR, hash)
		return
	}
	if c.MR[i] == nil {
		c.MR[i] = hash
		return
	}
	h := sha256.New()
	h.Write(c.MR[i][:])
	h.Write(hash)
	c.MR[i] = nil
	c.addToMR2(i+1, h.Sum(nil))
}

func (c *Chain) CloseMR() []byte {
	c.mux.Lock()
	lmr := len(c.MR)
	var bits uint
	for lmr > 0 {
		bits++
		lmr >>= 1
	}
	lmr = len(c.MR)
	po2 := 1 << (bits - 1)
	if po2 == lmr {
		if po2 == 0 {
			return nil
		}
		return c.MR[lmr-1]
	}
	po2 <<= 1
	fmt.Println("length of MR", len(c.MR))
	fmt.Println("Power of 2", po2)
	for len(c.MR)*2 < po2 {
		c.AddToMR(c.MR[len(c.MR)-1])
	}
	return c.MR[len(c.MR)-1]
	c.mux.Unlock()
}
