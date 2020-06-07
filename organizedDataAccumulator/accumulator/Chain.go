package accumulator

import (
	"sync"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

type Chain struct {
	MD       []*types.Hash // Array of hashes that represent the right edge of the Merkle tree
	Mux      sync.Mutex    // Allow the chain to be locked when closing
	HashList []types.Hash  // List of Hashes in the order added to the chain
}

func (c *Chain) AddToChain(hash types.Hash) {

	// We are going through through the MD list and combining hashes, so we have to record the hash first thing
	c.HashList = append(c.HashList, hash) // before it is combined with other hashes already added to MD[].

	// We make sure c.MD ends with a nil entry, because that cuts out most of the corner cases in adding hashes
	if len(c.MD) == 0 || c.MD[len(c.MD)-1] != nil { // If this is the first entry, or the last entry isn't nil
		c.MD = append(c.MD, nil) // then we need to add a nil to the end of c.MD
	}

	// Okay, now we go through c.MD and look for the first nil entry in MD and add our hash there.  But along the
	// way, we take every non-vil entry and combine it with the hash we are adding. Note we ALWAYS have a nil at the
	// end of c.MD so we don't have a end case to deal with.
	for i, v := range c.MD {

		// Look and see if the current spot in MD is open.
		if v == nil { // If it is open, put our hash here and continue.
			c.MD[i] = hash.Copy() // put a pointer to a copy of hash into c.MD
			return                // If we have added the hash to c.MD then we are done.
		}

		// If teh current spot is NOT open, we need to combine the hash we have with the hash on the "left", i.e.
		// the hash already in c.MD
		hash = *v.Combine(hash) // Combine v (left) and hash (right) to get a new combined hash to use forward
		c.MD[i] = nil           // Now that we have combined v and hash, this spot is now empty, so clear it.
	}
}

// Close off the Merkle Directed Acyclic Graph (Merkle DAG or MD)
// We take any trailing hashes in MD, hash them up and combine to create the Merkle Dag Root.
// Getting the closing MDRoot is non-destructive, which is useful for some use cases.
func (c *Chain) GetMDRoot() (MDRoot *types.Hash) {
	// We go through c.MD and combine any left over hashes in c.MD with each other and the MR.
	// If this is a power of two, that's okay because we will pick up the MR (a balanced MD) and
	// return that, the correct behavior
	for _, v := range c.MD {
		if MDRoot == nil { // We will pick up the first hash in c.MD no matter what.
			MDRoot = v // If we assign a nil over a nil, no harm no foul.  Fewer cases to test this way.
		} else if v != nil { // If MDRoot isn't nil and v isn't nil, we combine them.
			MDRoot = v.Combine(*MDRoot) // v is on the left, MDRoot candidate is on the right, for a new MDRoot
		}
	}
	// We drop out with a MDRoot unless c.MD is zero length, in which case we return a nil (correct)
	// If c.MD has the entries for a power of two, then only one hash (the last) is in c.MD, which we return (correct)
	// If c.MD has a railing nil, we return the trailing entries combined with the last entry in c.MD (correct)
	return MDRoot
}

// PrintMR
// For debugging purposes, it is nice to get a string that shows the nil and non nil entries in c.MD
// Note that the "low order" entries are first in the string, so the binary is going from low order on the left to
// high order going right in the string rather than how binary is normally represented.
func (c *Chain) PrintMR() (mr string) {
	for _, v := range c.MD {
		if v != nil {
			mr += "O"
			continue
		}
		mr += "_"
	}
	return mr
}
