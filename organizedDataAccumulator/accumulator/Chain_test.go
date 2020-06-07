package accumulator

import (
	"crypto/sha256"
	"testing"
)

func TestMerkleBuilding(t *testing.T) {
	hash := sha256.Sum256([]byte("testdata"))
	chain := new(Chain)

	// This test depends on the observation that the non blank entries in c.MD must be non-zero
	// for every set bit in the count of the entries added to c.MD.  So all we have to do to check the algorithm
	// here is to add entries, then check that the count of entries added predicts the slots in c.MD that are !nil
	for eCnt := 1; eCnt < 65000; eCnt++ {
		hash := sha256.Sum256(hash[:]) // Get a new hash
		chain.AddToChain(hash)         // Add the new hash to the chain

		cnt := eCnt // Get a count we can shift to compare the current count with c.MD
		for i, v := range chain.MD {
			if cnt&1 == 1 && v == nil { // If I have a bit set, v can't be nil
				t.Errorf("Expected the MD entry %d to be !nil. MD %s and entries %x ", i, chain.PrintMR(), eCnt)
				return
			} else if cnt&1 == 0 && v != nil { // If I have a bit clear, v can't be set
				t.Errorf("Expected the MD entry %d to be nil.  MD %s and entries %x ", i, chain.PrintMR(), eCnt)
				return
			}
			cnt = cnt >> 1
		}

	}
}

func TestMerkleInclusion(t *testing.T) {
	hash := sha256.Sum256([]byte("testdata"))
	chain := new(Chain)

	// This test leverages the fact that GetMDRoot() is non-destructive.  So we build up a
	// a MDRoot up to our limit, but after each additional entry, we redo the process with the entries
	// collected to build a second "copy" MDRoot.  We make sure these are always the same.  So
	// we are checking that a fresh build of a MDRoot produces the same result as one that builds
	// an intermediate MDRoot with each entry.  And we demonstrate that the same MDRoot results
	// given the same entries
	for eCnt := 1; eCnt < 2100; eCnt++ {
		hash := sha256.Sum256(hash[:]) // Get a new hash
		chain.AddToChain(hash)         // Add the new hash to the chain

		MDRoot := chain.GetMDRoot()

		copyChain := new(Chain)
		for _, v := range chain.HashList {
			copyChain.AddToChain(v)
		}
		copyMDRoot := chain.GetMDRoot()

		if *MDRoot != *copyMDRoot {
			t.Error("The same entries in the same order should have the same MDRoot")
			return
		}
	}

	// This test modifies one bit of each entry used to build a MDRoot and demonstrates that doing this
	// always produces a different MDRoot, i.e. we are sensitive to a single bit change to any entry used
	// to build a MDRoot
	for eCnt := 64; eCnt < 65; eCnt++ {
		hash := sha256.Sum256(hash[:]) // Get a new hash
		chain.AddToChain(hash)         // Add the new hash to the chain

		MDRoot := chain.GetMDRoot() // This is the MDRoot of the unmodified data

		for i := 0; i < eCnt; i++ { // Run eCnt tests (one for every entry in chain
			for j := 0; j < eCnt; j++ { // Modify each of the entries in chain and compute a MDRoot
				modChain := new(Chain)
				for i, v := range chain.HashList {
					if i == j {
						v[0] ^= 1 // Flip one bit only in the inputs into the new Chain
					}
					modChain.AddToChain(v)
				}
				modMDRoot := modChain.GetMDRoot()
				if *MDRoot == *modMDRoot {
					t.Error("We modified a bit in the inputs to modChain that was not detected in the MDRoot")
					return
				}
			}
		}
	}
}
