//  An accumulator accepts hashes from an input channel, and builds a merkle tree dynamically.
// At the end of a block, the merkle tree is balanced and the Merkle root produced.

package accumulator

import (
	"crypto/sha256"
	"sync"
)

type HashStream chan []byte

// Structure for accepting hashes to go into a chain.
type Chain struct {
	MR    [][]byte   // Array of hashes that represent the right edge of the Merkle tree
	Txs   HashStream // stream of hashes going into the Merkle tree
	Mux   sync.Mutex // Allow the chain to be locked when closing
	Count int64      // How many hashes we have processed.
}

// go routine that pulls transactions out of the hashstream and pushes them into adding to the Merkle tree
func (c *Chain) Run(txs HashStream) {
	c.Txs = txs
	for {
		c.AddToMR(<-txs)
	}
}

// Add a Hash to a building Merkle Tree
func (c *Chain) AddToMR(hash []byte) {
	c.Mux.Lock()
	c.addToMR2(0, hash) // Height of the Merkle tree to add a hash.  New data always inserted at zero
	c.Count++
	c.Mux.Unlock()
}

// Recursive function that drives the hash down the right side of the Merkle tree
func (c *Chain) addToMR2(i int, hash []byte) {
	if len(c.MR) == i { // If inserted where the height is equal to the length of the merkle tree, simply insert this hash
		c.MR = append(c.MR, hash)
		return
	}
	if c.MR[i] == nil { // The Merkle tree is nil if all the hashes at this level have been merged together
		c.MR[i] = hash //  In that case, just insert
		return
	}
	h := sha256.New()   // Combine this hash with the hash at this height in the MR.  Then recurse at the next
	h.Write(c.MR[i][:]) // level with the result.
	h.Write(hash)
	c.MR[i] = nil
	c.addToMR2(i+1, h.Sum(nil))
}

// CloseMR
// Pad the rest of the inputs into the Merkle tree until we match a power of two.
func (c *Chain) CloseMR() []byte {
	lmr := len(c.MR)
	var bits uint

	// Figure out how many bits are required to represent the current depth of the Merkle tree (resolved or not)
	for lmr > 0 {
		bits++
		lmr >>= 1
	}

	// Compute the 2^lmr (the number representing the width of the Merkle tree.
	lmr = len(c.MR)
	po2 := 1 << (bits - 1)
	if po2 == lmr { // If the Merkle tree is balanced, we are done.
		if po2 == 0 {
			return nil
		}
		return c.MR[lmr-1]
	}

	po2 <<= 1               // Go to the next power of two.
	for len(c.MR)*2 < po2 { // For every missing input,
		c.AddToMR(c.MR[len(c.MR)-1]) // Add the last hash again. (could add a zero hash too)
	} //
	return c.MR[len(c.MR)-1] // Return the merkle root
}
