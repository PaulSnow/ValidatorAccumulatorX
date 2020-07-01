package merkleDag

import (
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

// MD
// Collects Hashes from some source, and allows the creation of MD Roots as desired.
type MD struct {
	MD       []*types.Hash // Array of hashes that represent the right edge of the Merkle tree
	HashList []types.Hash  // List of Hashes in the order added to the chain
}

// GetHashList
// Returns the list of Hashes to be stored in the Database so we can create the MD for this
// chain for this block.  In Factom, this is much like a EntryBlock
func (m *MD) GetHashList() (list []byte) {
	list = types.Uint32Bytes(uint32(len(m.HashList)))
	for _, v := range m.HashList {
		list = append(list, v.Bytes()...)
	}
	return list
}

// AddToChain
// Add a Hash to the chain and incrementally build the MD
func (m *MD) AddToChain(hash types.Hash) {

	// We are going through through the MD list and combining hashes, so we have to record the hash first thing
	m.HashList = append(m.HashList, hash) // before it is combined with other hashes already added to MD[].

	// We make sure m.MD ends with a nil entry, because that cuts out most of the corner cases in adding hashes
	if len(m.MD) == 0 || m.MD[len(m.MD)-1] != nil { // If this is the first entry, or the last entry isn't nil
		m.MD = append(m.MD, nil) // then we need to add a nil to the end of m.MD
	}

	// Okay, now we go through m.MD and look for the first nil entry in MD and add our hash there.  But along the
	// way, we take every non-vil entry and combine it with the hash we are adding. Note we ALWAYS have a nil at the
	// end of m.MD so we don't have a end case to deal with.
	for i, v := range m.MD {

		// Look and see if the current spot in MD is open.
		if v == nil { // If it is open, put our hash here and continue.
			m.MD[i] = hash.Copy() // put a pointer to a copy of hash into m.MD
			return                // If we have added the hash to m.MD then we are done.
		}

		// If teh current spot is NOT open, we need to combine the hash we have with the hash on the "left", i.e.
		// the hash already in m.MD
		hash = *v.Combine(hash) // Combine v (left) and hash (right) to get a new combined hash to use forward
		m.MD[i] = nil           // Now that we have combined v and hash, this spot is now empty, so clear it.
	}
}

// GetMDRoot
// Close off the Merkle Directed Acyclic Graph (Merkle DAG or MD)
// We take any trailing hashes in MD, hash them up and combine to create the Merkle Dag Root.
// Getting the closing ListMDRoot is non-destructive, which is useful for some use cases.
func (m *MD) GetMDRoot() (MDRoot *types.Hash) {
	// We go through m.MD and combine any left over hashes in m.MD with each other and the MR.
	// If this is a power of two, that's okay because we will pick up the MR (a balanced MD) and
	// return that, the correct behavior
	for _, v := range m.MD {
		if MDRoot == nil { // We will pick up the first hash in m.MD no matter what.
			MDRoot = v // If we assign a nil over a nil, no harm no foul.  Fewer cases to test this way.
		} else if v != nil { // If MDRoot isn't nil and v isn't nil, we combine them.
			MDRoot = v.Combine(*MDRoot) // v is on the left, MDRoot candidate is on the right, for a new MDRoot
		}
	}
	// We drop out with a MDRoot unless m.MD is zero length, in which case we return a nil (correct)
	// If m.MD has the entries for a power of two, then only one hash (the last) is in m.MD, which we return (correct)
	// If m.MD has a railing nil, we return the trailing entries combined with the last entry in m.MD (correct)
	return MDRoot
}

// PrintMR
// For debugging purposes, it is nice to get a string that shows the nil and non nil entries in c.MD
// Note that the "low order" entries are first in the string, so the binary is going from low order on the left to
// high order going right in the string rather than how binary is normally represented.
func (m *MD) PrintMR() (mr string) {
	for _, v := range m.MD {
		if v != nil {
			mr += "O"
			continue
		}
		mr += "_"
	}
	return mr
}
