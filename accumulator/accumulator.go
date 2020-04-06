package accumulator

import (
	"crypto/sha256"
	"sync"
)

type HashStream chan []byte

type Chain struct {
	HashList [][]byte
	MR       [][]byte
	Txs      HashStream
	Mux      sync.Mutex
	Count    int64
}

func (c *Chain) Run(txs HashStream) {
	c.Txs = txs
	for {
		c.AddToMR(<-txs)
	}
}

// Add a Hash to a building Merkle Tree
func (c *Chain) AddToMR(hash []byte) {
	c.Mux.Lock()
	c.addToMR2(0, hash)
	c.Count++
	c.Mux.Unlock()
}

func (c *Chain) addToMR2(i int, hash []byte) {
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
	for len(c.MR)*2 < po2 {
		c.AddToMR(c.MR[len(c.MR)-1])
	}
	return c.MR[len(c.MR)-1]
}
