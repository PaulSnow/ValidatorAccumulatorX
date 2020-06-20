package types

import "crypto/sha256"

// Hash
// ===========================================================================
type Hash [32]byte // We are currently using sha256 hashes

// Copy
// Return a copy of the hash
func (h Hash) Copy() *Hash {
	return &h
}

// Bytes
// Return a []byte for the Hash
func (h Hash) Bytes() []byte {
	return h[:]
}

func (h *Hash) Extract(data []byte) []byte {
	copy((*h)[:], data[:32])
	return data[32:]
}

// Combine
// Hash this hash (the left hash) with the given right hash to produce a new hash
func (h Hash) Combine(right Hash) *Hash {
	sum := sha256.New()
	x := sha256.Sum256(h[:]) // Process the left side, i.e. v from this position in c.MD
	sum.Write(x[:])
	x = sha256.Sum256(right[:]) // Process the right side, i.e. whatever hash combinations we have in hash
	sum.Write(x[:])
	var combinedHash Hash
	copy(combinedHash[:], sum.Sum(nil))
	return &combinedHash
}

// GetChainID
// All the chainIDs under an Accumulator are unique to a DID in Factom, and unique
// from all chainIDs in other Accumulators and Factom itself.
//
// A Factom ChainID is H( H(subChainID[0] + H(subChainID[1] + .. + H(subChainID[n]) where
// the ChainID has n subChainIDs.
//
// An Accumulator ChainID is H (AccumulatorDID + 0x00 + H( H(subChainID[0] + H(subChainID[1] + .. + H(subChainID[n]))
// One cannot construct such a ChainID in Factom, nor on another AccumulatorDID.
func GetChainID(AccumulatorDID Hash, SubChainIDs []Hash) (chainID Hash) {
	sum := sha256.New()
	for _, sc := range SubChainIDs {
		h := sha256.Sum256(sc[:])
		sum.Write(h[:])
	}
	combine := sum.Sum(nil)
	sum = sha256.New()
	sum.Write(AccumulatorDID[:])
	sum.Write([]byte{0})
	sum.Write(combine[:])
	copy(chainID[:], sum.Sum(nil))
	return chainID
}
